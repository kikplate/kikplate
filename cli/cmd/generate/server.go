package generate

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/kickplate/cli/internal/config"
	"github.com/spf13/cobra"
)

func resolveAddr(cmd *cobra.Command) (string, string, error) {
	cfgPath, _ := cmd.Root().PersistentFlags().GetString("config")
	if cfgPath == "" {
		cfgPath = config.DefaultConfigPath
	}
	cfg, err := config.Load(cfgPath)
	if err != nil {
		return "", "", fmt.Errorf("config error: %w\nRun 'kikplate config init' first", err)
	}
	return cfg.Server.Address, cfg.Auth.Token, nil
}

func fetchServerSchema(cmd *cobra.Command, slug string) (*plateYAML, error) {
	addr, token, err := resolveAddr(cmd)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodGet, addr+"/generate/"+slug+"/schema", nil)
	if err != nil {
		return nil, fmt.Errorf("cannot build request: %w", err)
	}
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("cannot reach server: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("server returned %d: %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}

	var schema plateYAML
	if err := json.NewDecoder(resp.Body).Decode(&schema); err != nil {
		return nil, fmt.Errorf("cannot parse schema response: %w", err)
	}
	return &schema, nil
}

func generateFromServer(cmd *cobra.Command, slug string, values map[string]any) ([]byte, error) {
	addr, token, err := resolveAddr(cmd)
	if err != nil {
		return nil, err
	}

	body, err := marshalJSON(map[string]any{"values": values})
	if err != nil {
		return nil, fmt.Errorf("cannot encode values: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, addr+"/generate/"+slug, body)
	if err != nil {
		return nil, fmt.Errorf("cannot build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("cannot reach server: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		errBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("server returned %d: %s", resp.StatusCode, strings.TrimSpace(string(errBody)))
	}

	genID := resp.Header.Get("X-Generation-ID")
	if genID != "" {
		fmt.Printf("Generation ID: %s\n", genID)
	}

	return io.ReadAll(resp.Body)
}
