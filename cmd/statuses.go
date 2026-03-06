package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var statusesCmd = &cobra.Command{
	Use:   "statuses",
	Short: "ステータス一覧を取得する",
	RunE:  runStatuses,
}

func init() {
	rootCmd.AddCommand(statusesCmd)
}

func runStatuses(cmd *cobra.Command, args []string) error {
	c, err := loadClientFromProfile()
	if err != nil {
		return err
	}

	var result any
	if err := c.Get("/issue_statuses.json", nil, &result); err != nil {
		return fmt.Errorf("ステータス取得エラー: %w", err)
	}

	return outputJSON(result)
}
