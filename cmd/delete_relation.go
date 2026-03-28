package cmd

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
)

var deleteRelationCmd = &cobra.Command{
	Use:   "delete-relation <relation-id>",
	Short: "リレーションを削除する",
	Args:  cobra.ExactArgs(1),
	RunE:  runDeleteRelation,
}

func init() {
	rootCmd.AddCommand(deleteRelationCmd)
}

func runDeleteRelation(cmd *cobra.Command, args []string) error {
	relationID, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("無効なリレーションID: %s", args[0])
	}

	c, err := loadClientFromProfile()
	if err != nil {
		return err
	}

	path := fmt.Sprintf("/relations/%d.json", relationID)
	if err := c.Delete(path); err != nil {
		return fmt.Errorf("リレーション削除エラー: %w", err)
	}

	return outputJSON(map[string]string{"status": "deleted"})
}
