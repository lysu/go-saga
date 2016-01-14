package saga

type Storage interface {
	saveActivityRecord(activityID uint64, data string) error
	saveActionRecord(actionRecords []actionData) error
}

type MemStorage struct {
	data map[string]string
}

func NewMemStorage() (Storage, error) {
	return &MemStorage{
		data: make(map[string]string),
	}, nil
}

func (s *MemStorage) saveActivityRecord(activityID uint64, data string) error {
	return nil
}

func (s *MemStorage) saveActionRecord(actionRecords []actionData) error {
	return nil
}
