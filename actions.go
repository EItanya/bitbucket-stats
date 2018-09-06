package main

import (
	"bitbucket/api"
	"bitbucket/arrays"
	"bitbucket/cache"
	"bitbucket/gui"
	"bitbucket/logger"
	"bitbucket/stats"
	"errors"
	"fmt"
	"strings"

	"github.com/urfave/cli"
)

var errUser = errors.New("A user must be supplied in order to update data")
var errURL = errors.New("A url must be supplied in order to query the Instance")
var client *api.Client
var statsCtx *stats.Context
var err error

func mainAction(c *cli.Context) error {
	fmt.Println("Main Action is executing")
	return nil
}

func updateAction(c *cli.Context) error {
	// return client.Update()
	return nil
}

func getAllAction(c *cli.Context) error {
	client.GetFiles(make(map[string][]string))
	return nil
}

func getReposAction(c *cli.Context) error {
	client.GetRepos(make([]string, 0))
	return nil
}

func getProjectsAction(c *cli.Context) error {
	client.GetProjects(make([]string, 0))
	return nil
}

func statsAllAction(c *cli.Context) error {
	// totalFiles := statsCtx.TotalFileCount - statsCtx.RawFileData["Other"]
	table := &stats.Table{}
	table.CreateBasicFileTable(statsCtx.RawFileData, statsCtx.TotalFileCount)
	fmt.Println(table.Table.String())
	return nil
}

func statsReposAction(c *cli.Context) error {
	var filter []string
	if c.NArg() > 0 {
		filter = strings.Split(c.Args().First(), ",")
	}
	for _, val := range statsCtx.FileDataByRepo {
		if filter != nil && arrays.IndexOfSTR(filter, val.RepoSlug) == -1 {
			continue
		}
		fmt.Printf("\nProject key: (%s)\n", val.ProjectKey)
		fmt.Printf("Repo slug: (%s)\n", val.RepoSlug)
		table := &stats.Table{}
		table.CreateBasicFileTable(val.Stats.Data, val.Stats.Total)
		fmt.Println(table.Table.String())
	}
	return nil
}

func statsProjectsAction(c *cli.Context) error {
	var filter []string
	if c.NArg() > 0 {
		filter = strings.Split(c.Args().First(), ",")
	}

	for _, val := range statsCtx.FileDataByProject {
		if filter != nil && arrays.IndexOfSTR(filter, val.ProjectKey) == -1 {
			continue
		}
		fmt.Printf("\nProject key: (%s)\n", val.ProjectKey)
		table := &stats.Table{}
		table.CreateBasicFileTable(val.Stats.Data, val.Stats.Total)
		fmt.Println(table.Table.String())
	}
	return nil
}

func statsNodeModulesAction(c *cli.Context) error {
	fmt.Println(statsCtx.ReposWithNodeModules())
	return nil
}

func statsLangAction(c *cli.Context) error {
	var filter []string
	if c.NArg() > 0 {
		filter = strings.Split(c.Args().First(), ",")
	}
	for _, val := range statsCtx.GetDataByLanguage(filter) {
		fmt.Println(val)
	}
	return nil
}

func beforeStatsAction(c *cli.Context) error {
	err = setupClientAction(c)
	if err != nil {
		return err
	}

	statsCtx = &stats.Context{}
	err = statsCtx.Initialize(client)
	statsCtx.CountAllFiles()
	statsCtx.CountFilesByProject()
	statsCtx.CountFilesByRepo()
	if err != nil {
		return err
	}
	return nil
}

func guiAction(c *cli.Context) error {
	p := gui.Context{}
	p.Initialize()
	return nil
}

func setupClientAction(c *cli.Context) error {
	if !c.GlobalIsSet("user") {
		return errUser
	}
	if !c.GlobalIsSet("url") {
		return errURL
	}
	splitUsername := strings.Split(c.GlobalString("user"), ":")
	url := strings.Trim(c.GlobalString("url"), " \n")
	if len(splitUsername) != 2 {
		return errors.New("Inputted credentials was not in the proper format. Should be <username>:<password> was " + c.String("user"))
	} else if url == "" {
		return errors.New("Inputted URL is not in the proper format")
	}
	user := api.UserInfo{
		Username: splitUsername[0],
		Password: splitUsername[1],
	}
	forceReset := c.Command.Name == "update"
	redisCache, err := cache.NewRedisCache(cache.RedisConfig{
		Port:     "6379",
		Protocol: "tcp",
	})
	if err != nil {
		logger.Log.Fatal("Unable to establish connection to cache")
		return err
	}
	client, err = api.Initialize(&user, redisCache, url, forceReset)
	if err != nil {
		return err
	}
	logger.Log.Info("pre-run actions complete")
	return nil
}

func afterCommandAction(c *cli.Context) error {
	logger.Log.Info("Performing post action func")
	return nil
}

func onUsageError(c *cli.Context, err error, isSubcommand bool) error {
	if isSubcommand {
		return err
	}
	return nil
}
