package saga

import (
	"golang.org/x/net/context"
	"reflect"
	"strconv"
)

type ExecutionCoordinator struct {
	SubTxDefinitions  SubTxDefinitions
	ParamTypeRegister *paramTypeRegister
	LogStorage        Storage
}

func (e *ExecutionCoordinator) AddDefinition(subTxID string, action interface{}, compensate interface{}) *ExecutionCoordinator {
	e.ParamTypeRegister.addParams(action)
	e.ParamTypeRegister.addParams(compensate)
	e.SubTxDefinitions.AddDefinition(subTxID, action, compensate)
	return e
}

func (e *ExecutionCoordinator) FindSubTxDef(subTxID string) SubTxDefinition {
	define, ok := e.SubTxDefinitions.FindDefinition(subTxID)
	if !ok {
		panic("SubTxID: " + subTxID + " not found in context")
	}
	return define
}

func (e *ExecutionCoordinator) FindParamName(typ reflect.Type) string {
	return e.ParamTypeRegister.FindTypeName(typ)
}

func (e *ExecutionCoordinator) FindParamType(name string) reflect.Type {
	typ, ok := e.ParamTypeRegister.FindType(name)
	if !ok {
		panic("Find Param Type Panic: " + name)
	}
	return typ
}

func (e *ExecutionCoordinator) StartSaga(ctx context.Context, id uint64) *Saga {
	s := &Saga{
		id:          id,
		context:     ctx,
		coordinator: e,
		logID:       LogPrefix + strconv.FormatInt(int64(id), 10),
	}
	s.StartSaga()
	return s
}
