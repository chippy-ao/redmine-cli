package cmd

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
)

var deleteIssueCmd = &cobra.Command{
	Use:   "delete-issue <issue-id>",
	Short: "チケットを削除する",
	Args:  cobra.ExactArgs(1),
	RunE:  runDeleteIssue,
}

func init() {
	rootCmd.AddCommand(deleteIssueCmd)
}

func runDeleteIssue(cmd *cobra.Command, args []string) error {
	issueID, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("無効なチケットID: %s", args[0])
	}

	c, err := loadClientFromProfile()
	if err != nil {
		return err
	}

	path := fmt.Sprintf("/issues/%d.json", issueID)
	if err := c.Delete(path); err != nil {
		return fmt.Errorf("チケット削除エラー: %w", err)
	}

	return outputJSON(map[string]string{"status": "deleted"})
}
