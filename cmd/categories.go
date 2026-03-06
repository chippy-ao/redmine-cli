package cmd

import (
	"fmt"
	"net/url"

	"github.com/spf13/cobra"
)

var categoriesCmd = &cobra.Command{
	Use:   "categories",
	Short: "カテゴリ一覧を取得する",
	RunE:  runCategories,
}

func init() {
	categoriesCmd.Flags().String("project", "", "プロジェクトID (必須)")
	_ = categoriesCmd.MarkFlagRequired("project")

	rootCmd.AddCommand(categoriesCmd)
}

func runCategories(cmd *cobra.Command, args []string) error {
	c, err := loadClientFromProfile()
	if err != nil {
		return err
	}

	project, _ := cmd.Flags().GetString("project")

	var result any
	path := fmt.Sprintf("/projects/%s/issue_categories.json", url.PathEscape(project))
	if err := c.Get(path, nil, &result); err != nil {
		return fmt.Errorf("カテゴリ取得エラー: %w", err)
	}

	return outputJSON(result)
}
