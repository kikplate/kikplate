package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"go.yaml.in/yaml/v2"
)

type CLIConfig struct {
	Server ServerConfig `yaml:"server"`
	Auth   AuthConfig   `yaml:"auth"`
}

type ServerConfig struct {
	Address string `yaml:"address"`
}

type AuthConfig struct {
	Token string `yaml:"token"`
}

var defaultConfigPath = filepath.Join(os.Getenv("HOME"), ".kikplate", "config.yaml")

func LoadConfig(path string) (*CLIConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("cannot read config file %s: %w", path, err)
	}
	var cfg CLIConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("cannot parse config file: %w", err)
	}
	if cfg.Server.Address == "" {
		return nil, fmt.Errorf("server.address is required in config file")
	}
	return &cfg, nil
}

func SaveConfig(path string, cfg *CLIConfig) error {
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}

func resolveConfigPath(cmd *cobra.Command) string {
	path, _ := cmd.Root().PersistentFlags().GetString("config")
	if path == "" {
		path = defaultConfigPath
	}
	return path
}

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Create a default config file",
	Long:  "View or initialize the kikplate CLI configuration file.",
}

var configInitCmd = &cobra.Command{
	Use:   "init",
	Short: "Create a default config file",
	Long:  "Create a default kikplate CLI config file at ~/.kikplate/config.yaml",
	Example: `  kikplate config init
  kikplate config init --file /path/to/config.yaml`,
	RunE: func(cmd *cobra.Command, args []string) error {
		path, _ := cmd.Flags().GetString("file")
		if path == "" {
			path = defaultConfigPath
		}
		dir := filepath.Dir(path)
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return fmt.Errorf("cannot create directory %s: %w", dir, err)
		}
		if _, err := os.Stat(path); err == nil {
			return fmt.Errorf("config file already exists at %s", path)
		}
		defaultCfg := CLIConfig{
			Server: ServerConfig{
				Address: "http://localhost:3001",
			},
		}
		data, err := yaml.Marshal(&defaultCfg)
		if err != nil {
			return err
		}
		if err := os.WriteFile(path, data, 0o644); err != nil {
			return fmt.Errorf("cannot write config file: %w", err)
		}
		fmt.Printf("Config file created at %s\n", path)
		return nil
	},
}

var configViewCmd = &cobra.Command{
	Use:   "view",
	Short: "Print current config",
	Long:  "Print the contents of the active kikplate CLI config file.",
	RunE: func(cmd *cobra.Command, args []string) error {
		path, _ := cmd.Root().PersistentFlags().GetString("config")
		if path == "" {
			path = defaultConfigPath
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("cannot read config file %s: %w", path, err)
		}
		fmt.Print(string(data))
		return nil
	},
}

func init() {
	configInitCmd.Flags().StringP("file", "f", "", "Path for the config file (default: ~/.kikplate/config.yaml)")
	configCmd.AddCommand(configInitCmd)
	configCmd.AddCommand(configViewCmd)
	rootCmd.AddCommand(configCmd)
}
