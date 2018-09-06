package cache

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"reflect"
	"sync"
)

var projectsCacheTable *fileCacheMap
var repositoriesCacheTable *fileCacheMap
var filesCacheTable *fileCacheMap

var fileCacheMaps []*fileCacheMap

func init() {
	projectsCacheTable = new(fileCacheMap)
	projectsCacheTable.initialize(defaultDir, projectConst)

	repositoriesCacheTable = new(fileCacheMap)
	repositoriesCacheTable.initialize(defaultDir, repositoryConst)

	filesCacheTable = new(fileCacheMap)
	filesCacheTable.initialize(defaultDir, filesConst)

	fileCacheMaps = []*fileCacheMap{
		projectsCacheTable,
		repositoriesCacheTable,
		filesCacheTable,
	}
}

func (t *fileCacheMap) initialize(dir, name string) {
	t.data = make(fileCacheMapData)
	t.fileCacheTableBasic = &fileCacheTableBasic{
		filename: fmt.Sprintf("%s/%s.json", dir, name),
		title:    name,
		dataMtx:  &sync.RWMutex{},
		fileMtx:  &sync.RWMutex{},
	}
}

type fileCacheMap struct {
	*fileCacheTableBasic
	data fileCacheMapData
}
type fileCacheMapData map[fileCacheKey]interface{}

func (t *fileCacheMap) write(data interface{}) error {
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

func (t *fileCacheMap) read(data interface{}) error {
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

func (t *fileCacheMap) get(key fileCacheKey) (interface{}, error) {
	var data interface{}
	t.dataMtx.RLock()
	data = t.data[key]
	t.dataMtx.RUnlock()
	return data, nil
}

func (t *fileCacheMap) set(key fileCacheKey, dataToSave interface{}) error {
	t.dataMtx.Lock()
	t.data[key] = dataToSave
	t.dataMtx.Unlock()
	return nil
}

func (t *fileCacheMap) keys() []string {
	keys := make([]string, 0)
	if t.data != nil {
		for key := range t.data {
			keys = append(keys, fmt.Sprintf("%s:%s", key.Location, key.Key))
		}
	}
	return keys
}
