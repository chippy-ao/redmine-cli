package cmd

import (
	"fmt"
	"net/url"

	"github.com/spf13/cobra"
)

var versionsCmd = &cobra.Command{
	Use:   "versions",
	Short: "バージョン一覧を取得する",
	RunE:  runVersions,
}

func init() {
	versionsCmd.Flags().String("project", "", "プロジェクトID (必須)")
	_ = versionsCmd.MarkFlagRequired("project")

	rootCmd.AddCommand(versionsCmd)
}

func runVersions(cmd *cobra.Command, args []string) error {
	c, err := loadClientFromProfile()
	if err != nil {
		return err
	}

	project, _ := cmd.Flags().GetString("project")

	var result any
	path := fmt.Sprintf("/projects/%s/versions.json", url.PathEscape(project))
	if err := c.Get(path, nil, &result); err != nil {
		return fmt.Errorf("バージョン取得エラー: %w", err)
	}

	return outputJSON(result)
}
