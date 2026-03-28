package cmd

import (
	"fmt"
	"net/url"

	"github.com/spf13/cobra"
)

var membersCmd = &cobra.Command{
	Use:   "members",
	Short: "プロジェクトメンバー一覧を取得する",
	RunE:  runMembers,
}

func init() {
	membersCmd.Flags().String("project", "", "プロジェクトID (必須)")
	_ = membersCmd.MarkFlagRequired("project")

	rootCmd.AddCommand(membersCmd)
}

func runMembers(cmd *cobra.Command, args []string) error {
	c, err := loadClientFromProfile()
	if err != nil {
		return err
	}

	project, _ := cmd.Flags().GetString("project")

	var result any
	path := fmt.Sprintf("/projects/%s/memberships.json", url.PathEscape(project))
	if err := c.Get(path, nil, &result); err != nil {
		return fmt.Errorf("メンバー取得エラー: %w", err)
	}

	return outputJSON(result)
}
