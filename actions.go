package main

import (
	"bitbucket/api"
	"bitbucket/arrays"
	"bitbucket/gui"
	"bitbucket/stats"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/alexeyco/simpletable"

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
	return client.Update()
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
	totalFiles := statsCtx.TotalFileCount - statsCtx.RawFileData["Other"]
	data := make(stats.TableData, 0)
	for key, val := range statsCtx.RawFileData {
		if key != "Other" {
			data = append(data, stats.TableDatum{key, val, (float64(val) / float64(totalFiles)) * 100})
			// fmt.Printf("%s: %d/%d (%.2f%%)\n", key, val, totalFiles, (float64(val)/float64(totalFiles))*100)
		}
	}
	table := &stats.Table{
		Data:  data,
		Table: simpletable.New(),
	}
	table.AddHeader([]string{
		"FILETYPE",
		"# OF FILES",
		"% OF TOTAL",
	})
	table.AddFooter([]string{
		"", fmt.Sprintf("%d", statsCtx.TotalFileCount), "",
	})
	table.SetCellsForData(nil)
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

func guiAction(c *cli.Context) error {
	p := gui.Context{}
	p.Initialize()
	return nil
}

func checkUserBeforeAction(c *cli.Context) error {
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
	client, err = api.Initialize(&user, url)
	if err != nil {
		return err
	}
	log.Println("pre-run actions complete")
	return nil
}

func afterCommandAction(c *cli.Context) error {
	log.Println("Performing post action func")
	return nil
}

func onUsageError(c *cli.Context, err error, isSubcommand bool) error {
	if isSubcommand {
		return err
	}
	return nil
}
