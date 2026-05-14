package lib

import (
	"context"
	"database/sql"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/kickplate/api/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Database struct {
	*gorm.DB
}

type dbConnectionGuard struct {
	sqlDB     *sql.DB
	logger    Logger
	lastCheck atomic.Int64
	checking  atomic.Int32
}

func newDBConnectionGuard(sqlDB *sql.DB, logger Logger) *dbConnectionGuard {
	g := &dbConnectionGuard{
		sqlDB:  sqlDB,
		logger: logger,
	}
	g.lastCheck.Store(time.Now().UnixNano())
	return g
}

func (g *dbConnectionGuard) ensureConnected(ctx context.Context) error {
	const checkInterval = 5 * time.Second
	now := time.Now().UnixNano()
	if now-g.lastCheck.Load() < int64(checkInterval) {
		return nil
	}

	if !g.checking.CompareAndSwap(0, 1) {
		return nil
	}
	defer g.checking.Store(0)

	if err := g.pingWithTimeout(ctx, 2*time.Second); err == nil {
		g.lastCheck.Store(time.Now().UnixNano())
		return nil
	}

	g.logger.Warn("Database ping failed; attempting reconnect checks")
	var lastErr error
	for attempt := 1; attempt <= 5; attempt++ {
		if attempt > 1 {
			time.Sleep(time.Duration(attempt) * time.Second)
		}
		if err := g.pingWithTimeout(ctx, 3*time.Second); err == nil {
			g.lastCheck.Store(time.Now().UnixNano())
			g.logger.Infof("Database connection recovered after %d attempt(s)", attempt)
			return nil
		} else {
			lastErr = err
		}
	}

	return fmt.Errorf("database ping failed after retries: %w", lastErr)
}

func (g *dbConnectionGuard) pingWithTimeout(ctx context.Context, timeout time.Duration) error {
	pctx, cancel := context.WithTimeout(context.Background(), timeout)
	if ctx != nil {
		if deadline, ok := ctx.Deadline(); ok {
			if remaining := time.Until(deadline); remaining < timeout {
				cancel()
				pctx, cancel = context.WithTimeout(context.Background(), remaining)
			}
		}
	}
	defer cancel()
	return g.sqlDB.PingContext(pctx)
}

func installDBConnectionGuard(db *gorm.DB, guard *dbConnectionGuard, logger Logger) {
	hook := func(tx *gorm.DB) {
		if tx == nil || tx.Statement == nil {
			return
		}
		if err := guard.ensureConnected(tx.Statement.Context); err != nil {
			tx.AddError(err)
		}
	}

	if err := db.Callback().Create().Before("gorm:create").Register("kikplate:db_guard_create", hook); err != nil {
		logger.Warnf("failed to register DB guard callback %s: %v", "kikplate:db_guard_create", err)
	}
	if err := db.Callback().Query().Before("gorm:query").Register("kikplate:db_guard_query", hook); err != nil {
		logger.Warnf("failed to register DB guard callback %s: %v", "kikplate:db_guard_query", err)
	}
	if err := db.Callback().Update().Before("gorm:update").Register("kikplate:db_guard_update", hook); err != nil {
		logger.Warnf("failed to register DB guard callback %s: %v", "kikplate:db_guard_update", err)
	}
	if err := db.Callback().Delete().Before("gorm:delete").Register("kikplate:db_guard_delete", hook); err != nil {
		logger.Warnf("failed to register DB guard callback %s: %v", "kikplate:db_guard_delete", err)
	}
	if err := db.Callback().Row().Before("gorm:row").Register("kikplate:db_guard_row", hook); err != nil {
		logger.Warnf("failed to register DB guard callback %s: %v", "kikplate:db_guard_row", err)
	}
	if err := db.Callback().Raw().Before("gorm:raw").Register("kikplate:db_guard_raw", hook); err != nil {
		logger.Warnf("failed to register DB guard callback %s: %v", "kikplate:db_guard_raw", err)
	}
}

func NewDatabase(env Env, logger Logger) Database {
	logger.Info("Initializing database connection")
	username := env.DBUsername
	password := env.DBPassword
	host := env.DBHost
	port := env.DBPort
	dbname := env.DBName

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC", host, username, password, dbname, port)

	var db *gorm.DB
	var err error
	maxRetries := 5
	for i := 0; i < maxRetries; i++ {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
			Logger: logger.GetGormLogger(),
		})
		if err == nil {
			break
		}
		if i < maxRetries-1 {
			logger.Infof("Database connection attempt %d/%d failed, retrying in 2 seconds...", i+1, maxRetries)
			time.Sleep(2 * time.Second)
		}
	}

	if err != nil {
		logger.Info("DSN: ", dsn)
		logger.Panic(err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		logger.Panic(err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetConnMaxLifetime(30 * time.Minute)
	sqlDB.SetConnMaxIdleTime(5 * time.Minute)

	installDBConnectionGuard(db, newDBConnectionGuard(sqlDB, logger), logger)

	logger.Info("Database connection established")

	if err := db.AutoMigrate(
		&model.User{},
		&model.Account{},
		&model.Organization{},
		&model.EmailVerification{},
		&model.Plate{},
		&model.PlateTag{},
		&model.PlateMember{},
		&model.PlateReview{},
		&model.Badge{},
		&model.PlateBadge{},
	); err != nil {
		logger.Panicf("AutoMigrate failed: %v", err)
	}

	logger.Info("Database migrations applied")

	if err := db.Exec(`DROP TABLE IF EXISTS sync_log`).Error; err != nil {
		logger.Panicf("Failed to drop legacy sync_log table: %v", err)
	}

	if err := db.Exec(`CREATE EXTENSION IF NOT EXISTS pg_trgm`).Error; err != nil {
		logger.Warnf("pg_trgm extension: %v", err)
	}

	_ = db.Exec(`ALTER TABLE plate DROP COLUMN IF EXISTS search_vector`)
	if err := db.Exec(`
    ALTER TABLE plate
    ADD COLUMN IF NOT EXISTS search_vector tsvector
    GENERATED ALWAYS AS (
        setweight(to_tsvector('english', coalesce(name, '')), 'A') ||
        setweight(to_tsvector('english', coalesce(description, '')), 'B') ||
        setweight(to_tsvector('english', coalesce(category, '')), 'C')
    ) STORED
`).Error; err != nil {
		logger.Warnf("search_vector: %v", err)
	}

	if err := db.Exec(`
    CREATE INDEX IF NOT EXISTS idx_plate_search
    ON plate USING GIN(search_vector)
`).Error; err != nil {
		logger.Warnf("idx_plate_search: %v", err)
	}

	if err := db.Exec(`
    CREATE INDEX IF NOT EXISTS idx_plate_name_trgm
    ON plate USING GIN(name gin_trgm_ops)
`).Error; err != nil {
		logger.Warnf("idx_plate_name_trgm: %v", err)
	}
	if err := db.Exec(`
    CREATE INDEX IF NOT EXISTS idx_plate_desc_trgm
    ON plate USING GIN(description gin_trgm_ops)
`).Error; err != nil {
		logger.Warnf("idx_plate_desc_trgm: %v", err)
	}

	if err := db.Exec(`
    CREATE INDEX IF NOT EXISTS idx_plate_status_visibility_usecount
    ON plate (status, visibility, bookmark_count DESC)
`).Error; err != nil {
		logger.Warnf("idx_plate_status_visibility_usecount: %v", err)
	}
	if err := db.Exec(`
    CREATE INDEX IF NOT EXISTS idx_plate_tag_tag
    ON plate_tag (tag)
`).Error; err != nil {
		logger.Warnf("idx_plate_tag_tag: %v", err)
	}

	if err := db.Exec(`DELETE FROM plate WHERE type = 'file'`).Error; err != nil {
		logger.Warnf("cleanup file plates: %v", err)
	}
	if err := db.Exec(`ALTER TABLE plate DROP COLUMN IF EXISTS content`).Error; err != nil {
		logger.Warnf("drop plate.content: %v", err)
	}
	if err := db.Exec(`ALTER TABLE plate DROP COLUMN IF EXISTS filename`).Error; err != nil {
		logger.Warnf("drop plate.filename: %v", err)
	}

	logger.Info("Extended migrations applied")

	return Database{
		DB: db,
	}
}

func (d Database) Close() error {
	if d.DB == nil {
		return nil
	}
	sqlDB, err := d.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
