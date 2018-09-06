package main

import (
	"fmt"

	"github.com/urfave/cli"
)

var getCommand = cli.Command{
	Name:        "get",
	Aliases:     []string{"do"},
	Usage:       "bitbucket get",
	Description: "Gets raw data for bitbucket at ***REMOVED***",
	Flags: []cli.Flag{
		userFlag,
	},
}

var statsCommand = cli.Command{
	Name:        "stats",
	Usage:       "bitbucket stats",
	Description: "Gets language stats for bitbucket at ***REMOVED***",
	Flags: []cli.Flag{
		userFlag,
	},
	Action: statsAllAction,
	Before: beforeStatsAction,
	After: func(c *cli.Context) error {
		fmt.Println("Finished getting stats Data")
		return nil
	},
	Subcommands: []cli.Command{
		{
			Name:        "all",
			Usage:       "bitbucket stats all",
			Aliases:     []string{"a"},
			Description: "Gets language stats for all of bitbucket at ***REMOVED***",
			Action:      statsAllAction,
		},
		{
			Name:        "repos",
			Usage:       "bitbucket stats repos",
			Aliases:     []string{"r"},
			Description: "Gets language stats for repos on bitbucket at ***REMOVED***",
			Action:      statsReposAction,
		},
		{
			Name:        "projects",
			Usage:       "bitbucket stats projects",
			Aliases:     []string{"p"},
			Description: "Gets language stats for projects on bitbucket at ***REMOVED***",
			Action:      statsProjectsAction,
		},
		{
			Name:        "files",
			Usage:       "bitbucket stats files",
			Aliases:     []string{"f"},
			Description: "Gets language stats for files on bitbucket at ***REMOVED***",
			Action:      statsAllAction,
		},
		{
			Name:        "node_modules",
			Usage:       "bitbucket stats node_modules",
			Aliases:     []string{"n_m"},
			Description: "Gets repos which contain node_modules (SHAME ON YOU)",
			Action:      statsNodeModulesAction,
		},
	},
}

var updateCommand = cli.Command{
	Name:        "update",
	Aliases:     []string{"sync", "reload"},
	Usage:       "bitbucket update",
	Description: "Sync/Updates remote data",
	Flags: []cli.Flag{
		userFlag,
	},
	Action: updateAction,
	Before: checkUserBeforeAction,
	After: func(c *cli.Context) error {
		fmt.Println("Finished syncing/updating data")
		return nil
	},
}

var cliCommands = []cli.Command{
	statsCommand,
	getCommand,
	updateCommand,
}
