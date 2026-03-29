package cmd

import (
	"fmt"
	"net/http"

	"github.com/spf13/cobra"
)

var submitCmd = &cobra.Command{
	Use:   "submit [repo-url]",
	Short: "Submit a repository as a new plate",
	Long: `Submit a GitHub repository to the Kikplate server.
The repository must contain a kickplate.yaml manifest.
After submission the plate will be in "pending" status until verified.`,
	Example: `  kikplate submit https://github.com/org/repo
  kikplate submit https://github.com/org/repo --branch develop
  kikplate submit https://github.com/org/repo --org <org-id>`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		s, err := NewAuthSession(cmd)
		if err != nil {
			return err
		}

		repoURL := args[0]
		branch, _ := cmd.Flags().GetString("branch")
		orgID, _ := cmd.Flags().GetString("org")

		body := fmt.Sprintf(`{"repo_url":%q`, repoURL)
		if branch != "" {
			body += fmt.Sprintf(`,"branch":%q`, branch)
		}
		if orgID != "" {
			body += fmt.Sprintf(`,"organization_id":%q`, orgID)
		}
		body += "}"

		var plate SubmittedPlate
		if err := s.AuthPostJSON("/plates/repository", body, http.StatusCreated, &plate); err != nil {
			return err
		}

		fmt.Printf("Plate submitted successfully!\n\n")
		fmt.Printf("  Slug:     %s\n", plate.Slug)
		fmt.Printf("  Name:     %s\n", plate.Name)
		fmt.Printf("  Status:   %s\n", plate.Status)
		if plate.VerificationToken != nil {
			fmt.Printf("\nVerification token: %s\n", *plate.VerificationToken)
			fmt.Println("Add this token to your kickplate.yaml as 'verification_token',")
			fmt.Println("then run: kikplate verify " + plate.Slug)
		}
		return nil
	},
}

func init() {
	submitCmd.Flags().String("branch", "", "Git branch (default: main)")
	submitCmd.Flags().String("org", "", "Organization ID to submit under")
	rootCmd.AddCommand(submitCmd)
}
