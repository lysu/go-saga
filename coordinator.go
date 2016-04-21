package saga

import (
	"fmt"
	"github.com/juju/errors"
	"golang.org/x/net/context"
	"reflect"
	"strconv"
)

// ExecutionCoordinator presents Saga Execution Coordinator.
// It manages:
// - Saga log storage.
// - Sub-transaction definition with it's parameter info.
type ExecutionCoordinator struct {
	subTxDefinitions  subTxDefinitions
	paramTypeRegister *paramTypeRegister
}

// NewSEC creates Saga Execution Coordinator
// This method require supply a log Storage to save & lookup log during tx execute.
func NewSEC() ExecutionCoordinator {
	return ExecutionCoordinator{
		subTxDefinitions: make(subTxDefinitions),
		paramTypeRegister: &paramTypeRegister{
			nameToType: make(map[string]reflect.Type),
			typeToName: make(map[reflect.Type]string),
		},
	}
}

// AddSubTxDef create & add definition base on given subTxID, action and compensate, and return current SEC.
//
// subTxID identifies a sub-transaction type, it also be use to persist into saga-log and be lookup for retry
// action defines the action that sub-transaction will execute.
// compensate defines the compensate that sub-transaction will execute when sage aborted.
//
// action and compensate MUST a function that context.Context as first argument.
func (e *ExecutionCoordinator) AddSubTxDef(subTxID string, action interface{}, compensate interface{}) *ExecutionCoordinator {
	e.paramTypeRegister.addParams(action)
	e.paramTypeRegister.addParams(compensate)
	e.subTxDefinitions.addDefinition(subTxID, action, compensate)
	return e
}

// MustFindSubTxDef returns sub transaction definition by given subTxID.
// Panic if not found sub-transaction.
func (e *ExecutionCoordinator) MustFindSubTxDef(subTxID string) subTxDefinition {
	define, ok := e.subTxDefinitions.findDefinition(subTxID)
	if !ok {
		panic("SubTxID: " + subTxID + " not found in context")
	}
	return define
}

// MustFindParamName return param name by given reflect type.
// Panic if param name not found.
func (e *ExecutionCoordinator) MustFindParamName(typ reflect.Type) string {
	name, ok := e.paramTypeRegister.findTypeName(typ)
	if !ok {
		panic("Find Param Name Panic: " + name)
	}
	return name
}

// MustFindParamType return param type by given name.
// Panic if param type not found.
func (e *ExecutionCoordinator) MustFindParamType(name string) reflect.Type {
	typ, ok := e.paramTypeRegister.findType(name)
	if !ok {
		panic("Find Param Type Panic: " + name)
	}
	return typ
}

func (e *ExecutionCoordinator) StartCoordinator() error {
	logIDs, err := LogStorage().LogIDs()
	if err != nil {
		return errors.Annotate(err, "Fetch logs failure")
	}
	for _, logID := range logIDs {
		lastLogData, err := LogStorage().LastLog(logID)
		if err != nil {
			return errors.Annotate(err, "Fetch last log panic")
		}
		fmt.Println(lastLogData)
	}
	return nil
}

// StartSaga start a new saga, returns the saga was started.
// This method need execute context and UNIQUE id to identify saga instance.
func (e *ExecutionCoordinator) StartSaga(ctx context.Context, id uint64) *Saga {
	s := &Saga{
		id:      id,
		context: ctx,
		sec:     e,
		logID:   LogPrefix + strconv.FormatInt(int64(id), 10),
	}
	s.startSaga()
	return s
}
