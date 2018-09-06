package main

import (
	"github.com/urfave/cli"
)

var getCommand = cli.Command{
	Name:        "get",
	Aliases:     []string{"do"},
	Usage:       "bitbucket get",
	Description: "Gets raw data for bitbucket at supplied instance",
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
}

var updateCommand = cli.Command{
	Name:        "update",
	Aliases:     []string{"sync", "reload"},
	Usage:       "bitbucket update",
	Description: "Sync/Updates remote data",
	Action:      updateAction,
	Before:      setupClientAction,
	After:       afterCommandAction,
}

var guiCommand = cli.Command{
	Name:        "gui",
	Aliases:     []string{"cli"},
	Usage:       "bitbucket gui",
	Description: "Starts the interactive prompt",
	Action:      guiAction,
	Before:      setupClientAction,
	After:       afterCommandAction,
}

// var redisCommand = cli.Command{
// 	Name:        "redis",
// 	Aliases:     []string{"r"},
// 	Usage:       "bitbucket redis",
// 	Description: "Connects to redis",
// 	Action: func(c *cli.Context) error {
// 		redisCache := &cache.RedisCache{
// 			Config: &cache.RedisConfig{
// 				Port:     "6379",
// 				Protocol: "tcp",
// 			},
// 		}
// 		err := cache.InitializeCache(redisCache)
// 		if err != nil {
// 			return err
// 		}
// 		n, err := redisCache.Conn.Do("INCR", "counter")
// 		log.Println(n)
// 		return err
// 	},
// }

var cliCommands = []cli.Command{
	statsCommand,
	getCommand,
	updateCommand,
	guiCommand,
	// redisCommand,
}
