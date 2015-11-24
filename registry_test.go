package activity_test

import (
	"github.com/lysu/one-activity"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestRegisterFunc(t *testing.T) {

	f1 := func(a string) string { return a }

	r := activity.NewRegistry()
	r.Add("f1", f1)

	funcID := r.FindFuncID(reflect.ValueOf(f1))
	assert.Equal(t, funcID, "f1", "Find func ID by func value")

	fv := r.FindFunction("f1")
	assert.True(t, fv.IsValid(), "Find function by func ID")

	param := "abc"
	rt := fv.Call([]reflect.Value{reflect.ValueOf(param)})
	assert.Equal(t, "abc", rt[0].String(), "Call funcion")

	typeName := r.FindTypeName(reflect.ValueOf(param).Type())
	assert.Equal(t, "string", typeName)

	typ, ok := r.FindType("string")
	assert.True(t, ok)
	assert.Equal(t, "string", typ.Name())

}
