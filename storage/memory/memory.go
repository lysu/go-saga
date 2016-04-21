package memory

import (
	"github.com/juju/errors"
	"github.com/lysu/go-saga"
	"github.com/lysu/go-saga/storage"
	"sync"
)

var storageInstance storage.Storage
var memoryInit sync.Once

func init() {
	saga.StorageProvider = func(cfg storage.StorageConfig) storage.Storage {
		memoryInit.Do(func() {
			var err error
			storageInstance, err = newMemStorage()
			if err != nil {
				panic(err)
			}
		})
		return storageInstance
	}
}

type memStorage struct {
	data map[string][]string
}

// NewMemStorage creates log storage base on memory.
// This storage use simple `map[string][]string`, just for TestCase used.
// NOT use this in product.
func newMemStorage() (storage.Storage, error) {
	return &memStorage{
		data: make(map[string][]string),
	}, nil
}

// AppendLog appends log into queue under given logID.
func (s *memStorage) AppendLog(logID string, data string) error {
	logQueue, ok := s.data[logID]
	if !ok {
		logQueue = []string{}
		s.data[logID] = logQueue
	}
	s.data[logID] = append(s.data[logID], data)
	return nil
}

// Lookup lookups log under given logID.
func (s *memStorage) Lookup(logID string) ([]string, error) {
	return s.data[logID], nil
}

// Close uses to close storage and release resources.
func (s *memStorage) Close() error {
	return nil
}

// LogIDs uses to take all Log ID av in current storage
func (s *memStorage) LogIDs() ([]string, error) {
	ids := make([]string, 0, len(s.data))
	for id := range s.data {
		ids = append(ids, id)
	}
	return ids, nil
}

func (s *memStorage) Cleanup(logID string) error {
	delete(s.data, logID)
	return nil
}

func (s *memStorage) LastLog(logID string) (string, error) {
	logData, ok := s.data[logID]
	if !ok {
		err := errors.NewErr("LogData %s not found", logID)
		return "", &err
	}
	sizeOfLog := len(logData)
	if sizeOfLog == 0 {
		return "", errors.New("LogData is empty")
	}
	lastLog := logData[sizeOfLog-1]
	return lastLog, nil
}
