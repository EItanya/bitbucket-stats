package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// filesCmd represents the files command
var filesCmd = &cobra.Command{
	Use:     "files",
	Short:   "Root command for stats command revolving around files",
	Aliases: []string{"f", "F"},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("files called")
	},
}

func init() {
	statsCmd.AddCommand(filesCmd)
}
