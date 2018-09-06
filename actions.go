package main

import (
	"bitbucket/api"
	"bitbucket/arrays"
	"bitbucket/stats"
	"errors"
	"fmt"
	"log"
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
	client.Update()
	return nil
}

func getAllAction(c *cli.Context) error {
	client.GetFiles(make(map[string][]string))
	return nil
}

func getReposAction(c *cli.Context) error {
	client.GetRepos(make(map[string][]string))
	return nil
}

func getProjectsAction(c *cli.Context) error {
	client.GetProjects(make([]string, 0))
	return nil
}

func statsAllAction(c *cli.Context) error {
	// statsJSON, ok := statsCtx.ToJSON("RawFileData")
	// if ok {
	// 	fmt.Println(statsJSON)
	// 	fmt.Println(statsCtx.TotalFileCount)
	// }
	totalFiles := statsCtx.TotalFileCount - statsCtx.RawFileData["Other"]
	for key, val := range statsCtx.RawFileData {
		if key != "Other" {
			fmt.Printf("%s: %d/%d (%.2f%%)\n", key, val, totalFiles, (float64(val)/float64(totalFiles))*100)
		}
	}
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
		fmt.Printf("(%s:%s)\n", val.ProjectKey, val.RepoSlug)
		for lang, total := range val.Stats.Data {
			fmt.Printf("  %s: %d/%d (%.2f%%)\n", lang, total, val.Stats.Total, (float64(total)/float64(val.Stats.Total))*100)
		}
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
		fmt.Printf("(%s)\n", val.ProjectKey)
		for lang, total := range val.Stats.Data {
			fmt.Printf("  %s: %d/%d (%.2f%%)\n", lang, total, val.Stats.Total, (float64(total)/float64(val.Stats.Total))*100)
		}
	}
	return nil
}

func statsNodeModulesAction(c *cli.Context) error {
	fmt.Println(statsCtx.ReposWithNodeModules())
	return nil
}

func statsLangAction(c *cli.Context) error {
	for _, val := range statsCtx.GetDataByLanguage([]string{"go"}) {
		fmt.Println(*val)
	}
	return nil
}

func beforeStatsAction(c *cli.Context) error {
	err = checkUserBeforeAction(c)
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

func checkUserBeforeAction(c *cli.Context) error {
	if !c.IsSet("user") {
		return errUser
	}
	if !c.IsSet("url") {
		return errURL
	}
	splitUsername := strings.Split(c.String("user"), ":")
	url := c.String("url")
	if len(splitUsername) != 2 {
		return errors.New("Inputted username, credentials was not in the proper format. Should be <username>:<password> was " + c.String("user"))
	} else if url == "" {
		return errors.New("Inputted URL is not in the proper format")
	}
	user := api.UserInfo{
		Username: splitUsername[0],
		Password: splitUsername[1],
	}
	client, err = api.Initialize(&user, url)
	if err != nil {
		return err
	}
	fmt.Println("Good News, User exists")
	return nil
}

func onUsageError(c *cli.Context, err error, isSubcommand bool) error {
	if isSubcommand {
		return err
	}
	if err == errUser {
		log.Panic(err)
	}
	return nil
}
