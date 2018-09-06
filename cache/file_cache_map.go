package cache

import (
	"bitbucket/models"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"sync"
	"time"
)

var projectsCacheTable *fileCacheMap
var repositoriesCacheTable *fileCacheMap
var filesCacheTable *fileCacheMap

var fileCacheMaps []*fileCacheMap

func init() {
	projectsCacheTable = new(fileCacheMap)
	projectsCacheTable.initialize(defaultDir, projectConst)
	projectsCacheTable.read()
	projectsCacheTable.unmarshalEntity = func(c CacheEntity) (interface{}, error) {
		dat := models.Project{}
		err := c.Unmarshal(&dat)
		return dat, err
	}
	projectsCacheTable.marshalEntity = func(dat interface{}) (CacheEntity, error) {
		result := &models.Project{}
		err := result.Marshal(dat)
		return result, err
	}

	repositoriesCacheTable = new(fileCacheMap)
	repositoriesCacheTable.initialize(defaultDir, repositoryConst)
	repositoriesCacheTable.read()
	repositoriesCacheTable.unmarshalEntity = func(c CacheEntity) (interface{}, error) {
		dat := models.Repository{}
		err := c.Unmarshal(&dat)
		return dat, err
	}
	repositoriesCacheTable.marshalEntity = func(dat interface{}) (CacheEntity, error) {
		result := &models.Repository{}
		err := result.Marshal(dat)
		return result, err
	}

	filesCacheTable = new(fileCacheMap)
	filesCacheTable.initialize(defaultDir, filesConst)
	filesCacheTable.read()
	filesCacheTable.unmarshalEntity = func(c CacheEntity) (interface{}, error) {
		dat := models.FilesID{}
		err := c.Unmarshal(&dat)
		return dat, err
	}
	filesCacheTable.marshalEntity = func(dat interface{}) (CacheEntity, error) {
		result := &models.FilesID{}
		err := result.Marshal(dat)
		return result, err
	}

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
	data            fileCacheMapData
	marshalEntity   func(interface{}) (CacheEntity, error)
	unmarshalEntity func(CacheEntity) (interface{}, error)
}
type fileCacheMapData map[fileCacheKey]interface{}

func (t *fileCacheMap) write() error {
	if t.filename == "" {
		return errors.New("Filename of Cache Table cannot be null")
	}
	if t.data == nil {
		return errors.New("Data object to be written cannot be null")
	}
	t.fileMtx.Lock()
	dataToWrite := make(writtenCacheMapData, 0)
	dataToWrite.prepareDataForWrite(t.data)
	byt, err := json.Marshal(dataToWrite)
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

func (t *fileCacheMap) read() error {
	if t.filename == "" {
		return errors.New("Filename of Cache Table cannot be null")
	}
	if info, err := os.Stat(t.filename); err == nil {
		if time.Since(t.lastSync) > time.Since(info.ModTime()) {
			t.fileMtx.RLock()
			byt, err := ioutil.ReadFile(t.filename)
			if err != nil {
				return err
			}
			// var readData interface{}
			var readData writtenCacheMapData
			err = json.Unmarshal(byt, &readData)
			t.data = readData.translateDataFromRead()
			if err != nil {
				return err
			}
			t.fileMtx.RUnlock()
		}
		return nil
	}
	return fmt.Errorf("File (%s) could not be found for the (%s) file cache table", t.filename, t.title)
}

func (t *fileCacheMap) get(key fileCacheKey) (CacheEntity, error) {
	t.dataMtx.RLock()
	data, err := t.marshalEntity(t.data[key])
	t.dataMtx.RUnlock()
	return data, err
}

func (t *fileCacheMap) set(key fileCacheKey, dataToSave CacheEntity) error {
	t.dataMtx.Lock()
	data, err := t.unmarshalEntity(dataToSave)
	if err != nil {
		return err
	}
	t.data[key] = data
	t.dataMtx.Unlock()
	return nil
}

func (t *fileCacheMap) clear() error {
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

type writtenCacheMapData map[string]interface{}

func (t *writtenCacheMapData) prepareDataForWrite(data fileCacheMapData) {
	for key, value := range data {
		stringKey := fmt.Sprintf("%s:%s", key.Location, key.Key)
		(*t)[stringKey] = value
	}
}

func (t *writtenCacheMapData) translateDataFromRead() fileCacheMapData {
	result := make(fileCacheMapData)
	for key, val := range *t {
		splitKey := strings.SplitN(key, ":", 2)
		cacheKey := fileCacheKey{
			Key:      splitKey[1],
			Location: splitKey[0],
		}
		result[cacheKey] = val
	}
	return result
}
