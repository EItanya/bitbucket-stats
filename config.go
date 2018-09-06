package main

import (
	"bitbucket/logger"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/urfave/cli"
)

// Config structure of config.json
type Config struct {
	Username string `json:"username"`
	Password string `json:"password"`
	URL      string `json:"url"`
}

func (c *Config) Read(filename string) error {
	byt, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	err = json.Unmarshal(byt, c)
	if err != nil {
		return err
	}
	return nil
}

func beforeAppSetup(c *cli.Context) error {
	logger.Log.Info("Running Setup")
	config := Config{}
	config.Read(c.GlobalString(strings.Split(configFlag.Name, ",")[0]))
	if config.Username != "" && config.Password != "" {
		err := c.GlobalSet("user", fmt.Sprintf("%s:%s", config.Username, config.Password))
		if err != nil {
			return err
		}
	}
	if config.URL != "" {
		err = c.GlobalSet("url", config.URL)
		if err != nil {
			return err
		}
	}
	return nil
}
