package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var createIssueCmd = &cobra.Command{
	Use:   "create-issue",
	Short: "Redmine にチケットを作成する",
	RunE:  runCreateIssue,
}

func init() {
	createIssueCmd.Flags().String("project", "", "プロジェクト ID or 識別子")
	createIssueCmd.Flags().String("subject", "", "チケット件名")
	createIssueCmd.Flags().Int("tracker-id", 0, "トラッカー ID")
	createIssueCmd.Flags().Int("status-id", 0, "ステータス ID")
	createIssueCmd.Flags().Int("priority-id", 0, "優先度 ID")
	createIssueCmd.Flags().String("description", "", "チケット説明")
	createIssueCmd.Flags().Int("category-id", 0, "カテゴリ ID")
	createIssueCmd.Flags().Int("version-id", 0, "対象バージョン ID")
	createIssueCmd.Flags().Int("assigned-to-id", 0, "担当者 ID")
	createIssueCmd.Flags().Int("parent-issue-id", 0, "親チケット ID")
	createIssueCmd.Flags().Float64("estimated-hours", 0, "予定工数")
	createIssueCmd.Flags().Bool("private", false, "プライベートチケット")

	_ = createIssueCmd.MarkFlagRequired("project")
	_ = createIssueCmd.MarkFlagRequired("subject")

	rootCmd.AddCommand(createIssueCmd)
}

func runCreateIssue(cmd *cobra.Command, args []string) error {
	c, err := loadClientFromProfile()
	if err != nil {
		return err
	}

	project, _ := cmd.Flags().GetString("project")
	subject, _ := cmd.Flags().GetString("subject")

	issue := map[string]any{
		"project_id": project,
		"subject":    subject,
	}

	if cmd.Flags().Changed("tracker-id") {
		v, _ := cmd.Flags().GetInt("tracker-id")
		issue["tracker_id"] = v
	}
	if cmd.Flags().Changed("status-id") {
		v, _ := cmd.Flags().GetInt("status-id")
		issue["status_id"] = v
	}
	if cmd.Flags().Changed("priority-id") {
		v, _ := cmd.Flags().GetInt("priority-id")
		issue["priority_id"] = v
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
	if cmd.Flags().Changed("assigned-to-id") {
		v, _ := cmd.Flags().GetInt("assigned-to-id")
		issue["assigned_to_id"] = v
	}
	if cmd.Flags().Changed("parent-issue-id") {
		v, _ := cmd.Flags().GetInt("parent-issue-id")
		issue["parent_issue_id"] = v
	}
	if cmd.Flags().Changed("estimated-hours") {
		v, _ := cmd.Flags().GetFloat64("estimated-hours")
		issue["estimated_hours"] = v
	}
	if cmd.Flags().Changed("private") {
		v, _ := cmd.Flags().GetBool("private")
		issue["is_private"] = v
	}

	body := map[string]any{"issue": issue}
	var result any
	if err := c.Post("/issues.json", body, &result); err != nil {
		return fmt.Errorf("チケット作成エラー: %w", err)
	}

	return outputJSON(result)
}
