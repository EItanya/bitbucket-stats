package cache

import (
	"bitbucket/models"
	"errors"
	"unsafe"
)

type RedisEntity struct {
	RawData interface{}
}

// Retrieve data
func (re *RedisEntity) unmarshal(dat interface{}) error {
	switch typedData := dat.(type) {
	case *models.Project:
		*typedData = re.RawData.(models.Project)
		// *dat.(*models.Project) = re.RawData.(models.Project)
	case *models.Repository:
		*dat.(*models.Repository) = re.RawData.(models.Repository)
	case *models.FilesID:
		*dat.(*models.FilesID) = re.RawData.(models.FilesID)
	case *interface{}:
		e := (*interface{})(unsafe.Pointer(typedData))
		*e = re.RawData
	default:
		return errors.New("Must pass pointer to interface into unmarshal method")
	}
	return nil
}

func (re *RedisEntity) marshal(dat interface{}) error {
	switch typedData := dat.(type) {
	case models.Repository:
		re.RawData = typedData
	case models.Project:
		re.RawData = typedData
	case models.FilesID:
		re.RawData = typedData
	case int:
		re.RawData = typedData
	case float64:
		re.RawData = typedData
	case string:
		re.RawData = typedData
	default:
		return errors.New("Data passed to Redis Entity was not in correct form")
	}
	return nil
}
