package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var validRelationTypes = []string{
	"relates", "duplicates", "blocks", "precedes", "follows", "copied_to",
}

var addRelationCmd = &cobra.Command{
	Use:   "add-relation",
	Short: "チケット間のリレーションを作成する",
	RunE:  runAddRelation,
}

func init() {
	addRelationCmd.Flags().Int("issue-id", 0, "元チケット ID")
	addRelationCmd.Flags().Int("related-id", 0, "関連先チケット ID")
	addRelationCmd.Flags().String("type", "", "リレーション種別 (relates, duplicates, blocks, precedes, follows, copied_to)")
	addRelationCmd.Flags().Int("delay", 0, "遅延日数 (precedes/follows のみ)")

	_ = addRelationCmd.MarkFlagRequired("issue-id")
	_ = addRelationCmd.MarkFlagRequired("related-id")
	_ = addRelationCmd.MarkFlagRequired("type")

	rootCmd.AddCommand(addRelationCmd)
}

func isValidRelationType(t string) bool {
	for _, v := range validRelationTypes {
		if v == t {
			return true
		}
	}
	return false
}

func runAddRelation(cmd *cobra.Command, args []string) error {
	c, err := loadClientFromProfile()
	if err != nil {
		return err
	}

	issueID, _ := cmd.Flags().GetInt("issue-id")
	relatedID, _ := cmd.Flags().GetInt("related-id")
	relationType, _ := cmd.Flags().GetString("type")

	if !isValidRelationType(relationType) {
		return fmt.Errorf("無効なリレーション種別: %s（有効な値: %v）", relationType, validRelationTypes)
	}

	relation := map[string]any{
		"issue_to_id":   relatedID,
		"relation_type": relationType,
	}

	if cmd.Flags().Changed("delay") && (relationType == "precedes" || relationType == "follows") {
		delay, _ := cmd.Flags().GetInt("delay")
		relation["delay"] = delay
	}

	body := map[string]any{"relation": relation}
	var result any
	path := fmt.Sprintf("/issues/%d/relations.json", issueID)
	if err := c.Post(path, body, &result); err != nil {
		return fmt.Errorf("リレーション作成エラー: %w", err)
	}

	return outputJSON(result)
}
