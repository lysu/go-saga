package saga

import "reflect"

// SubTxDefinitions maintains SubTx that use in current application.
// You MUST init it as singleton, and register SubTxDefinition into it
// before use other Saga function
type SubTxDefinitions map[string]SubTxDefinition

// SubTxDefinition defines sub-transaction detail.
type SubTxDefinition struct {

	// SubTxID identifies a sub-transaction type.
	// it also be use to persist into saga-log and be lookup
	// when transaction retry or recovery
	SubTxID string

	// Action defines the action that sub-transaction will execute.
	// it will be the reflect.Value of a function
	Action reflect.Value

	// Action defines the compensate that sub-transaction will execute when sage aborted.
	// it will be the reflect.Value of a function
	Compensate reflect.Value
}

// AddDefinition create definition on the given subTxID, action and compensate
// then add it into SubTxDefinitions, and return definitions.
// Action and compensate MUST a function that SagaContext as first argument.
func (s SubTxDefinitions) AddDefinition(subTxID string, action interface{}, compensate interface{}) SubTxDefinitions {
	actionMethod := subTxMethod(action)
	compensateMethod := subTxMethod(compensate)
	s[subTxID] = SubTxDefinition{
		SubTxID:    subTxID,
		Action:     actionMethod,
		Compensate: compensateMethod,
	}
	return s
}

// FindDefinition returns definition by given subTxID and whether definition found.
func (s SubTxDefinitions) FindDefinition(subTxID string) (SubTxDefinition, bool) {
	define, ok := s[subTxID]
	return define, ok
}

func subTxMethod(obj interface{}) reflect.Value {
	funcValue := reflect.ValueOf(obj)
	if funcValue.Kind() != reflect.Func {
		panic("Regist object must be a func")
	}
	if funcValue.Type().NumIn() < 1 ||
		funcValue.Type().In(0) != sagaContextType {
		panic("First argument must use SagaContext.")
	}
	return funcValue
}
