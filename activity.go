package activity

import (
	"bytes"
	"encoding/json"
	"reflect"
	"time"
)

type Func func(ctx ActivityContext) error

type Activity struct {
	ID        uint64
	Status    ActivityStatus
	StartTime time.Time
	EndTime   time.Time
	Actions   []Action
	Registry  *Registry
	Storage   Storage
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

func Start(storage Storage, reg *Registry, biz int) *Activity {
	return &Activity{
		ID:        1,
		Status:    ActivityStarted,
		StartTime: time.Now(),
		Actions:   []Action{},
		Registry:  reg,
		Storage:   storage,
	}
}

func (a *Activity) Then(doFunc Func, args ...interface{}) func(backFunc Func, args ...interface{}) *Activity {
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
	return func(backFunc Func, args ...interface{}) *Activity {
		var backParams []reflect.Value
		for _, arg := range args {
			backParams = append(backParams, reflect.ValueOf(arg))
		}
		newAction.RollbackFunc = reflect.ValueOf(backFunc)
		newAction.RollbackParams = backParams
		return a
	}
}

func (a *Activity) Exec(ctx ActivityContext) {
	carg := reflect.ValueOf(ctx)
	for _, action := range a.Actions {
		action.DoFunc.Call([]reflect.Value{carg})
	}
}

func (a *Activity) Run(ctx ActivityContext) error {
	err := a.SaveLog()
	if err != nil {
		return err
	}
	a.Exec(ctx)
	return nil
}

func (a *Activity) SaveLog() error {
	ar := a.activeToRecord()
	err := a.Storage.saveActivityRecord(ar.ActivityID, toJson(ar))
	if err != nil {
		return err
	}
	ars := a.actionsToRecord()
	rs := make([]actionData, len(ars))
	for _, ar := range ars {
		rs = append(rs, actionData{
			actionID: ar.ActionID,
			data:     toJson(ar),
		})
	}
	err = a.Storage.saveActionRecord(rs)
	if err != nil {
		return err
	}
	return nil
}

func (a *Activity) activeToRecord() *ActivityRecord {
	r := &ActivityRecord{
		ActivityID: a.ID,
		Status:     a.Status,
		StartTime:  a.StartTime,
		EndTime:    a.EndTime,
	}
	return r
}

func (a *Activity) actionsToRecord() []ActionRecord {
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

func (a *Activity) valueArrayToString(values []reflect.Value) string {
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
	return string(s)
}
