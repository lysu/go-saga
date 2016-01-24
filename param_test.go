package saga_test

import (
	"fmt"
	"golang.org/x/net/context"
	"reflect"
	"testing"
)

func Param1(name string, aga int) {
	fmt.Printf("%s----%d\n", name, aga)
}

func TestMarshalParam(t *testing.T) {
	initIt(OK)
	ctx := context.Background()

	var sagaID uint64 = 1
	s := sec.StartSaga(ctx, sagaID)
	pd := s.MarshalParam([]interface{}{"a", 1})
	rv := s.UnmarshalParam(pd)

	f := reflect.ValueOf(Param1)

	f.Call(rv)

}

func Param2(name *string, aga int) {
	fmt.Printf("%s----%d\n", name, aga)
}

func TestMarshalPtr(t *testing.T) {
	initIt(OK)
	ctx := context.Background()

	var sagaID uint64 = 1
	s := sec.StartSaga(ctx, sagaID)
	x := "a"
	pd := s.MarshalParam([]interface{}{&x, 1})
	rv := s.UnmarshalParam(pd)

	f := reflect.ValueOf(Param2)

	f.Call(rv)

}
