package main

import (
	"bitbucket/logger"
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Llongfile)
}

func main() {
	// start := time.Now()

	defer logger.Log.Sync()
	logger.Log.Info("Starting up app")

	app := cli.NewApp()
	app.Before = beforeAppSetup
	app.Flags = cliFlags
	app.Commands = cliCommands

	app.Name = "Bitbucket Stats"
	app.Usage = "Gather bitbucket stats"
	app.Action = mainAction
	app.OnUsageError = onUsageError
	app.ExitErrHandler = func(c *cli.Context, err error) {
		if err != nil {
			fmt.Println(err)
			logger.Log.Fatal("Program Error out")
		}
	}
	app.Run(os.Args)
	logger.Log.Info("Program Executed Successfully")

	// fmt.Println(time.Since(start))
}
