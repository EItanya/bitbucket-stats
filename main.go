package main

import (
	"log"
	"os"

	"github.com/urfave/cli"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Llongfile)
}

func main() {
	// start := time.Now()

	app := cli.NewApp()
	app.Before = beforeAppSetup
	app.Flags = cliFlags
	app.Commands = cliCommands

	app.Name = "Bitbucket Stats"
	app.Usage = "Gather bitbucket stats"
	app.Action = mainAction
	// app.OnUsageError = onUsageError
	app.ExitErrHandler = func(c *cli.Context, err error) {
		if err != nil {
			log.Fatalln(err)
		}
	}
	app.Run(os.Args)
	log.Println("Program Executed Successfully")

	// fmt.Println(time.Since(start))
}
