package saga

import (
	"golang.org/x/net/context"
	"reflect"
	"time"
)

const (
	LogPrefix = "saga_"
)

type Saga struct {
	id          uint64
	logID       string
	context     context.Context
	coordinator *ExecutionCoordinator
}

func (s *Saga) StartSaga() {
	log := &Log{
		Type: SagaStart,
		Time: time.Now(),
	}
	s.coordinator.LogStorage.AppendLog(s.logID, log.MustMarshal())
}

func (s *Saga) SubTx(subTxID string, args ...interface{}) *Saga {
	subTxDef := s.coordinator.FindSubTxDef(subTxID)
	log := &Log{
		Type:    ActionStart,
		SubTxID: subTxID,
		Time:    time.Now(),
		Params:  s.MarshalParam(args),
	}
	s.coordinator.LogStorage.AppendLog(s.logID, log.MustMarshal())

	params := make([]reflect.Value, 0, len(args)+1)
	params = append(params, reflect.ValueOf(s.context))
	for _, arg := range args {
		params = append(params, reflect.ValueOf(arg))
	}
	result := subTxDef.Action.Call(params)
	if isReturnError(result) {
		s.Abort()
		return s
	}

	log = &Log{
		Type:    ActionEnd,
		SubTxID: subTxID,
		Time:    time.Now(),
	}
	s.coordinator.LogStorage.AppendLog(s.logID, log.MustMarshal())
	return s
}

func (s *Saga) EndSaga() {
	log := &Log{
		Type: SagaEnd,
		Time: time.Now(),
	}
	s.coordinator.LogStorage.AppendLog(s.logID, log.MustMarshal())
}

func (s *Saga) Abort() {
	logs, err := s.coordinator.LogStorage.Lookup(s.logID)
	if err != nil {
		panic("Abort Panic")
	}
	for i := len(logs) - 1; i >= 0; i-- {
		logData := logs[i]
		log := MustUnmarshal(logData)
		if log.Type == ActionStart {
			if err := s.Compensate(log); err != nil {
				panic("Compensate Failure..")
			}
		}
	}
}

func (s *Saga) Compensate(tlog Log) error {
	clog := &Log{
		Type:    CompensateStart,
		SubTxID: tlog.SubTxID,
		Time:    time.Now(),
	}
	s.coordinator.LogStorage.AppendLog(s.logID, clog.MustMarshal())

	args := s.UnmarshalParam(tlog.Params)

	params := make([]reflect.Value, 0, len(args)+1)
	params = append(params, reflect.ValueOf(s.context))
	for _, arg := range args {
		params = append(params, arg)
	}

	subDef := s.coordinator.FindSubTxDef(tlog.SubTxID)
	result := subDef.Compensate.Call(params)
	if isReturnError(result) {
		s.Abort()
	}

	clog = &Log{
		Type:    CompensateEnd,
		SubTxID: tlog.SubTxID,
		Time:    time.Now(),
	}
	s.coordinator.LogStorage.AppendLog(s.logID, clog.MustMarshal())
	return nil
}

func isReturnError(result []reflect.Value) bool {
	if len(result) == 1 && !result[0].IsNil() {
		return true
	}
	return false
}
