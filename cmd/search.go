package cmd

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/chippy-ao/redmine-cli/internal/query"
	"github.com/spf13/cobra"
)

var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "Redmine チケットを検索する",
	RunE:  runSearch,
}

func init() {
	searchCmd.Flags().String("keyword", "", "検索キーワード")
	searchCmd.Flags().String("project", "", "プロジェクトID")
	searchCmd.Flags().String("status", "", "ステータスID (open, closed, *, または数値)")
	searchCmd.Flags().String("assigned-to", "", "担当者ID")
	searchCmd.Flags().Int("tracker-id", 0, "トラッカーID")
	searchCmd.Flags().Int("category-id", 0, "カテゴリID")
	searchCmd.Flags().Int("version-id", 0, "バージョンID")
	searchCmd.Flags().String("sort", "", "ソート順 (例: updated_on:desc)")
	searchCmd.Flags().Int("offset", 0, "オフセット")
	searchCmd.Flags().Int("limit", 25, "取得件数")

	rootCmd.AddCommand(searchCmd)
}

func runSearch(cmd *cobra.Command, args []string) error {
	c, err := loadClientFromProfile()
	if err != nil {
		return err
	}

	keyword, _ := cmd.Flags().GetString("keyword")
	project, _ := cmd.Flags().GetString("project")
	status, _ := cmd.Flags().GetString("status")
	assignedTo, _ := cmd.Flags().GetString("assigned-to")
	trackerID, _ := cmd.Flags().GetInt("tracker-id")
	categoryID, _ := cmd.Flags().GetInt("category-id")
	versionID, _ := cmd.Flags().GetInt("version-id")
	sort, _ := cmd.Flags().GetString("sort")
	offset, _ := cmd.Flags().GetInt("offset")
	limit, _ := cmd.Flags().GetInt("limit")

	sp := query.SearchParams{
		ProjectID:      project,
		StatusID:       status,
		AssignedToID:   assignedTo,
		TrackerID:      trackerID,
		CategoryID:     categoryID,
		FixedVersionID: versionID,
		Sort:           sort,
		Offset:         offset,
		Limit:          limit,
	}

	var result interface{}

	if keyword != "" {
		filterQuery := query.BuildFilterQuery(keyword, sp)
		params := query.BuildSearchParams(sp)

		// フィルタで処理済みのキーを通常パラメータから削除
		delete(params, "tracker_id")
		delete(params, "assigned_to_id")
		delete(params, "category_id")
		delete(params, "fixed_version_id")
		if !query.IsStatusSpecial(sp.StatusID) {
			delete(params, "status_id")
		}

		// 残りのパラメータをクエリ文字列に変換
		var remainingParts []string
		if len(params) > 0 {
			v := url.Values{}
			for key, val := range params {
				v.Set(key, val)
			}
			remainingParts = append(remainingParts, v.Encode())
		}

		// フィルタクエリと残りパラメータを結合
		fullQuery := filterQuery
		if len(remainingParts) > 0 {
			fullQuery = strings.Join(append([]string{filterQuery}, remainingParts...), "&")
		}

		if err := c.GetRawQuery("/issues.json", fullQuery, &result); err != nil {
			return fmt.Errorf("検索エラー: %w", err)
		}
	} else {
		params := query.BuildSearchParams(sp)
		if err := c.Get("/issues.json", params, &result); err != nil {
			return fmt.Errorf("検索エラー: %w", err)
		}
	}

	return outputJSON(result)
}
