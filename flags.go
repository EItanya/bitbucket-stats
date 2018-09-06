package main

import "github.com/urfave/cli"

var userFlag = cli.StringFlag{
	Name:     "user, u",
	Usage:    "Credentials for the bitbucket instance `<username>:<password>`",
	EnvVar:   "CREDS",
	FilePath: "./.creds",
}

var configFlag = cli.StringFlag{
	Name:  "config, c",
	Usage: "Load configuration from `FILE`",
}

var cliFlags = []cli.Flag{
	userFlag,
	configFlag,
}
