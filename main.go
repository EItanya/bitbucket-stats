package main

import (
	"bitbucket-stats/cmd"
	"bitbucket/logger"
	"fmt"
	"log"
	"time"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Llongfile)
}

func main() {
	start := time.Now()

	defer logger.Log.Sync()
	logger.Log.Info("Starting up app")
	cmd.Execute()
	logger.Log.Info("Program Executed Successfully")

	fmt.Println(time.Since(start))
}
