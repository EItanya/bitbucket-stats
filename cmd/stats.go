// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"bitbucket-stats/logger"
	"bitbucket-stats/stats"
	"fmt"

	"github.com/spf13/cobra"
)

var statsCtx *stats.Context

// statsCmd represents the stats command
var statsCmd = &cobra.Command{
	Use:               "stats",
	Short:             "Root command for repo statistics",
	Aliases:           []string{"s", "S"},
	PersistentPreRunE: statsPreRun,
	Run:               statsRun,
}

func init() {
	rootCmd.AddCommand(statsCmd)
}

func statsPreRun(cmd *cobra.Command, args []string) error {
	err := cobraSetup(cmd, args)
	if err != nil {
		return err
	}
	statsCtx = &stats.Context{}
	err = statsCtx.Initialize(client)
	if err != nil {
		return err
	}
	logger.Log.Info("Initialized stats context object")
	statsCtx.CountAllFiles()
	statsCtx.CountFilesByProject()
	statsCtx.CountFilesByRepo()
	logger.Log.Info("Finished compiling all necessary stats")

	return nil
}

func statsRun(cmd *cobra.Command, args []string) {
	table := &stats.Table{}
	table.CreateBasicFileTable(statsCtx.RawFileData, statsCtx.TotalFileCount)
	fmt.Println(table.Table.String())
}

var nodeModulesCmd = &cobra.Command{
	Use:   "n_m",
	Short: "Gets a list of all repos which have included their node_modules (Shame on you)",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(statsCtx.ReposWithNodeModules())
	},
}

func init() {
	statsCmd.AddCommand(nodeModulesCmd)
}
