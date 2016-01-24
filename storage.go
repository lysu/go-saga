package saga

// Storage uses to support save and lookup saga log.
type Storage interface {

	// AppendLog appends log data into log under given logID
	AppendLog(logID string, data string) error

	// Lookup uses to lookup all log under given logID
	Lookup(logID string) ([]string, error)
}

type memStorage struct {
	data map[string][]string
}

// NewMemStorage creates log storage base on memory.
// This storage use simple `map[string][]string`, just for TestCase used.
// NOT use this in product.
func NewMemStorage() (Storage, error) {
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
