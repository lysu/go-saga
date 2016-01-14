package saga_test

import (
	"fmt"
	"github.com/lysu/go-saga"
	"strconv"
	"testing"
)

func Call1(ctx saga.ActivityContext, abc int) error {
	fmt.Println("call1 " + strconv.Itoa(abc))
	return fmt.Errorf("1212")
}

func Rollback1(ctx saga.ActivityContext) error {
	fmt.Println("rolled")
	return nil
}

var reg *saga.Registry

var storage saga.Storage

func initIt() {
	storage, _ = saga.NewMemStorage()
	reg = saga.NewRegistry().
		Add("call1", Call1).Add("rollback1", Rollback1)
}

func TestOneActivityExec(t *testing.T) {

	// initInStartup
	initIt()

	// execute activities
	ctx := saga.ActivityContext{}
	saga.Start(storage, reg, 3).
		Then(Call1, 1)(Rollback1).
		Run(ctx)

}
