package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var trackersCmd = &cobra.Command{
	Use:   "trackers",
	Short: "トラッカー一覧を取得する",
	RunE:  runTrackers,
}

func init() {
	rootCmd.AddCommand(trackersCmd)
}

func runTrackers(cmd *cobra.Command, args []string) error {
	c, err := loadClientFromProfile()
	if err != nil {
		return err
	}

	var result any
	if err := c.Get("/trackers.json", nil, &result); err != nil {
		return fmt.Errorf("トラッカー取得エラー: %w", err)
	}

	return outputJSON(result)
}
