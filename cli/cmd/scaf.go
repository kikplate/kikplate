package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

func isGitURL(s string) bool {
	return strings.HasPrefix(s, "https://") || strings.HasPrefix(s, "git@") || strings.HasPrefix(s, "ssh://")
}

func runGit(dir string, args ...string) error {
	c := exec.Command("git", args...)
	c.Dir = dir
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	return c.Run()
}

var scafCmd = &cobra.Command{
	Use:   "scaffold [slug] [target]",
	Short: "Scaffold a new project from a plate",
	Args:  cobra.RangeArgs(1, 2),
	RunE: func(cmd *cobra.Command, args []string) error {
		slug := args[0]
		local, _ := cmd.Flags().GetBool("local")
		force, _ := cmd.Flags().GetBool("force")

		plate, err := resolvePlate(cmd, slug)
		if err != nil {
			return err
		}
		if plate.RepoURL == "" {
			return fmt.Errorf("plate %q has no repository URL", slug)
		}

		if local {
			parts := strings.Split(slug, "/")
			return scafLocal(plate, parts[len(parts)-1])
		}
		if len(args) < 2 {
			return fmt.Errorf("provide a target name/URL or use --local")
		}

		target := args[1]
		if isGitURL(target) {
			return scafRemote(plate, target, force)
		}
		return scafLocal(plate, target)
	},
}

func scafLocal(plate *LocalPlate, targetDir string) error {
	absTarget, err := filepath.Abs(targetDir)
	if err != nil {
		return err
	}
	if _, err := os.Stat(absTarget); err == nil {
		return fmt.Errorf("directory %q already exists", absTarget)
	}

	fmt.Printf("Cloning %s into %s ...\n", plate.RepoURL, targetDir)
	if err := runGit(".", "clone", plate.RepoURL, targetDir); err != nil {
		return fmt.Errorf("git clone failed: %w", err)
	}

	cleanKikplateYaml(absTarget)
	writeKikplateOrigin(absTarget, plate)
	stampReadme(absTarget, plate)
	fmt.Printf("Done! Project scaffolded in %s\n", targetDir)
	return nil
}

func remoteHasKikplateOrigin(remoteURL string) bool {
	tmpDir, err := os.MkdirTemp("", "kikplate-check-*")
	if err != nil {
		return false
	}
	defer os.RemoveAll(tmpDir)

	c := exec.Command("git", "clone", "--depth", "1", "--no-checkout", remoteURL, tmpDir)
	c.Stdout = nil
	c.Stderr = nil
	if err := c.Run(); err != nil {
		return false
	}

	checkout := exec.Command("git", "checkout", "HEAD", "--", ".kikplate.origin")
	checkout.Dir = tmpDir
	checkout.Stdout = nil
	checkout.Stderr = nil
	if err := checkout.Run(); err != nil {
		return false
	}

	_, err = os.Stat(filepath.Join(tmpDir, ".kikplate.origin"))
	return err == nil
}

func scafRemote(plate *LocalPlate, remoteURL string, force bool) error {
	if remoteHasKikplateOrigin(remoteURL) {
		if !force {
			return fmt.Errorf("remote %q was already scaffolded from a kikplate plate\nUse --force to overwrite it", remoteURL)
		}
		fmt.Printf("Warning: remote %q already has a scaffold, force-pushing will overwrite it\n", remoteURL)
	}

	tmpDir, err := os.MkdirTemp("", "kikplate-scaf-*")
	if err != nil {
		return fmt.Errorf("cannot create temp directory: %w", err)
	}
	defer os.RemoveAll(tmpDir)

	fmt.Printf("Cloning %s ...\n", plate.RepoURL)
	if err := runGit(".", "clone", plate.RepoURL, tmpDir); err != nil {
		return fmt.Errorf("git clone failed: %w", err)
	}

	cleanKikplateYaml(tmpDir)
	writeKikplateOrigin(tmpDir, plate)
	stampReadme(tmpDir, plate)

	if err := os.RemoveAll(filepath.Join(tmpDir, ".git")); err != nil {
		return fmt.Errorf("cannot remove .git: %w", err)
	}
	if err := runGit(tmpDir, "init"); err != nil {
		return fmt.Errorf("git init failed: %w", err)
	}
	if err := runGit(tmpDir, "add", "."); err != nil {
		return fmt.Errorf("git add failed: %w", err)
	}
	commitMsg := fmt.Sprintf("Initial commit — scaffolded from kikplate plate %q (%s)", plate.Slug, plate.RepoURL)
	if err := runGit(tmpDir, "commit", "-m", commitMsg); err != nil {
		return fmt.Errorf("git commit failed: %w", err)
	}
	if err := runGit(tmpDir, "branch", "-M", "main"); err != nil {
		return fmt.Errorf("git branch rename failed: %w", err)
	}
	if err := runGit(tmpDir, "remote", "add", "origin", remoteURL); err != nil {
		return fmt.Errorf("git remote add failed: %w", err)
	}

	fmt.Printf("Pushing to %s ...\n", remoteURL)
	pushArgs := []string{"push", "-u", "origin", "main"}
	if force {
		pushArgs = []string{"push", "--force", "-u", "origin", "main"}
	}
	if err := runGit(tmpDir, pushArgs...); err != nil {
		return fmt.Errorf("git push failed: %w", err)
	}

	fmt.Printf("Done! Repository created at %s\n", remoteURL)
	return nil
}

func writeKikplateOrigin(dir string, plate *LocalPlate) {
	content := fmt.Sprintf(
		"# This project was scaffolded by kikplate.\n# Source plate: %s\n# Source repo: %s\n# Created at: %s\n",
		plate.Slug, plate.RepoURL, time.Now().UTC().Format(time.RFC3339),
	)
	if err := os.WriteFile(filepath.Join(dir, ".kikplate.origin"), []byte(content), 0o644); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: could not write .kikplate.origin: %v\n", err)
	}
}

func stampReadme(dir string, plate *LocalPlate) {
	readmePath := filepath.Join(dir, "README.md")
	existing, _ := os.ReadFile(readmePath)
	banner := fmt.Sprintf(
		"> **Scaffolded from [%s](%s)** via [kikplate](https://github.com/kickplate/kikplate)\n\n",
		plate.Slug, plate.RepoURL,
	)
	if err := os.WriteFile(readmePath, []byte(banner+string(existing)), 0o644); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: could not stamp README.md: %v\n", err)
	}
}

func cleanKikplateYaml(dir string) {
	path := filepath.Join(dir, "kikplate.yaml")
	if _, err := os.Stat(path); err == nil {
		if err := os.Remove(path); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: could not remove kikplate.yaml: %v\n", err)
		} else {
			fmt.Println("Removed kikplate.yaml")
		}
	}
}

func resolvePlate(cmd *cobra.Command, slug string) (*LocalPlate, error) {
	plates, err := loadLocalPlates()
	if err == nil {
		for _, p := range plates {
			if p.Slug == slug {
				return &p, nil
			}
		}
	}

	s, err := NewSession(cmd)
	if err != nil {
		return nil, fmt.Errorf("plate %q not found locally and config not available: %w\nRun 'kikplate config init' first", slug, err)
	}

	plate, err := s.FetchPlateBySlug(slug)
	if err != nil {
		return nil, fmt.Errorf("plate %q not found locally or on server: %w", slug, err)
	}
	return plate, nil
}

func init() {
	scafCmd.Flags().Bool("local", false, "Clone into the current directory using the plate name")
	scafCmd.Flags().Bool("force", false, "Overwrite existing scaffold target")
	rootCmd.AddCommand(scafCmd)
}
