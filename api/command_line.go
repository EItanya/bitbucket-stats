package api

import (
	"errors"
	"os"
	"strings"
)

func translateArgs() (*UserInfo, error) {
	args := os.Args[1:]
	if len(args) >= 2 && args[0] == "-c" {
		creds := strings.Split(args[1], ":")
		if len(creds) >= 2 {
			user := UserInfo{Username: creds[0], Password: creds[1]}
			return &user, nil
		}
	}
	return nil, errors.New("Could not properly translate command line args")
}
