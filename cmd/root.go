package cmd

import (
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
