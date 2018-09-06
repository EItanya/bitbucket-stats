package main

import (
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
	f, err := os.OpenFile("debug.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	log.SetOutput(f)
	log.Println("Starting up app")

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
			log.Fatalln("Program Error out")
		}
	}
	app.Run(os.Args)
	log.Println("Program Executed Successfully")

	// fmt.Println(time.Since(start))
}
