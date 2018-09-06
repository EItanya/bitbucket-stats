package main

import "github.com/urfave/cli"

var userFlag = cli.StringFlag{
	Name:   "user, u",
	Usage:  "Credentials for the bitbucket instance `<username>:<password>`",
	EnvVar: "CREDS",
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

var cacheFlag = cli.StringFlag{
	Name:   "cache",
	Usage:  "Type of cache to be used, `redis` or `file`",
	EnvVar: "CACHE",
	Value:  "redis",
}

var configFlag = cli.StringFlag{
	Name:   "config, c",
	Usage:  "Load configuration from `FILE`",
	EnvVar: "CONFIG",
	Value:  "config.json",
}

var urlFLag = cli.StringFlag{
	Name:   "url",
	Usage:  "`URL` of the bitbucket instance",
	EnvVar: "URL",
}

var cliFlags = []cli.Flag{
	configFlag,
	urlFLag,
	userFlag,
	cacheFlag,
}
