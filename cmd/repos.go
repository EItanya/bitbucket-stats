package cmd

import (
	"bitbucket-stats/arrays"
	"bitbucket-stats/stats"
	"fmt"

	"github.com/spf13/cobra"
)

// reposCmd represents the repos command
var reposCmd = &cobra.Command{
	Use:     "repos",
	Short:   "Root command for stats command revolving around repositories",
	Aliases: []string{"r", "R"},
	Run:     statsReposRun,
}

func init() {
	statsCmd.AddCommand(reposCmd)
}

func statsReposRun(cmd *cobra.Command, args []string) {
	for _, val := range statsCtx.FileDataByRepo {
		if len(args) > 0 && arrays.IndexOfSTR(args, val.RepoSlug) == -1 {
			continue
		}
		fmt.Printf("\nProject key: (%s)\n", val.ProjectKey)
		fmt.Printf("Repo slug: (%s)\n", val.RepoSlug)
		table := &stats.Table{}
		table.CreateBasicFileTable(val.Stats.Data, val.Stats.Total)
		fmt.Println(table.Table.String())
	}
}
