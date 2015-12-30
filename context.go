package activity

import (
	"golang.org/x/net/context"
	"reflect"
)

var activityContextType = reflect.TypeOf(ActivityContext{})

type ActivityContext struct {
	context.Context
}
