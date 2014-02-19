package spinner

import (
	"io/ioutil"
	"path/filepath"
	"sync"
)

type cachedFileContent struct {
	mutex sync.Mutex
	files map[string][]byte
}

func (c *cachedFileContent) ReadFile(filename string) ([]byte, error) {
	if filepath.IsAbs(filename) {
		fp := filepath.Clean(filename)
		c.mutex.Lock()
		content, ok := c.files[fp]
		c.mutex.Unlock()
		if ok {
			return content, nil
		}
		content, err := ioutil.ReadFile(fp)
		if err == nil {
			c.mutex.Lock()
			c.files[fp] = content
			c.mutex.Unlock()
		}
		return content, err
	}
	return ioutil.ReadFile(filename)
}

var readOnce = &cachedFileContent{
	files: make(map[string][]byte),
}

func ReadOnce(filename string) ([]byte, error) {
	return readOnce.ReadFile(filename)
}
