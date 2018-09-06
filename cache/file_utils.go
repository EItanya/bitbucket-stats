package cache

import (
	"bitbucket-stats/logger"
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"strings"
)

func removeAllLocalData(dirname string) error {
	fileInfo, err := ioutil.ReadDir(dirname)
	if err != nil && os.IsNotExist(err) {
		logger.Log.Infof("Directory: (%s) does not exist to be deleted. Trying default directory", dirname)
		fileInfo, err = ioutil.ReadDir(defaultDir)
		if err != nil {
			return err
		}
		err = deleteFiles(fileInfo)
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	} else {
		err = deleteFiles(fileInfo)
		if err != nil {
			return err
		}
	}

	return nil
}

func deleteFiles(fileInfo []os.FileInfo) error {
	for _, item := range fileInfo {
		var deleteErr error
		if item.IsDir() {
			deleteErr = os.Remove(item.Name())
		} else {
			deleteErr = os.Remove(item.Name())
		}
		if deleteErr != nil {
			return deleteErr
		}
	}
	return nil
}

func writeJSONToFile(data interface{}, filename string) error {
	if !strings.Contains(filename, ".json") {
		return errors.New("write file must be a valid json filename")
	}
	byt, err := json.MarshalIndent(&data, "", "  ")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(filename, byt, 0600)
	if err != nil {
		return err
	}
	return nil
}

func readJSONFromFile(dataType interface{}, filename string) error {
	if !strings.Contains(filename, ".json") {
		return errors.New("Read file must be a valid json filename")
	}
	byt, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	err = json.Unmarshal(byt, dataType)
	if err != nil {
		return err
	}
	return nil
}
