package main

import "github.com/urfave/cli"

var userFlag = cli.StringFlag{
	Name:     "user, u",
	Usage:    "Credentials for the bitbucket instance `<username>:<password>`",
	EnvVar:   "CREDS",
	FilePath: "./.creds",
}

var reposFlag = cli.StringFlag{
	Name:   "repos, r",
	Usage:  "List of comma seperated repos to query",
	EnvVar: "REPOS",
}

var projectsFlag = cli.StringFlag{
	Name:   "projects, p",
	Usage:  "List of comma seperated projects to query",
	EnvVar: "PROJECTS",
}

var configFlag = cli.StringFlag{
	Name:  "config, c",
	Usage: "Load configuration from `FILE`",
}

var cliFlags = []cli.Flag{
	userFlag,
	configFlag,
}
