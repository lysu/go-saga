package saga

import (
	"golang.org/x/net/context"
	"reflect"
	"time"
)

const logPrefix = "saga_"

// Saga presents current execute transaction.
// A Saga constituted by small sub-transactions.
type Saga struct {
	id      uint64
	logID   string
	context context.Context
	sec     *ExecutionCoordinator
}

func (s *Saga) startSaga() {
	log := &Log{
		Type: SagaStart,
		Time: time.Now(),
	}
	err := s.sec.logStorage.AppendLog(s.logID, log.mustMarshal())
	if err != nil {
		panic("Add log Failure")
	}
}

// SubTx executes a sub-transaction for given subTxID(which define in SEC initialize) and arguments.
// it returns current Saga.
func (s *Saga) SubTx(subTxID string, args ...interface{}) *Saga {
	subTxDef := s.sec.MustFindSubTxDef(subTxID)
	log := &Log{
		Type:    ActionStart,
		SubTxID: subTxID,
		Time:    time.Now(),
		Params:  MarshalParam(s.sec, args),
	}
	err := s.sec.logStorage.AppendLog(s.logID, log.mustMarshal())
	if err != nil {
		panic("Add log Failure")
	}

	params := make([]reflect.Value, 0, len(args)+1)
	params = append(params, reflect.ValueOf(s.context))
	for _, arg := range args {
		params = append(params, reflect.ValueOf(arg))
	}
	result := subTxDef.action.Call(params)
	if isReturnError(result) {
		s.Abort()
		return s
	}

	log = &Log{
		Type:    ActionEnd,
		SubTxID: subTxID,
		Time:    time.Now(),
	}
	err = s.sec.logStorage.AppendLog(s.logID, log.mustMarshal())
	if err != nil {
		panic("Add log Failure")
	}
	return s
}

// EndSaga finishes a Saga's execution.
func (s *Saga) EndSaga() {
	log := &Log{
		Type: SagaEnd,
		Time: time.Now(),
	}
	err := s.sec.logStorage.AppendLog(s.logID, log.mustMarshal())
	if err != nil {
		panic("Add log Failure")
	}
}

// Abort stop and compensate to rollback to start situation.
// This method will stop continue sub-transaction and do Compensate for executed sub-transaction.
// SubTx will call this method internal.
func (s *Saga) Abort() {
	logs, err := s.sec.logStorage.Lookup(s.logID)
	if err != nil {
		panic("Abort Panic")
	}
	for i := len(logs) - 1; i >= 0; i-- {
		logData := logs[i]
		log := mustUnmarshalLog(logData)
		if log.Type == ActionStart {
			if err := s.compensate(log); err != nil {
				panic("Compensate Failure..")
			}
		}
	}
}

func (s *Saga) compensate(tlog Log) error {
	clog := &Log{
		Type:    CompensateStart,
		SubTxID: tlog.SubTxID,
		Time:    time.Now(),
	}
	err := s.sec.logStorage.AppendLog(s.logID, clog.mustMarshal())
	if err != nil {
		panic("Add log Failure")
	}

	args := UnmarshalParam(s.sec, tlog.Params)

	params := make([]reflect.Value, 0, len(args)+1)
	params = append(params, reflect.ValueOf(s.context))
	for _, arg := range args {
		params = append(params, arg)
	}

	subDef := s.sec.MustFindSubTxDef(tlog.SubTxID)
	result := subDef.compensate.Call(params)
	if isReturnError(result) {
		s.Abort()
	}

	clog = &Log{
		Type:    CompensateEnd,
		SubTxID: tlog.SubTxID,
		Time:    time.Now(),
	}
	err = s.sec.logStorage.AppendLog(s.logID, clog.mustMarshal())
	if err != nil {
		panic("Add log Failure")
	}
	return nil
}

func isReturnError(result []reflect.Value) bool {
	if len(result) == 1 && !result[0].IsNil() {
		return true
	}
	return false
}
