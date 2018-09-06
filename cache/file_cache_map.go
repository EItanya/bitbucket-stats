package cache

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"reflect"
)

var projectsCacheTable *fileCacheMap
var repositoriesCacheTable *fileCacheMap
var filesCacheTable *fileCacheMap

var fileCacheMaps []*fileCacheMap

func init() {
	projectsCacheTable = &fileCacheMap{
		data: make(fileCacheMapData),
	}
	projectsCacheTable.filename = fmt.Sprintf("%s/%s.json", defaultDir, projectConst)
	projectsCacheTable.title = projectConst
	repositoriesCacheTable = &fileCacheMap{
		data: make(fileCacheMapData),
	}
	repositoriesCacheTable.filename = fmt.Sprintf("%s/%s.json", defaultDir, repositoryConst)
	repositoriesCacheTable.title = projectConst
	filesCacheTable = &fileCacheMap{
		data: make(fileCacheMapData),
	}
	filesCacheTable.filename = fmt.Sprintf("%s/%s.json", defaultDir, filesConst)
	filesCacheTable.title = projectConst
	fileCacheMaps = []*fileCacheMap{
		projectsCacheTable,
		repositoriesCacheTable,
		filesCacheTable,
	}
}

type fileCacheMap struct {
	*fileCacheTableBasic
	data fileCacheMapData
}
type fileCacheMapData map[fileCacheKey]interface{}

func (t *fileCacheTableBasic) write(data interface{}) error {
	if t.filename == "" {
		return errors.New("Filename of Cache Table cannot be null")
	}
	if data == nil {
		return errors.New("Data object to be written cannot be null")
	}
	t.fileMtx.Lock()
	byt, err := json.MarshalIndent(&data, "", "  ")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(t.filename, byt, 0600)
	if err != nil {
		return err
	}
	t.fileMtx.Unlock()
	return nil
}

func (t *fileCacheTableBasic) read(data interface{}) error {
	if t.filename == "" {
		return errors.New("Filename of Cache Table cannot be null")
	}
	if data == nil {
		return errors.New("Data object to be read into cannot be null")
	}
	if reflect.ValueOf(data).Kind() != reflect.Ptr {
		return errors.New("Data object to be read into must be a pointer")
	}
	t.fileMtx.RLock()
	byt, err := ioutil.ReadFile(t.filename)
	if err != nil {
		return err
	}
	err = json.Unmarshal(byt, data)
	if err != nil {
		return err
	}
	t.fileMtx.RUnlock()
	return nil
}

func (t *fileCacheTableBasic) get(key fileCacheKey, currentCacheData map[fileCacheKey]interface{}) (interface{}, error) {
	var data interface{}
	t.dataMtx.RLock()
	data = currentCacheData[key]
	t.dataMtx.RUnlock()
	return data, nil
}

func (t *fileCacheTableBasic) set(key fileCacheKey, currentCacheData, dataToSave interface{}) error {
	switch typedData := currentCacheData.(type) {
	case map[fileCacheKey]interface{}:
		t.dataMtx.Lock()
		typedData[key] = dataToSave
		t.dataMtx.Unlock()
	default:
		return errors.New("No behavior specified for given data type in cache set")
	}
	return nil
}
