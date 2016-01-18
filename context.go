package saga

import (
	"golang.org/x/net/context"
	"reflect"
)

var sagaContextType = reflect.TypeOf(SagaContext{})

type SagaContext struct {
	context.Context
}
