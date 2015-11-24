package activity

import (
	"reflect"
	"time"
)

type Func func(ctx ActivityContext) error

type Activitor struct {
	ID        uint64
	Status    ActivityStatus
	StartTime time.Time
	EndTime   time.Time
	Actions   []Action
	Registry  *Registry
	Storage   *DBStorage
}

type Action struct {
	Status         ActionStatus
	StartTime      time.Time
	EndTime        time.Time
	DoFunc         reflect.Value
	DoParams       []reflect.Value
	RollbackFunc   reflect.Value
	RollbackParams []reflect.Value
}

func Start(storage *DBStorage, reg *Registry, biz int) *Activitor {
	return &Activitor{
		ID:        1,
		Status:    ActivityStarted,
		StartTime: time.Now(),
		Actions:   []Action{},
		Registry:  reg,
		Storage:   storage,
	}
}

func (a *Activitor) Then(doFunc Func, args ...interface{}) func(backFunc Func, args ...interface{}) *Activitor {
	var doParams []reflect.Value
	for _, arg := range args {
		doParams = append(doParams, reflect.ValueOf(arg))
	}
	newAction := Action{
		Status:    ActionStarted,
		StartTime: time.Now(),
		DoFunc:    reflect.ValueOf(doFunc),
		DoParams:  doParams,
	}
	a.Actions = append(a.Actions, newAction)
	return func(backFunc Func, args ...interface{}) *Activitor {
		var backParams []reflect.Value
		for _, arg := range args {
			backParams = append(backParams, reflect.ValueOf(arg))
		}
		newAction.RollbackFunc = reflect.ValueOf(backFunc)
		newAction.RollbackParams = backParams
		return a
	}
}

func (a *Activitor) Run(ctx ActivityContext) error {
	err := a.SaveLog()
	if err != nil {
		return err
	}
	return nil
}

func (a *Activitor) SaveLog() error {
	ar := activeToRecord(a)
	err := a.Storage.saveActivityRecord(&ar)
	if err != nil {
		return err
	}
	ars := actionsToRecord(a)
	err = a.Storage.saveActionRecord(ars)
	if err != nil {
		return err
	}
	return nil
}

func activeToRecord(a *Activitor) ActivityRecord {
	r := ActivityRecord{
		ID:        a.ID,
		Status:    a.Status,
		StartTime: a.StartTime,
		EndTime:   a.EndTime,
	}
	return r
}

func actionsToRecord(a *Activitor) []ActionRecord {
	registry := a.Registry
	var rs []ActionRecord
	for _, action := range a.Actions {
		rs = append(rs, ActionRecord{
			Status:         action.Status,
			StartTime:      action.StartTime,
			EndTime:        action.EndTime,
			ActivityID:     a.ID,
			DoFuncID:       registry.FindFuncID(action.DoFunc),
			DoParams:       valueArrayToString(action.DoParams),
			RollbackFuncID: registry.FindFuncID(action.RollbackFunc),
			RollbackParams: valueArrayToString(action.RollbackParams),
		})
	}
	return rs
}

func valueArrayToString(values []reflect.Value) string {
	for _, value := range values {
		value.Type()
	}
	return ""
}
