package cmd

import (
	"encoding/json"
	"os"

	"github.com/chippy-ao/redmine-cli/internal/client"
	"github.com/chippy-ao/redmine-cli/internal/config"
	"github.com/spf13/cobra"
)

var profile string

var rootCmd = &cobra.Command{
	Use:   "redmine-cli",
	Short: "Redmine CLI - Redmine REST API を操作する CLI ツール",
}

func init() {
	rootCmd.PersistentFlags().StringVar(&profile, "profile", "", "使用するプロファイル名")
}

func Execute() error {
	return rootCmd.Execute()
}

func loadClientFromProfile() (*client.Client, error) {
	cfg, err := config.LoadConfig(config.DefaultConfigPath())
	if err != nil {
		return nil, err
	}
	p, err := cfg.GetProfile(profile)
	if err != nil {
		return nil, err
	}
	return client.New(p.URL, p.APIKey), nil
}

func outputJSON(v any) error {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(v)
}
