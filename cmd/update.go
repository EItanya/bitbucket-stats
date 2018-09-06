package cmd

import (
	"bitbucket-stats/logger"
	"errors"

	"go.uber.org/zap"

	"github.com/spf13/cobra"
)

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Updates local cache data",
	Long: `Updates to the cache will only occur automatically when the cache is empty.
  This is a way to force update it immediately`,
	Aliases: []string{"u", "U", "sync"},
	RunE:    updateRun,
}

func init() {
	rootCmd.AddCommand(updateCmd)
}

func updateRun(cmd *cobra.Command, args []string) error {
	logger.Log.Info("Update called, beginning update")
	err := client.Update()
	if err != nil {
		logger.Log.DPanicw(
			"Update failed",
			zap.Error(err),
		)
		return errors.New("Update function failed, check logs for details")
	}
	logger.Log.Info("Successfully updated cache")
	return nil
}
