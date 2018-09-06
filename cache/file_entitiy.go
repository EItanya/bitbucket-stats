package cache

import (
	"bitbucket/models"
	"errors"
)

// Type of entity saved to File Cache
type FileEntity struct {
	RawData interface{}
}

func (fe *FileEntity) unmarshal(dat interface{}) error {
	switch typedData := dat.(type) {
	case *models.Project:
		*typedData = fe.RawData.(models.Project)
		// *dat.(*models.Project) = fe.RawData.(models.Project)
	case *models.Repository:
		*dat.(*models.Repository) = fe.RawData.(models.Repository)
	case *models.FilesID:
		*dat.(*models.FilesID) = fe.RawData.(models.FilesID)
	// case *interface{}:
	// 	e := (*interface{})(unsafe.Pointer(typedData))
	// 	*e = fe.RawData
	default:
		return errors.New("Must pass pointer to interface into unmarshal method")
	}
	return nil
}

func (fe *FileEntity) marshal(dat interface{}) error {
	switch typedData := dat.(type) {
	case models.Repository:
		fe.RawData = typedData
	case models.Project:
		fe.RawData = typedData
	case models.FilesID:
		fe.RawData = typedData
	case models.Files:
		fe.RawData = typedData
	case []string:
		fe.RawData = typedData
	case int:
		fe.RawData = typedData
	case float64:
		fe.RawData = typedData
	case string:
		fe.RawData = typedData
	default:
		return errors.New("Data passed to Redis Entity was not in correct form")
	}
	return nil
}
