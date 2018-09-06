package main

import (
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
	app.OnUsageError = onUsageError
	app.ExitErrHandler = func(c *cli.Context, err error) {
		log.Fatalln(err)
	}
	app.Run(os.Args)

	// fmt.Println(time.Since(start))
}
