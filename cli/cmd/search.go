package cmd

import (
	"fmt"
	"net/url"

	"github.com/spf13/cobra"
)

var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search plates on the server",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		s, err := NewSession(cmd)
		if err != nil {
			return err
		}

		q := url.Values{}
		if v, _ := cmd.Flags().GetString("name"); v != "" {
			q.Set("search", v)
		}
		if v, _ := cmd.Flags().GetString("category"); v != "" {
			q.Set("category", v)
		}
		if v, _ := cmd.Flags().GetString("tag"); v != "" {
			q.Set("tag", v)
		}
		if v, _ := cmd.Flags().GetInt("limit"); v > 0 {
			q.Set("limit", fmt.Sprintf("%d", v))
		}
		if v, _ := cmd.Flags().GetInt("page"); v > 0 {
			q.Set("page", fmt.Sprintf("%d", v))
		}

		var result PlateListResult
		if err := s.GetJSON("/plates", q, &result); err != nil {
			return err
		}

		if len(result.Data) == 0 {
			fmt.Println("No plates found.")
			return nil
		}

		t := NewTable("SLUG", "NAME", "CATEGORY", "RATING", "VERIFIED")
		for _, p := range result.Data {
			t.Row(p.Slug, p.Name, p.Category, fmt.Sprintf("%.1f", p.AvgRating), boolYesNo(p.IsVerified))
		}
		t.PrintWithTotal(result.Total)
		return nil
	},
}

func init() {
	searchCmd.Flags().String("name", "", "Search by name")
	searchCmd.Flags().String("category", "", "Filter by category")
	searchCmd.Flags().String("tag", "", "Filter by tag")
	searchCmd.Flags().Int("limit", 20, "Results per page")
	searchCmd.Flags().Int("page", 1, "Page number")
	rootCmd.AddCommand(searchCmd)
}
