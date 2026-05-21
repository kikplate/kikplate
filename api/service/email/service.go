package email

import (
	"bytes"
	"crypto/tls"
	"embed"
	"errors"
	"fmt"
	"html/template"
	"net"
	"net/smtp"
	"path"
	"strconv"
	"strings"

	"github.com/kickplate/api/lib"
)

//go:embed templates/*.html
var templatesFS embed.FS

type Service struct {
	smtp    lib.SMTPConfig
	logger  lib.Logger
	logoURL string
}

func NewService(env lib.Env, logger lib.Logger) *Service {
	return &Service{
		smtp:    env.SMTP,
		logger:  logger,
		logoURL: resolveLogoURL(env),
	}
}

func (s *Service) SendWelcome(to, name string) error {
	return s.send(to, "Welcome!", renderTemplate("welcome.html", s.templateData(map[string]any{
		"Name": name,
	})))
}

func (s *Service) SendVerificationEmail(to, name, verifyURL string) error {
	return s.send(to, "Verify your Kikplate account", renderTemplate("verification_email.html", s.templateData(map[string]any{
		"Name":      name,
		"VerifyURL": verifyURL,
	})))
}

func (s *Service) SendPasswordChanged(to string) error {
	return s.send(to, "Your password was changed", renderTemplate("password_change.html", s.templateData(nil)))
}

func (s *Service) SendLikeNotification(to, likedBy string) error {
	return s.send(to, "Someone liked you!", renderTemplate("like_notification.html", s.templateData(map[string]any{
		"LikedBy": likedBy,
	})))
}

func (s *Service) SendPlateSubmitted(to, name, plateName string) error {
	return s.send(to, "Plate submitted - continue verification", renderTemplate("plate_submitted.html", s.templateData(map[string]any{
		"Name":      name,
		"PlateName": plateName,
	})))
}

func (s *Service) SendPlateVerified(to, name, plateName string) error {
	return s.send(to, "Your plate is verified", renderTemplate("plate_verified.html", s.templateData(map[string]any{
		"Name":      name,
		"PlateName": plateName,
	})))
}

func (s *Service) SendPlateRated(to, plateName, ratedBy string, rating int16) error {
	return s.send(to, "New rating on your plate", renderTemplate("plate_rated.html", s.templateData(map[string]any{
		"PlateName": plateName,
		"RatedBy":   ratedBy,
		"Rating":    rating,
	})))
}

func (s *Service) SendPlateSyncIssue(to, plateName, issue string) error {
	return s.send(to, "Issue syncing your plate", renderTemplate("plate_sync_issue.html", s.templateData(map[string]any{
		"PlateName": plateName,
		"Issue":     issue,
	})))
}

func (s *Service) templateData(extra map[string]any) map[string]any {
	data := map[string]any{
		"LogoURL": s.logoURL,
	}
	for k, v := range extra {
		data[k] = v
	}
	return data
}

func resolveLogoURL(env lib.Env) string {
	const defaultLogoURL = "https://kikplate.dev/kikplate-logo-on-dark.svg"

	logo := strings.TrimSpace(env.Customization.Logo)
	if logo == "" {
		return defaultLogoURL
	}

	if strings.HasPrefix(strings.ToLower(logo), "http://") || strings.HasPrefix(strings.ToLower(logo), "https://") {
		return logo
	}

	base := strings.TrimRight(strings.TrimSpace(env.FrontendURL), "/")
	if base == "" {
		return defaultLogoURL
	}

	if strings.HasPrefix(logo, "/") {
		return base + logo
	}

	return base + "/" + logo
}

func (s *Service) send(to, subject, bodyHTML string) error {
	if !s.smtp.IsConfigured() {
		s.logger.Warnf("smtp is not configured; skipping email to %s (%s)", to, subject)
		return nil
	}

	from := s.smtp.FromEmail
	fromHeader := from
	if s.smtp.FromName != "" {
		fromHeader = fmt.Sprintf("%s <%s>", s.smtp.FromName, s.smtp.FromEmail)
	}

	message := strings.Join([]string{
		fmt.Sprintf("From: %s", fromHeader),
		fmt.Sprintf("To: %s", to),
		fmt.Sprintf("Subject: %s", subject),
		"MIME-Version: 1.0",
		"Content-Type: text/html; charset=UTF-8",
		"",
		bodyHTML,
	}, "\r\n")

	return sendSMTPMail(s.smtp, from, []string{to}, []byte(message))
}

func renderTemplate(name string, data map[string]any) string {
	tpl, err := template.ParseFS(templatesFS, path.Join("templates", name))
	if err != nil {
		return "<p>Unable to render email template.</p>"
	}

	var out bytes.Buffer
	if err := tpl.Execute(&out, data); err != nil {
		return "<p>Unable to render email template.</p>"
	}

	return out.String()
}

func sendSMTPMail(cfg lib.SMTPConfig, from string, to []string, message []byte) error {
	if !cfg.IsConfigured() {
		return nil
	}

	addr := net.JoinHostPort(cfg.Host, strconv.Itoa(cfg.Port))
	tlsConfig := &tls.Config{ServerName: cfg.Host}

	var (
		client *smtp.Client
		err    error
	)

	if cfg.UseStartTL {
		conn, dialErr := net.Dial("tcp", addr)
		if dialErr != nil {
			return dialErr
		}

		client, err = smtp.NewClient(conn, cfg.Host)
		if err != nil {
			_ = conn.Close()
			return err
		}

		if ok, _ := client.Extension("STARTTLS"); ok {
			if err := client.StartTLS(tlsConfig); err != nil {
				_ = client.Close()
				return err
			}
		}
	} else {
		tlsConn, dialErr := tls.Dial("tcp", addr, tlsConfig)
		if dialErr != nil {
			return dialErr
		}

		client, err = smtp.NewClient(tlsConn, cfg.Host)
		if err != nil {
			_ = tlsConn.Close()
			return err
		}
	}

	defer client.Close()

	auth := smtp.PlainAuth("", cfg.Username, cfg.Password, cfg.Host)
	if ok, _ := client.Extension("AUTH"); ok {
		if err := client.Auth(auth); err != nil {
			return err
		}
	}

	if err := client.Mail(from); err != nil {
		return err
	}
	for _, recipient := range to {
		if err := client.Rcpt(recipient); err != nil {
			return err
		}
	}

	writer, err := client.Data()
	if err != nil {
		return err
	}

	if _, err := writer.Write(message); err != nil {
		_ = writer.Close()
		return err
	}

	if err := writer.Close(); err != nil {
		return err
	}

	if err := client.Quit(); err != nil && !errors.Is(err, net.ErrClosed) {
		return err
	}

	return nil
}
