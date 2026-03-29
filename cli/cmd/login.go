package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

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

	rootCmd.AddCommand(loginCmd)
	rootCmd.AddCommand(logoutCmd)
	rootCmd.AddCommand(whoamiCmd)
}
