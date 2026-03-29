package cmd

import (
	"fmt"
	"net/url"

	"github.com/spf13/cobra"
)

var myCmd = &cobra.Command{
	Use:   "my",
	Short: "View your plates, bookmarks, and organizations",
}

var myPlatesCmd = &cobra.Command{
	Use:   "plates",
	Short: "List your plates",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		s, err := NewAuthSession(cmd)
		if err != nil {
			return err
		}

		var me struct {
			AccountID string `json:"account_id"`
		}
		if err := s.AuthGetJSON("/me", nil, &me); err != nil {
			return err
		}

		q := url.Values{}
		q.Set("owner_id", me.AccountID)
		q.Set("limit", "100")

		var result PlateListResult
		if err := s.AuthGetJSON("/plates", q, &result); err != nil {
			return err
		}

		if len(result.Data) == 0 {
			fmt.Println("You have no plates.")
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

var myBookmarksCmd = &cobra.Command{
	Use:   "bookmarks",
	Short: "List your bookmarked plates",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		s, err := NewAuthSession(cmd)
		if err != nil {
			return err
		}

		var result PlateListResult
		if err := s.AuthGetJSON("/plates/bookmarked", nil, &result); err != nil {
			return err
		}

		if len(result.Data) == 0 {
			fmt.Println("No bookmarked plates.")
			return nil
		}

		t := NewTable("SLUG", "NAME", "CATEGORY", "RATING")
		for _, p := range result.Data {
			t.Row(p.Slug, p.Name, p.Category, fmt.Sprintf("%.1f", p.AvgRating))
		}
		t.PrintWithTotal(result.Total)
		return nil
	},
}

var myOrgsCmd = &cobra.Command{
	Use:   "orgs",
	Short: "List your organizations",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		s, err := NewAuthSession(cmd)
		if err != nil {
			return err
		}

		var orgs []OrgItem
		if err := s.AuthGetJSON("/organizations/me", nil, &orgs); err != nil {
			return err
		}

		if len(orgs) == 0 {
			fmt.Println("No organizations.")
			return nil
		}

		t := NewTable("ID", "NAME", "DESCRIPTION")
		for _, o := range orgs {
			t.Row(o.ID, o.Name, truncate(o.Description, 50))
		}
		t.Print()
		return nil
	},
}

var mySyncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Show sync status of your plates",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		s, err := NewAuthSession(cmd)
		if err != nil {
			return err
		}

		var me struct {
			AccountID string `json:"account_id"`
		}
		if err := s.AuthGetJSON("/me", nil, &me); err != nil {
			return err
		}

		q := url.Values{}
		q.Set("owner_id", me.AccountID)
		q.Set("limit", "100")

		var result MyPlateSyncResult
		if err := s.AuthGetJSON("/plates", q, &result); err != nil {
			return err
		}

		if len(result.Data) == 0 {
			fmt.Println("You have no plates.")
			return nil
		}

		t := NewTable("SLUG", "NAME", "STATUS", "SYNC", "LAST SYNCED", "NEXT SYNC")
		for _, p := range result.Data {
			t.Row(
				p.Slug,
				p.Name,
				p.Status,
				deref(p.SyncStatus),
				deref(p.LastSynced),
				deref(p.NextSyncAt),
			)
		}
		t.PrintWithTotal(result.Total)
		return nil
	},
}

func init() {
	myCmd.AddCommand(myPlatesCmd)
	myCmd.AddCommand(myBookmarksCmd)
	myCmd.AddCommand(myOrgsCmd)
	myCmd.AddCommand(mySyncCmd)
	rootCmd.AddCommand(myCmd)
}
