package cache

import (
	"encoding/json"
	"os"
)

type FileCache struct {
	Path string
	data map[string]string
}

func NewFileCache(path string) (*FileCache, error) {
	c := &FileCache{Path: path, data: make(map[string]string)}

	if _, err := os.Stat(path); err == nil {
		bytes, _ := os.ReadFile(path)
		json.Unmarshal(bytes, &c.data)
	}
	return c, nil
}

func (c *FileCache) Get(key string) (string, bool) {
	val, ok := c.data[key]
	return val, ok
}

func (c *FileCache) Set(key string, val string) error {
	c.data[key] = val
	bytes, _ := json.MarshalIndent(c.data, "", "  ")
	return os.WriteFile(c.Path, bytes, 0644)
}
