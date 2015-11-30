package activity

import (
	"bytes"
	"encoding/json"
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
	ar := a.activeToRecord()
	err := a.Storage.saveActivityRecord(&ar)
	if err != nil {
		return err
	}
	ars := a.actionsToRecord()
	err = a.Storage.saveActionRecord(ars)
	if err != nil {
		return err
	}
	return nil
}

func (a *Activitor) activeToRecord() ActivityRecord {
	r := ActivityRecord{
		ID:        a.ID,
		Status:    a.Status,
		StartTime: a.StartTime,
		EndTime:   a.EndTime,
	}
	return r
}

func (a *Activitor) actionsToRecord() []ActionRecord {
	registry := a.Registry
	var rs []ActionRecord
	for _, action := range a.Actions {
		rs = append(rs, ActionRecord{
			Status:         action.Status,
			StartTime:      action.StartTime,
			EndTime:        action.EndTime,
			ActivityID:     a.ID,
			DoFuncID:       registry.FindFuncID(action.DoFunc),
			DoParams:       a.valueArrayToString(action.DoParams),
			RollbackFuncID: registry.FindFuncID(action.RollbackFunc),
			RollbackParams: a.valueArrayToString(action.RollbackParams),
		})
	}
	return rs
}

func (a *Activitor) valueArrayToString(values []reflect.Value) string {
	var buf bytes.Buffer
	for i, value := range values {
		if i != 0 {
			buf.WriteString(",")
		}
		typ := a.Registry.FindTypeName(value.Type())
		buf.WriteString(typ)
		buf.WriteString(":")
		buf.WriteString(toJson(value.Interface()))
	}
	return buf.String()
}

func toJson(value interface{}) string {
	s, err := json.Marshal(value)
	if err != nil {
		panic(err)
	}
	return s
}
