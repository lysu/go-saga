package activity

import "database/sql"

type DBStorage struct {
	db *sql.DB
}

func NewStorage(host string) (*DBStorage, error) {
	d, err := sql.Open("mysql", host)
	if err != nil {
		return nil, err
	}
	return &DBStorage{d}, nil
}

func (s *DBStorage) saveActivityRecord(r *ActivityRecord) error {
	return nil
}

func (s *DBStorage) saveActionRecord(rs []ActionRecord) error {
	return nil
}
