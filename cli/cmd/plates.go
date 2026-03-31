package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

type LocalPlate struct {
	Slug        string `json:"slug"`
	Name        string `json:"name"`
	Description string `json:"description"`
	RepoURL     string `json:"repo_url"`
	ServerURL   string `json:"server_url"`
}

func localPlatesPath() string {
	return filepath.Join(os.Getenv("HOME"), ".kikplate", "plates.json")
}

func loadLocalPlates() ([]LocalPlate, error) {
	path := localPlatesPath()
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return []LocalPlate{}, nil
		}
		return nil, err
	}
	var plates []LocalPlate
	if err := json.Unmarshal(data, &plates); err != nil {
		return nil, err
	}
	return plates, nil
}

func saveLocalPlates(plates []LocalPlate) error {
	path := localPlatesPath()
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(plates, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}

var platesCmd = &cobra.Command{
	Use:   "plates",
	Short: "Manage local plates",
	Long:  "Add, list, and remove plates tracked locally by the CLI.",
}

var platesAddCmd = &cobra.Command{
	Use:   "add [slug]",
	Short: "Add a plate locally by its slug",
	Long: `Look up a plate on the Kikplate server by its slug (e.g. owner/repo)
and add it to your local plate list.`,
	Example: `  kikplate plates add myorg/my-template`,
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		s, err := NewSession(cmd)
		if err != nil {
			return err
		}

		plate, err := s.FetchPlateBySlug(args[0])
		if err != nil {
			return err
		}

		plates, err := loadLocalPlates()
		if err != nil {
			return fmt.Errorf("cannot load local plates: %w", err)
		}

		for _, p := range plates {
			if p.Slug == plate.Slug {
				fmt.Printf("Plate %q is already added.\n", plate.Slug)
				return nil
			}
		}

		plates = append(plates, *plate)
		if err := saveLocalPlates(plates); err != nil {
			return fmt.Errorf("cannot save plates: %w", err)
		}
		fmt.Printf("Added plate %q (%s)\n", plate.Name, plate.Slug)
		return nil
	},
}

var platesListCmd = &cobra.Command{
	Use:   "list",
	Short: "List locally added plates",
	Long:  "Show all plates that have been added to the local CLI store.",
	RunE: func(cmd *cobra.Command, args []string) error {
		plates, err := loadLocalPlates()
		if err != nil {
			return fmt.Errorf("cannot load local plates: %w", err)
		}
		if len(plates) == 0 {
			fmt.Println("No plates added yet. Use 'kikplate plates add <slug>' to add one.")
			return nil
		}
		t := NewTable("SLUG", "NAME", "DESCRIPTION", "SERVER")
		for _, p := range plates {
			t.Row(p.Slug, p.Name, truncate(p.Description, 50), p.ServerURL)
		}
		t.Print()
		return nil
	},
}

var platesRemoveCmd = &cobra.Command{
	Use:     "remove [slug]",
	Short:   "Remove a plate from the local list",
	Example: `  kikplate plates remove myorg/my-template`,
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		slug := args[0]
		plates, err := loadLocalPlates()
		if err != nil {
			return fmt.Errorf("cannot load local plates: %w", err)
		}
		found := false
		filtered := plates[:0]
		for _, p := range plates {
			if p.Slug == slug {
				found = true
				continue
			}
			filtered = append(filtered, p)
		}
		if !found {
			return fmt.Errorf("plate %q not found in local list", slug)
		}
		if err := saveLocalPlates(filtered); err != nil {
			return fmt.Errorf("cannot save plates: %w", err)
		}
		fmt.Printf("Removed plate %q\n", slug)
		return nil
	},
}

func init() {
	platesCmd.AddCommand(platesAddCmd)
	platesCmd.AddCommand(platesListCmd)
	platesCmd.AddCommand(platesRemoveCmd)
	rootCmd.AddCommand(platesCmd)
}
