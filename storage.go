package saga

type ConsumeOffset int

type Storage interface {
	AppendLog(logID string, data string) error
	Lookup(logID string) ([]string, error)
}

type MemStorage struct {
	data map[string][]string
}

func NewMemStorage() (Storage, error) {
	return &MemStorage{
		data: make(map[string][]string),
	}, nil
}

func (s *MemStorage) AppendLog(logID string, data string) error {
	logQueue, ok := s.data[logID]
	if !ok {
		logQueue = []string{}
		s.data[logID] = logQueue
	}
	s.data[logID] = append(s.data[logID], data)
	return nil
}

func (s *MemStorage) Lookup(logID string) ([]string, error) {
	return s.data[logID], nil
}
