package cmd

type PlateListItem struct {
	Slug       string   `json:"slug"`
	Name       string   `json:"name"`
	Category   string   `json:"category"`
	AvgRating  float64  `json:"avg_rating"`
	IsVerified bool     `json:"is_verified"`
	Tags       []TagObj `json:"tags"`
}

type TagObj struct {
	Tag string `json:"tag"`
}

type PlateDetail struct {
	ID          string   `json:"id"`
	Slug        string   `json:"slug"`
	Name        string   `json:"name"`
	Description *string  `json:"description"`
	Category    string   `json:"category"`
	Status      string   `json:"status"`
	Visibility  string   `json:"visibility"`
	RepoURL     *string  `json:"repo_url"`
	Branch      *string  `json:"branch"`
	AvgRating   float64  `json:"avg_rating"`
	StarCount   int      `json:"star_count"`
	IsVerified  bool     `json:"is_verified"`
	Tags        []TagObj `json:"tags"`
	Owner       *struct {
		Username    *string `json:"username"`
		DisplayName *string `json:"display_name"`
	} `json:"owner"`
	Organization *struct {
		Name string `json:"name"`
	} `json:"organization"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type PlateListResult struct {
	Data  []PlateListItem `json:"data"`
	Total int             `json:"total"`
}

type OrgItem struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type SubmittedPlate struct {
	ID                string  `json:"id"`
	Slug              string  `json:"slug"`
	Name              string  `json:"name"`
	Status            string  `json:"status"`
	Visibility        string  `json:"visibility"`
	VerificationToken *string `json:"verification_token"`
	RepoURL           *string `json:"repo_url"`
	Branch            *string `json:"branch"`
}

type MyPlateSync struct {
	Slug       string  `json:"slug"`
	Name       string  `json:"name"`
	Status     string  `json:"status"`
	SyncStatus *string `json:"sync_status"`
	SyncError  *string `json:"sync_error"`
	NextSyncAt *string `json:"next_sync_at"`
	LastSynced *string `json:"last_synced_at"`
	IsVerified bool    `json:"is_verified"`
}

type MyPlateSyncResult struct {
	Data  []MyPlateSync `json:"data"`
	Total int           `json:"total"`
}
