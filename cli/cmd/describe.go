package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

var describeCmd = &cobra.Command{
	Use:   "describe [slug]",
	Short: "Show details of a plate",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		s, err := NewSession(cmd)
		if err != nil {
			return err
		}

		var p PlateDetail
		if err := s.GetJSON("/plates/"+args[0], nil, &p); err != nil {
			return err
		}

		fmt.Printf("Name:         %s\n", p.Name)
		fmt.Printf("Slug:         %s\n", p.Slug)
		fmt.Printf("Category:     %s\n", p.Category)
		fmt.Printf("Status:       %s\n", p.Status)
		fmt.Printf("Visibility:   %s\n", p.Visibility)

		if d := deref(p.Description); d != "" {
			fmt.Printf("Description:  %s\n", d)
		}
		if r := deref(p.RepoURL); r != "" {
			fmt.Printf("Repo:         %s\n", r)
		}
		if b := deref(p.Branch); b != "" {
			fmt.Printf("Branch:       %s\n", b)
		}

		fmt.Printf("Rating:       %.1f\n", p.AvgRating)
		fmt.Printf("Stars:        %d\n", p.StarCount)
		fmt.Printf("Verified:     %s\n", boolYesNo(p.IsVerified))

		if len(p.Tags) > 0 {
			tags := make([]string, len(p.Tags))
			for i, t := range p.Tags {
				tags[i] = t.Tag
			}
			fmt.Printf("Tags:         %s\n", strings.Join(tags, ", "))
		}

		if p.Owner != nil {
			name := deref(p.Owner.Username)
			if name == "" {
				name = deref(p.Owner.DisplayName)
			}
			if name != "" {
				fmt.Printf("Owner:        %s\n", name)
			}
		}

		if p.Organization != nil {
			fmt.Printf("Organization: %s\n", p.Organization.Name)
		}

		fmt.Printf("Created:      %s\n", p.CreatedAt)
		fmt.Printf("Updated:      %s\n", p.UpdatedAt)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(describeCmd)
}
