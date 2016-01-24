package saga_test

import (
	"fmt"
	"github.com/lysu/go-saga"
	"reflect"
	"testing"
)

func Param1(name string, aga int) {
	fmt.Printf("%s----%d\n", name, aga)
}

func TestMarshalParam(t *testing.T) {
	initIt(OK)

	pd := saga.MarshalParam(&sec, []interface{}{"a", 1})
	rv := saga.UnmarshalParam(&sec, pd)

	f := reflect.ValueOf(Param1)

	f.Call(rv)

}

func Param2(name *string, aga int) {
	fmt.Printf("%v----%d\n", name, aga)
}

func TestMarshalPtr(t *testing.T) {
	initIt(OK)
	x := "a"
	pd := saga.MarshalParam(&sec, []interface{}{&x, 1})
	rv := saga.UnmarshalParam(&sec, pd)

	f := reflect.ValueOf(Param2)

	f.Call(rv)

}
