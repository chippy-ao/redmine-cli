package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var projectsCmd = &cobra.Command{
	Use:   "projects",
	Short: "プロジェクト一覧を取得する",
	RunE:  runProjects,
}

func init() {
	projectsCmd.Flags().String("include", "", "含める関連データ")
	projectsCmd.Flags().Int("offset", 0, "オフセット")
	projectsCmd.Flags().Int("limit", 25, "取得件数")

	rootCmd.AddCommand(projectsCmd)
}

func runProjects(cmd *cobra.Command, args []string) error {
	c, err := loadClientFromProfile()
	if err != nil {
		return err
	}

	params := make(map[string]string)

	include, _ := cmd.Flags().GetString("include")
	if include != "" {
		params["include"] = include
	}

	offset, _ := cmd.Flags().GetInt("offset")
	params["offset"] = fmt.Sprintf("%d", offset)

	limit, _ := cmd.Flags().GetInt("limit")
	params["limit"] = fmt.Sprintf("%d", limit)

	var result any
	if err := c.Get("/projects.json", params, &result); err != nil {
		return fmt.Errorf("プロジェクト取得エラー: %w", err)
	}

	return outputJSON(result)
}
