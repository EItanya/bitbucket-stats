package main

import (
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli"
)

func main() {
	// start := time.Now()
	app := cli.NewApp()
	app.Flags = cliFlags
	app.Commands = cliCommands

	app.Name = "Bitbucket Stats"
	app.Usage = "Gather bitbucket stats"
	app.Action = mainAction
	app.OnUsageError = func(c *cli.Context, err error, isSubcommand bool) error {
		if isSubcommand {
			return err
		}

		fmt.Printf("Error occurred: %#v\n", err)
		return nil
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}

	// client, err := api.Initialize(nil)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// // err = client.Update()
	// files, err := client.GetFiles(make(map[string][]string))
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// repos, err := client.GetRepos(make(map[string][]string))
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// projects, err := client.GetProjects(make([]string, 0))
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// count := stats.Languages{}
	// count.CountAllFiles(&files)
	// count.CountFilesByRepo(&files, &repos)
	// count.CountFilesByProject(&files, &projects)
	// fmt.Println(time.Since(start))
}
