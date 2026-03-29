package cmd

import (
	"fmt"
	"net/http"

	"github.com/spf13/cobra"
)

var verifyCmd = &cobra.Command{
	Use:   "verify [slug]",
	Short: "Verify a submitted plate",
	Long: `Verify a plate by checking its kickplate.yaml verification token.
The plate must have been submitted first and the verification_token
must be present in the repository's kickplate.yaml file.
On success the plate becomes approved, public, and verified.`,
	Example: `  kikplate verify myorg/my-template`,
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		s, err := NewAuthSession(cmd)
		if err != nil {
			return err
		}

		var detail PlateDetail
		if err := s.AuthGetJSON("/plates/"+args[0], nil, &detail); err != nil {
			return err
		}

		var result SubmittedPlate
		if err := s.AuthPostJSON("/plates/"+detail.ID+"/verify", "{}", http.StatusOK, &result); err != nil {
			return err
		}

		fmt.Printf("Plate verified!\n\n")
		fmt.Printf("  Slug:       %s\n", result.Slug)
		fmt.Printf("  Name:       %s\n", result.Name)
		fmt.Printf("  Status:     %s\n", result.Status)
		fmt.Printf("  Visibility: %s\n", result.Visibility)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(verifyCmd)
}
