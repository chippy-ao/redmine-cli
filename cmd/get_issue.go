package cmd

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
)

var getIssueCmd = &cobra.Command{
	Use:   "get-issue <id>",
	Short: "Redmine チケットの詳細を取得する",
	Args:  cobra.ExactArgs(1),
	RunE:  runGetIssue,
}

func init() {
	getIssueCmd.Flags().String("include", "", "含める関連データ (journals,children,relations,attachments,changesets,watchers,allowed_statuses)")

	rootCmd.AddCommand(getIssueCmd)
}

func runGetIssue(cmd *cobra.Command, args []string) error {
	id, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("無効なチケットID: %s", args[0])
	}

	c, err := loadClientFromProfile()
	if err != nil {
		return err
	}

	params := make(map[string]string)
	include, _ := cmd.Flags().GetString("include")
	if include != "" {
		params["include"] = include
	}

	var result any
	path := fmt.Sprintf("/issues/%d.json", id)
	if err := c.Get(path, params, &result); err != nil {
		return fmt.Errorf("チケット取得エラー: %w", err)
	}

	return outputJSON(result)
}
