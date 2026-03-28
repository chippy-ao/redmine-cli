package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var prioritiesCmd = &cobra.Command{
	Use:   "priorities",
	Short: "優先度一覧を取得する",
	RunE:  runPriorities,
}

func init() {
	rootCmd.AddCommand(prioritiesCmd)
}

func runPriorities(cmd *cobra.Command, args []string) error {
	c, err := loadClientFromProfile()
	if err != nil {
		return err
	}

	var result any
	if err := c.Get("/enumerations/issue_priorities.json", nil, &result); err != nil {
		return fmt.Errorf("優先度取得エラー: %w", err)
	}

	return outputJSON(result)
}
