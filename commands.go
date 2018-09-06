package main

import (
	"github.com/urfave/cli"
)

var getCommand = cli.Command{
	Name:        "get",
	Aliases:     []string{"do"},
	Usage:       "bitbucket get",
	Description: "Gets raw data for bitbucket at supplied instance",
	// Flags: []cli.Flag{
	// 	userFlag,
	// },
}

// var statsCommandFlags = []cli.Flag{
// 	userFlag,
// 	urlFLag,
// 	configFlag,
// }
var statsCommand = cli.Command{
	Name:        "stats",
	Usage:       "bitbucket stats",
	Description: "Gets language stats for bitbucket at supplied instance",
	// Flags:       statsCommandFlags,
	Action: statsAllAction,
	Before: beforeStatsAction,
	After:  afterCommandAction,
	Subcommands: []cli.Command{
		{
			Name:        "all",
			Usage:       "bitbucket stats all",
			Aliases:     []string{"a"},
			Description: "Gets language stats for all of bitbucket at supplied instance",
			Action:      statsAllAction,
		},
		{
			Name:        "repos",
			Usage:       "bitbucket stats repos",
			Aliases:     []string{"r"},
			Description: "Gets language stats for repos on bitbucket at supplied instance",
			Action:      statsReposAction,
		},
		{
			Name:        "projects",
			Usage:       "bitbucket stats projects",
			Aliases:     []string{"p"},
			Description: "Gets language stats for projects on bitbucket at supplied instance",
			Action:      statsProjectsAction,
		},
		{
			Name:        "files",
			Usage:       "bitbucket stats files",
			Aliases:     []string{"f"},
			Description: "Gets language stats for files on bitbucket at supplied instance",
			Action:      statsAllAction,
		},
		{
			Name:        "languages",
			Usage:       "bitbucket stats files",
			Aliases:     []string{"l", "langs", "lang"},
			Description: "Gets stats for supplied languages",
			Action:      statsLangAction,
		},
		{
			Name:        "node_modules",
			Usage:       "bitbucket stats node_modules",
			Aliases:     []string{"n_m"},
			Description: "Gets repos which contain node_modules (SHAME ON YOU)",
			Action:      statsNodeModulesAction,
		},
	},
	OnUsageError: onUsageError,
}

var updateCommand = cli.Command{
	Name:        "update",
	Aliases:     []string{"sync", "reload"},
	Usage:       "bitbucket update",
	Description: "Sync/Updates remote data",
	// Flags: []cli.Flag{
	// 	userFlag,
	// 	urlFLag,
	// },
	Action:       updateAction,
	Before:       checkUserBeforeAction,
	After:        afterCommandAction,
	OnUsageError: onUsageError,
}

var guiCommand = cli.Command{
	Name:        "gui",
	Aliases:     []string{"cli"},
	Usage:       "bitbucket gui",
	Description: "Starts the interactive prompt",
	Action:      guiAction,
	// Flags: []cli.Flag{
	// 	userFlag,
	// 	urlFLag,
	// },
	Before:       checkUserBeforeAction,
	After:        afterCommandAction,
	OnUsageError: onUsageError,
}

var cliCommands = []cli.Command{
	statsCommand,
	getCommand,
	updateCommand,
	guiCommand,
}
