package cmd

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os/exec"
	"runtime"
	"sort"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Authenticate with a Kikplate server",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		email, _ := cmd.Flags().GetString("email")
		password, _ := cmd.Flags().GetString("password")

		if email == "" || password == "" {
			return fmt.Errorf("provide --email and --password")
		}

		s, err := NewSession(cmd)
		if err != nil {
			return err
		}

		body := fmt.Sprintf(`{"email":%q,"password":%q}`, email, password)
		resp, err := s.Post("/auth/login", body)
		if err != nil {
			return fmt.Errorf("cannot reach server: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			raw, _ := io.ReadAll(resp.Body)
			return fmt.Errorf("login failed (%d): %s", resp.StatusCode, string(raw))
		}

		var result struct {
			Token string `json:"token"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return fmt.Errorf("cannot parse response: %w", err)
		}
		if result.Token == "" {
			return fmt.Errorf("server returned empty token")
		}

		s.Config.Auth.Token = result.Token
		if err := s.SaveConfig(); err != nil {
			return fmt.Errorf("cannot save config: %w", err)
		}
		fmt.Println("Login successful. Token saved.")
		return nil
	},
}

var loginSSOCmd = &cobra.Command{
	Use:   "sso <provider>",
	Short: "Authenticate with a configured SSO provider",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		s, err := NewSession(cmd)
		if err != nil {
			return err
		}

		var providersResponse struct {
			Providers []string `json:"providers"`
		}
		if err := s.GetJSON("/auth/providers", nil, &providersResponse); err != nil {
			return fmt.Errorf("cannot fetch SSO providers: %w", err)
		}

		providers := providersResponse.Providers
		sort.Strings(providers)
		if len(providers) == 0 {
			return fmt.Errorf("no SSO providers are configured on this server")
		}

		if len(args) == 0 {
			return fmt.Errorf("provider is required; available providers: %s", strings.Join(providers, ", "))
		}

		provider, ok := resolveProvider(args[0], providers)
		if !ok {
			return fmt.Errorf("unknown provider %q; available providers: %s", args[0], strings.Join(providers, ", "))
		}

		listener, err := net.Listen("tcp", "127.0.0.1:0")
		if err != nil {
			return fmt.Errorf("cannot start local callback listener: %w", err)
		}

		callbackURL := fmt.Sprintf("http://%s/callback", listener.Addr().String())
		authURL, err := buildSSOLoginURL(s.Addr(), provider, callbackURL)
		if err != nil {
			_ = listener.Close()
			return err
		}

		fmt.Printf("Callback URL: %s\n", callbackURL)
		fmt.Printf("Auth URL: %s\n", authURL)
		fmt.Printf("Opening browser for %s SSO login...\n", provider)
		if err := openBrowser(authURL); err != nil {
			fmt.Printf("Could not open browser automatically. Open this URL manually:\n%s\n", authURL)
		}
		fmt.Println("Waiting for callback...")

		token, err := awaitSSOToken(listener, 3*time.Minute)
		if err != nil {
			return err
		}
		if token == "" {
			return fmt.Errorf("received empty token from SSO callback")
		}

		s.Config.Auth.Token = token
		if err := s.SaveConfig(); err != nil {
			return fmt.Errorf("cannot save config: %w", err)
		}

		fmt.Println("SSO login successful. Token saved.")
		return nil
	},
}

var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Remove stored authentication token",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		s, err := NewSession(cmd)
		if err != nil {
			return err
		}
		if s.Config.Auth.Token == "" {
			fmt.Println("Not logged in.")
			return nil
		}
		s.Config.Auth.Token = ""
		if err := s.SaveConfig(); err != nil {
			return fmt.Errorf("cannot save config: %w", err)
		}
		fmt.Println("Logged out. Token removed.")
		return nil
	},
}

var whoamiCmd = &cobra.Command{
	Use:   "whoami",
	Short: "Show the currently authenticated user",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		s, err := NewAuthSession(cmd)
		if err != nil {
			return err
		}

		var me struct {
			AccountID   string  `json:"account_id"`
			Username    *string `json:"username"`
			DisplayName *string `json:"display_name"`
			Email       *string `json:"email"`
		}
		if err := s.AuthGetJSON("/me", nil, &me); err != nil {
			return err
		}

		if me.Username != nil {
			fmt.Printf("Username:  %s\n", *me.Username)
		}
		if me.DisplayName != nil {
			fmt.Printf("Name:      %s\n", *me.DisplayName)
		}
		if me.Email != nil {
			fmt.Printf("Email:     %s\n", *me.Email)
		}
		fmt.Printf("Account:   %s\n", me.AccountID)
		return nil
	},
}

func init() {
	loginCmd.Flags().String("email", "", "Email address")
	loginCmd.Flags().String("password", "", "Password")
	loginCmd.AddCommand(loginSSOCmd)

	rootCmd.AddCommand(loginCmd)
	rootCmd.AddCommand(logoutCmd)
	rootCmd.AddCommand(whoamiCmd)
}

func resolveProvider(input string, providers []string) (string, bool) {
	for _, p := range providers {
		if strings.EqualFold(input, p) {
			return p, true
		}
	}
	return "", false
}

func buildSSOLoginURL(serverAddr, provider, callbackURL string) (string, error) {
	base, err := url.Parse(strings.TrimRight(serverAddr, "/"))
	if err != nil {
		return "", fmt.Errorf("invalid server address %q: %w", serverAddr, err)
	}
	base.Path = strings.TrimRight(base.Path, "/") + "/auth/" + provider + "/redirect"
	q := base.Query()
	q.Set("cli_callback", base64.RawURLEncoding.EncodeToString([]byte(callbackURL)))
	base.RawQuery = q.Encode()
	return base.String(), nil
}

func awaitSSOToken(listener net.Listener, timeout time.Duration) (string, error) {
	tokenCh := make(chan string, 1)
	flowErrCh := make(chan error, 1)
	serverErrCh := make(chan error, 1)

	mux := http.NewServeMux()
	mux.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		token := strings.TrimSpace(r.URL.Query().Get("token"))
		errParam := strings.TrimSpace(r.URL.Query().Get("error"))

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if errParam != "" {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte("<html><body><h3>SSO login failed</h3><p>You can close this window and retry in the terminal.</p></body></html>"))
			select {
			case flowErrCh <- fmt.Errorf("sso login failed: %s", errParam):
			default:
			}
			return
		}

		if token == "" {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte("<html><body><h3>Missing token</h3><p>No token received from server. This usually means the OAuth callback URL wasn't recognized. Check terminal for details and try again.</p></body></html>"))
			select {
			case flowErrCh <- fmt.Errorf("missing token in SSO callback"):
			default:
			}
			return
		}

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("<html><body><h3>Login successful</h3><p style='color:green'><strong>✓ Token received!</strong></p><p>You can close this window and return to the terminal.</p></body></html>"))
		select {
		case tokenCh <- token:
		default:
		}
	})

	server := &http.Server{Handler: mux}
	go func() {
		if err := server.Serve(listener); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErrCh <- err
		}
	}()

	select {
	case token := <-tokenCh:
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		_ = server.Shutdown(shutdownCtx)
		return token, nil
	case err := <-flowErrCh:
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		_ = server.Shutdown(shutdownCtx)
		return "", err
	case err := <-serverErrCh:
		return "", fmt.Errorf("local callback server failed: %w", err)
	case <-time.After(timeout):
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		_ = server.Shutdown(shutdownCtx)
		return "", fmt.Errorf("timeout waiting for SSO callback")
	}
}

func openBrowser(targetURL string) error {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", targetURL)
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", targetURL)
	default:
		cmd = exec.Command("xdg-open", targetURL)
	}
	return cmd.Start()
}
