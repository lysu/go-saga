package activity

import (
	"bytes"
	"encoding/json"
	"reflect"
	"time"
)

type Activity struct {
	ID        uint64
	Status    ActivityStatus
	StartTime time.Time
	EndTime   time.Time
	Actions   []*Action
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
		Actions:   []*Action{},
		Registry:  reg,
		Storage:   storage,
	}
}

func (a *Activity) Then(doFunc interface{}, args ...interface{}) func(backFunc interface{}, args ...interface{}) *Activity {
	var doParams []reflect.Value
	for _, arg := range args {
		doParams = append(doParams, reflect.ValueOf(arg))
	}
	doFuncValue := reflect.ValueOf(doFunc)
	if doFuncValue.Kind() != reflect.Func {
		panic("Regist object must be a func")
	}
	if doFuncValue.Type().NumIn() < 1 ||
		doFuncValue.Type().NumIn() != len(doParams)+1 ||
		doFuncValue.Type().In(0) != activityContextType {
		panic("First argument must use ActivityContext.")
	}
	newAction := &Action{
		Status:    ActionStarted,
		StartTime: time.Now(),
		DoFunc:    doFuncValue,
		DoParams:  doParams,
	}
	a.Actions = append(a.Actions, newAction)
	return func(backFunc interface{}, args ...interface{}) *Activity {
		var backParams []reflect.Value
		for _, arg := range args {
			backParams = append(backParams, reflect.ValueOf(arg))
		}
		backFuncValue := reflect.ValueOf(backFunc)
		if backFuncValue.Kind() != reflect.Func {
			panic("Regist object must be a func")
		}
		if backFuncValue.Type().NumIn() < 1 ||
			backFuncValue.Type().NumIn() != len(backParams)+1 ||
			backFuncValue.Type().In(0) != activityContextType {
			panic("First argument must use ActivityContext.")
		}
		newAction.RollbackFunc = backFuncValue
		newAction.RollbackParams = backParams
		return a
	}
}

func (a *Activity) Exec(ctx ActivityContext) {
	carg := reflect.ValueOf(ctx)
	for step, action := range a.Actions {
		args := append([]reflect.Value{carg}, action.DoParams...)
		result := action.DoFunc.Call(args)
		if isReturnError(result) {
			a.Rollback(ctx, step)
		}
	}
}

func (a *Activity) Rollback(ctx ActivityContext, fromStep int) {
	carg := reflect.ValueOf(ctx)
	for i := fromStep; i >= 0; i-- {
		action := a.Actions[i]
		args := append([]reflect.Value{carg}, action.RollbackParams...)
		result := action.RollbackFunc.Call(args)
		if isReturnError(result) {
			panic("!212")
		}
	}
}

func isReturnError(result []reflect.Value) bool {
	if len(result) == 1 && !result[0].IsNil() {
		return true
	}
	return false
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
