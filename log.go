package saga

import (
	"encoding/json"
	"time"
)

// LogType present type flag for Log
type LogType int

const (
	// SagaStart flag saga stared log
	SagaStart LogType = iota + 1
	// SagaEnd flag saga ended log
	SagaEnd
	// SagaAbort flag saga aborted
	SagaAbort
	// ActionStart flag action start log
	ActionStart
	// ActionEnd flag action end log
	ActionEnd
	// CompensateStart flag compensate start log
	CompensateStart
	// CompensateEnd flag compensate end log
	CompensateEnd
)

// Log presents Saga Log.
// Saga Log used to log execute status for saga,
// and SEC use it to compensate and retry.
type Log struct {
	Type    LogType     `json:"type,omitempty"`
	SubTxID string      `json:"subTxID,omitempty"`
	Time    time.Time   `json:"time,omitempty"`
	Params  []ParamData `json:"params,omitempty"`
}

func (l *Log) mustMarshal() string {
	return mustMarshal(l)
}

func mustUnmarshalLog(data string) Log {
	var log Log
	mustUnmarshal([]byte(data), &log)
	return log
}

func mustMarshal(value interface{}) string {
	s, err := json.Marshal(value)
	if err != nil {
		panic("Marshal Failure")
	}
	return string(s)
}

func mustUnmarshal(data []byte, v interface{}) {
	err := json.Unmarshal([]byte(data), v)
	if err != nil {
		panic("Unmarshal Failure")
	}
}
