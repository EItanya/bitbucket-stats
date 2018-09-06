package cache

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"strings"
)

func removeAllLocalData() error {
	// err := removeLocalFilesData()
	// if err != nil {
	// 	return err
	// }
	// err := removeLocalReposData()
	// if err != nil {
	// 	return err
	// }
	// err = removeLocalProjectsData()
	// if err != nil {
	// 	return err
	// }
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
