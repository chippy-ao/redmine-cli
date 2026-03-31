package cmd

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
)

var updateIssueCmd = &cobra.Command{
	Use:   "update-issue <issue-id>",
	Short: "チケットを更新する",
	Args:  cobra.ExactArgs(1),
	RunE:  runUpdateIssue,
}

func init() {
	updateIssueCmd.Flags().Int("status-id", 0, "ステータス ID")
	updateIssueCmd.Flags().Int("assigned-to-id", 0, "担当者 ID")
	updateIssueCmd.Flags().Int("tracker-id", 0, "トラッカー ID")
	updateIssueCmd.Flags().Int("priority-id", 0, "優先度 ID")
	updateIssueCmd.Flags().String("subject", "", "題名")
	updateIssueCmd.Flags().String("description", "", "説明")
	updateIssueCmd.Flags().Int("category-id", 0, "カテゴリ ID")
	updateIssueCmd.Flags().Int("version-id", 0, "対象バージョン ID")
	updateIssueCmd.Flags().String("notes", "", "コメント")
	updateIssueCmd.Flags().Bool("private-notes", false, "コメントを非公開にする")

	rootCmd.AddCommand(updateIssueCmd)
}

func runUpdateIssue(cmd *cobra.Command, args []string) error {
	issueID, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("無効なチケットID: %s", args[0])
	}

	c, err := loadClientFromProfile()
	if err != nil {
		return err
	}

	issue := map[string]any{}

	if cmd.Flags().Changed("status-id") {
		v, _ := cmd.Flags().GetInt("status-id")
		issue["status_id"] = v
	}
	if cmd.Flags().Changed("assigned-to-id") {
		v, _ := cmd.Flags().GetInt("assigned-to-id")
		issue["assigned_to_id"] = v
	}
	if cmd.Flags().Changed("tracker-id") {
		v, _ := cmd.Flags().GetInt("tracker-id")
		issue["tracker_id"] = v
	}
	if cmd.Flags().Changed("priority-id") {
		v, _ := cmd.Flags().GetInt("priority-id")
		issue["priority_id"] = v
	}
	if cmd.Flags().Changed("subject") {
		v, _ := cmd.Flags().GetString("subject")
		issue["subject"] = v
	}
	if cmd.Flags().Changed("description") {
		v, _ := cmd.Flags().GetString("description")
		issue["description"] = v
	}
	if cmd.Flags().Changed("category-id") {
		v, _ := cmd.Flags().GetInt("category-id")
		issue["category_id"] = v
	}
	if cmd.Flags().Changed("version-id") {
		v, _ := cmd.Flags().GetInt("version-id")
		issue["fixed_version_id"] = v
	}
	if cmd.Flags().Changed("notes") {
		v, _ := cmd.Flags().GetString("notes")
		issue["notes"] = v
	}
	if cmd.Flags().Changed("private-notes") {
		v, _ := cmd.Flags().GetBool("private-notes")
		issue["private_notes"] = v
	}

	if len(issue) == 0 {
		return fmt.Errorf("更新するフィールドを1つ以上指定してください")
	}

	body := map[string]any{"issue": issue}
	path := fmt.Sprintf("/issues/%d.json", issueID)
	if err := c.Put(path, body); err != nil {
		return fmt.Errorf("チケット更新エラー: %w", err)
	}

	return outputJSON(map[string]string{"status": "updated"})
}
