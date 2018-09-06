package main

import (
	"bitbucket/api"
	"bitbucket/stats"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/urfave/cli"
)

var errUser = errors.New("A user must be supplied in order to update data")
var client *api.Client
var statsCtx *stats.Context
var err error

func mainAction(c *cli.Context) error {
	fmt.Println("Main Action is executing")
	return nil
}

func updateAction(c *cli.Context) error {
	splitUsername := strings.Split(c.String("user"), ":")
	user := api.UserInfo{
		Username: splitUsername[0],
		Password: splitUsername[1],
	}
	client, err := api.Initialize(&user)
	if err != nil {
		return err
	}
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
	if splitUsername := strings.Split(c.String("user"), ":"); len(splitUsername) == 2 {
		user := api.UserInfo{
			Username: splitUsername[0],
			Password: splitUsername[1],
		}
		client, err = api.Initialize(&user)
		if err != nil {
			return err
		}
		fmt.Println("Good News, User exists")
		return nil
	}
	return errors.New("Inputted username, credentials was not in the proper format. Should be <username>: password was " + c.String("user"))
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
