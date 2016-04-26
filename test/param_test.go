package saga_test

import (
	"fmt"
	"github.com/lysu/go-saga"
	"golang.org/x/net/context"
	"reflect"
	"testing"
)

func Param1(ctx context.Context, name string, aga int) {
	fmt.Printf("%s----%d\n", name, aga)
}

func initParam() {
	saga.AddSubTxDef("param1", Param1, Param2)
}

func TestMarshalParam(t *testing.T) {
	initParam()
	pd := saga.MarshalParam(&saga.DefaultSEC, []interface{}{"a", 1})
	rv := saga.UnmarshalParam(&saga.DefaultSEC, pd)

	p := []reflect.Value{}
	p = append(p, reflect.ValueOf(context.Background()))
	p = append(p, rv...)

	f := reflect.ValueOf(Param1)

	f.Call(p)

}

func Param2(ctx context.Context, name *string, aga int) {
	fmt.Printf("%v----%d\n", name, aga)
}

func TestMarshalPtr(t *testing.T) {
	initParam()
	x := "a"
	pd := saga.MarshalParam(&saga.DefaultSEC, []interface{}{&x, 1})
	rv := saga.UnmarshalParam(&saga.DefaultSEC, pd)

	p := []reflect.Value{}
	p = append(p, reflect.ValueOf(context.Background()))
	p = append(p, rv...)

	f := reflect.ValueOf(Param2)

	f.Call(p)

}
