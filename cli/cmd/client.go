package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

var httpClient = &http.Client{Timeout: 10 * time.Second}

type Session struct {
	Config  *CLIConfig
	CfgPath string
}

func NewSession(cmd *cobra.Command) (*Session, error) {
	path := resolveConfigPath(cmd)
	cfg, err := LoadConfig(path)
	if err != nil {
		return nil, fmt.Errorf("config error: %w\nRun 'kikplate config init' first", err)
	}
	return &Session{Config: cfg, CfgPath: path}, nil
}

func NewAuthSession(cmd *cobra.Command) (*Session, error) {
	s, err := NewSession(cmd)
	if err != nil {
		return nil, err
	}
	if s.Config.Auth.Token == "" {
		return nil, fmt.Errorf("not logged in — run 'kikplate login' first")
	}
	return s, nil
}

func (s *Session) Addr() string {
	return s.Config.Server.Address
}

func (s *Session) Token() string {
	return s.Config.Auth.Token
}

func (s *Session) Get(path string, query url.Values) (*http.Response, error) {
	u := s.Addr() + path
	if len(query) > 0 {
		u += "?" + query.Encode()
	}
	return httpClient.Get(u)
}

func (s *Session) AuthGet(path string, query url.Values) (*http.Response, error) {
	u := s.Addr() + path
	if len(query) > 0 {
		u += "?" + query.Encode()
	}
	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+s.Token())
	return httpClient.Do(req)
}

func (s *Session) Post(path string, body string) (*http.Response, error) {
	return httpClient.Post(s.Addr()+path, "application/json", strings.NewReader(body))
}

func (s *Session) AuthPost(path string, body string) (*http.Response, error) {
	req, err := http.NewRequest("POST", s.Addr()+path, strings.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.Token())
	return httpClient.Do(req)
}

func (s *Session) AuthPostJSON(path string, body string, expect int, target any) error {
	resp, err := s.AuthPost(path, body)
	if err != nil {
		return fmt.Errorf("cannot reach server: %w", err)
	}
	return decodeJSONStatus(resp, expect, target)
}

func (s *Session) SaveConfig() error {
	return SaveConfig(s.CfgPath, s.Config)
}

func decodeJSON(resp *http.Response, target any) error {
	return decodeJSONStatus(resp, http.StatusOK, target)
}

func decodeJSONStatus(resp *http.Response, expect int, target any) error {
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusUnauthorized {
		return fmt.Errorf("unauthorized — run 'kikplate login' first")
	}
	if resp.StatusCode != expect {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("server returned %d: %s", resp.StatusCode, string(body))
	}
	return json.NewDecoder(resp.Body).Decode(target)
}

func (s *Session) GetJSON(path string, query url.Values, target any) error {
	resp, err := s.Get(path, query)
	if err != nil {
		return fmt.Errorf("cannot reach server: %w", err)
	}
	return decodeJSON(resp, target)
}

func (s *Session) AuthGetJSON(path string, query url.Values, target any) error {
	resp, err := s.AuthGet(path, query)
	if err != nil {
		return fmt.Errorf("cannot reach server: %w", err)
	}
	return decodeJSON(resp, target)
}

func (s *Session) FetchPlateBySlug(slug string) (*LocalPlate, error) {
	resp, err := s.Get("/plates/"+slug, nil)
	if err != nil {
		return nil, fmt.Errorf("cannot reach server: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("plate %q not found on server", slug)
	}
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("server returned %d: %s", resp.StatusCode, string(body))
	}

	var plate struct {
		Slug        string  `json:"slug"`
		Name        string  `json:"name"`
		Description string  `json:"description"`
		RepoURL     *string `json:"repo_url"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&plate); err != nil {
		return nil, fmt.Errorf("cannot parse server response: %w", err)
	}

	return &LocalPlate{
		Slug:        plate.Slug,
		Name:        plate.Name,
		Description: plate.Description,
		RepoURL:     deref(plate.RepoURL),
		ServerURL:   s.Addr(),
	}, nil
}
