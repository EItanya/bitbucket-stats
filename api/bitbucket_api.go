package api

import (
	"bitbucket/cache"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/gosuri/uiprogress"
)

// Client basic type of bitbucket api client
type Client struct {
	api   *API
	cache cache.Cache
	User  UserInfo
}

func (client *Client) checkUser() error {
	return nil
}

// Update retrieves all data and saves
func (client *Client) Update() error {
	uiprogress.Start() // start rendering
	if client.api == nil {
		return errors.New("Must initialize client before attempting any retrievals")
	}
	if client.cache == nil {
		return errors.New("Must initialize cache before attempting update")
	}
	log.Println("Clearing data cache")
	err := cache.ClearCache(client.cache)
	if err != nil {
		log.Println("Error while clearing cache\n Repopulation might be slightly incorrect, rerun for assurance")
	}
	log.Println("Data cache cleared successfully, Beginning download")
	log.Println("Downloading data to cache")
	if err != nil && !strings.Contains(err.Error(), "no such file or directory") {
		fmt.Println(err.Error())
		return err
	}

	_, err = client.GetProjects(make([]string, 0))
	if err != nil {
		return err
	}
	_, err = client.GetRepos(make([]string, 0))
	if err != nil {
		return err
	}
	_, err = client.GetFiles(make(map[string][]string))
	if err != nil {
		return err
	}
	uiprogress.Stop()
	return nil
}

func Initialize(user *UserInfo, cache cache.Cache, url string, forceReset bool) (*Client, error) {
	if user.Username != "" && user.Password != "" {
		api := &API{
			BaseURL: url,
			user:    *user,
		}
		client := &Client{
			api:   api,
			cache: cache,
			User:  *user,
		}
		err := client.checkUser()
		if err != nil {
			return nil, err
		}
		if forceReset {
			err = client.Update()
			if err != nil {
				return nil, err
			}
		}
		// else {
		// 	client.checkLocalFiles()
		// }

		// api.Timeout = 15 * time.Second
		return client, nil
	}
	err := errors.New("Proper credentials were not supplied to the program")
	return nil, err
}

// func (client *Client) checkLocalFiles() {
// 	errs := make([]error, 0)
// 	_, err := os.Stat(projectsFilePath)
// 	errs = append(errs, err)
// 	_, err = os.Stat(reposFilePath)
// 	errs = append(errs, err)
// 	_, err = os.Stat(filesFilePath)
// 	errs = append(errs, err)
// 	for _, e := range errs {
// 		if e != nil && os.IsNotExist(e) {
// 			client.Update()
// 			break
// 		}
// 	}
// }
