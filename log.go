package saga

import (
	"encoding/json"
	"time"
)

type LogType int

const (
	SagaStart LogType = iota + 1
	SagaEnd
	ActionStart
	ActionEnd
	CompensateStart
	CompensateEnd
)

type Log struct {
	Type    LogType     `json:"type,omitempty"`
	SubTxID string      `json:"subTxID,omitempty"`
	Time    time.Time   `json:"time,omitempty"`
	Params  []ParamData `json:"params,omitempty"`
}

func (l *Log) MustMarshal() string {
	return mustMarshal(l)
}

func MustUnmarshal(data string) Log {
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
