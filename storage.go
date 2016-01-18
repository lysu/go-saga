package saga

type Storage interface {
	saveActivityLog(activityID uint64, data string) error
	saveActionLogs(actionRecords []actionData) error
}

type MemStorage struct {
	data map[string]string
}

func NewMemStorage() (Storage, error) {
	return &MemStorage{
		data: make(map[string]string),
	}, nil
}

func (s *MemStorage) saveActivityLog(activityID uint64, data string) error {
	return nil
}

func (s *MemStorage) saveActionLogs(actionDatas []actionData) error {
	return nil
}
