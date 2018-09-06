package cache

// Cache interface for Caches
type Cache interface {
	write(string) error
	clear() error
	read(string) (string, error)
	initialize() error
}

func SaveToCache(c Cache, key string) error {
	return c.write(key)
}

func ReadFromCache(c Cache, key string) (string, error) {
	return c.read(key)
}

func ClearCache(c Cache) error {
	return c.clear()
}

func InitializeCache(c Cache) error {
	return c.initialize()
}
