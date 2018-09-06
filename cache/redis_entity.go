package cache

import (
	"bitbucket/models"
	"errors"
)

type RedisEntity struct {
	RawData interface{}
}

// Retrieve data
func (re *RedisEntity) unmarshal(dat interface{}) error {
	switch dat.(type) {
	case *interface{}:
		*dat.(*interface{}) = re.RawData
	case *models.Project:
		*dat.(*models.Project) = re.RawData.(models.Project)
	case *models.Repository:
		*dat.(*models.Repository) = re.RawData.(models.Repository)
	case *models.Files:
		*dat.(*models.Files) = re.RawData.(models.Files)
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

type RedisList struct {
	RawData []interface{}
}

func (rl *RedisList) unmarshal(dat interface{}) error {
	switch dat.(type) {
	case *interface{}:
		*dat.(*interface{}) = rl.RawData
	default:
		return errors.New("Must pass pointer to interface into unmarshal method")
	}
	return nil
}

func (rl *RedisList) marshal(dat interface{}) error {
	switch typedData := dat.(type) {
	case []models.Repository:
		convertedArray := make([]interface{}, len(typedData))
		for i, val := range typedData {
			convertedArray[i] = val
		}
		rl.RawData = convertedArray
	case []models.Project:
		convertedArray := make([]interface{}, len(typedData))
		for i, val := range typedData {
			convertedArray[i] = val
		}
		rl.RawData = convertedArray
	default:
		return errors.New("Data passed to Redis List must be in list form")
	}
	return nil
}
