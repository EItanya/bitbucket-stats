package cmd

import (
	"bitbucket-stats/arrays"
	"bitbucket-stats/stats"
	"fmt"

	"github.com/spf13/cobra"
)

// projectsCmd represents the projects command
var projectsCmd = &cobra.Command{
	Use:     "projects",
	Short:   "Root command for stats command revolving around projects",
	Aliases: []string{"p", "P"},
	Run:     statsProjectsRun,
}

func init() {
	statsCmd.AddCommand(projectsCmd)
}

func statsProjectsRun(cmd *cobra.Command, args []string) {
	for _, val := range statsCtx.FileDataByProject {
		if len(args) > 0 && arrays.IndexOfSTR(args, val.ProjectKey) == -1 {
			continue
		}
		fmt.Printf("\nProject key: (%s)\n", val.ProjectKey)
		table := &stats.Table{}
		table.CreateBasicFileTable(val.Stats.Data, val.Stats.Total)
		fmt.Println(table.Table.String())
	}
}
