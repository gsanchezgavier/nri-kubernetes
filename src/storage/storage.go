package storage

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"
	"time"
)

const (
	fileExt  = ".json"
	filePerm = 0644
)

var now = time.Now

// Storage defines the interface of a Key-Value storage system, which is able to store the timestamp
// where the key was stored.
type Storage interface {
	// Write stores a value for a given key. Implementors must save also the time when it was stored.
	// The value can be any type.
	Write(key string, value interface{}) error
	// Read gets the value associated to a given key and stores in the value referenced by the pointer passed as argument.
	// It returns the Unix timestamp when the value was stored (in seconds), or an error if the Read operation failed.
	// It may return any type of value.
	Read(key string, valuePtr interface{}) (int64, error)
}

// JSONDiskStorage is a Storage implementation that uses the file system as persistence backend, storing
// the objects as JSON.
// This requires that any object that has to be stored is Marshallable and Unmarshallable.
type JSONDiskStorage struct {
	rootPath string
}

// Holder for any entry in the JSON disk storage
type jsonEntry struct {
	Timestamp int64
	Value     interface{}
}

// NewJSONDiskStorage returns a JSONDiskStorage using the rootPath argument as root folder for the persistent entities.
func NewJSONDiskStorage(rootPath string) JSONDiskStorage {
	return JSONDiskStorage{rootPath: rootPath}
}

// Write stores a value for a given key. Implementors must save also the time when it was stored.
// This implementation adds a restriction to the key name: it must be a valid file name (without extension).
func (j JSONDiskStorage) Write(key string, value interface{}) error {
	entry := jsonEntry{
		Timestamp: now().Unix(),
		Value:     value,
	}
	bytes, err := json.Marshal(&entry)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filepath.Join(j.rootPath, key+fileExt), bytes, filePerm)
}

// Read gets the value associated to a given key and stores it in the value referenced by the pointer passed as
// second argument
// This implementation adds a restriction to the key name: it must be a valid file name (without extension).
func (j JSONDiskStorage) Read(key string, valuePtr interface{}) (int64, error) {
	bytes, err := ioutil.ReadFile(filepath.Join(j.rootPath, key+fileExt))
	if err != nil {
		return 0, err
	}
	var entry jsonEntry
	entry.Value = valuePtr
	err = json.Unmarshal(bytes, &entry)
	if err != nil {
		return 0, err
	}
	return entry.Timestamp, nil
}
