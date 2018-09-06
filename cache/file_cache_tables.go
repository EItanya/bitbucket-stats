package cache

import (
	"bitbucket/models"
	"errors"
)

type bitbucketProjectsCacheTable struct {
	*fileCacheTableBasic
	data map[fileCacheKey]models.Project
}

func (t *bitbucketProjectsCacheTable) get(key fileCacheKey, dat interface{}) error {
	if t.data == nil {
		err := t.read()
		if err != nil {
			return err
		}
	}
	err, retrievedData := t.fileCacheTableBasic.get(key, t.data)
	return nil
}

func (t *bitbucketProjectsCacheTable) set(key fileCacheKey, dat interface{}) error {
	if t.data == nil {
		err := t.read()
		if err != nil {
			return err
		}
	}
	switch typedData := dat.(type) {
	case models.Project:
		t.dataMtx.Lock()
		t.data[key] = typedData
		t.dataMtx.Unlock()
		err := t.write()
		if err != nil {
			return err
		}
	default:
		return errors.New("Incorrect data type passed into Cache set for projects table")
	}
	return nil
}

func (t *bitbucketProjectsCacheTable) read() error {
	// if t.data == nil {
	t.data = make(map[fileCacheKey]models.Project)
	// }
	err := t.fileCacheTableBasic.read(&t.data)
	if err != nil {
		return err
	}
	return nil
}

func (t *bitbucketProjectsCacheTable) write() error {
	err := t.fileCacheTableBasic.write(t.data)
	if err != nil {
		return err
	}
	return nil
}
