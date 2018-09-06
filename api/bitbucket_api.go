package api

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

const baseURL = "https://***REMOVED***/rest/api/1.0"

// Client basic type of bitbucket api client
type Client struct {
	api  *API
	User UserInfo
}

func (client *Client) checkUser() error {
	return nil
}

// Initialize sets up bitbucket API
func Initialize(user *UserInfo) (*Client, error) {
	if user != nil {
		return setupClient(user)
	} else if len(os.Args) > 1 {
		user, err := translateArgs()
		if err != nil {
			return nil, err
		}
		return setupClient(user)
	}
	return setupClient(nil)
}

// Update retrieves all data and saves
func (client *Client) Update() error {
	if client.api == nil {
		return errors.New("Must initialize client before attempting any retrievals")
	}
	err := removeAllLocalData()
	if err != nil && !strings.Contains(err.Error(), "no such file or directory") {
		fmt.Println(err.Error())
		return err
	}
	_, err = client.GetProjects(make([]string, 0))
	if err != nil {
		return err
	}
	_, err = client.GetRepos(make(map[string][]string))
	if err != nil {
		return err
	}
	_, err = client.GetFiles(make(map[string][]string))
	if err != nil {
		return err
	}
	return nil
}

func setupClient(user *UserInfo) (*Client, error) {
	if user.Username != "" && user.Password != "" {
		api := &API{
			BaseURL:  baseURL,
			username: user.Username,
			password: user.Password,
		}
		client := &Client{
			api:  api,
			User: *user,
		}
		err := client.checkUser()
		if err != nil {
			return nil, err
		}
		client.checkLocalFiles()

		// api.Timeout = 15 * time.Second
		return client, nil
	}
	err := errors.New("Proper credentials were not supplied to the program")
	return nil, err
}

func (client *Client) checkLocalFiles() {
	errs := make([]error, 0)
	_, err := os.Stat(projectsFilePath)
	errs = append(errs, err)
	_, err = os.Stat(reposFilePath)
	errs = append(errs, err)
	_, err = os.Stat(filesFilePath)
	errs = append(errs, err)
	for _, e := range errs {
		if e != nil && os.IsNotExist(e) {
			client.Update()
			break
		}
	}
}
