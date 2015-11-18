package activity

import (
	"golang.org/x/net/context"
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

func Start(biz int) *Activitor {
	return &Activitor{
		ID:        1,
		Status:    ActivityStarted,
		StartTime: time.Now(),
		Actions:   []Action{},
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

func (a *Activitor) Run(ctx context.Context) error {
	return nil
}
