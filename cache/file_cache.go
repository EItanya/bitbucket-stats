package cache

type FileCache struct {
	// RawFileData       *api.SavedFiles
	// FileDataByRepo    []*api.RepoModel
	// FileDataByProject []*api.ProjectModel
}
type FileConfig struct {
	Dir string
}

func (c *FileCache) write(key string) error {
	return nil
}

func (c *FileCache) read(keys []string) ([]CacheEntity, error) {
	result := make([]CacheEntity, len(keys))
	return result, nil
}

func (c *FileCache) clear() error {
	return nil
}

func (r *FileCache) initialize() error {
	return nil
}
