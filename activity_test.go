package activity_test

import (
	"fmt"
	"github.com/lysu/one-activity"
	"strconv"
	"testing"
)

func Call1(ctx activity.ActivityContext, abc int) error {
	fmt.Println("call1 " + strconv.Itoa(abc))
	return fmt.Errorf("1212")
}

func Rollback1(ctx activity.ActivityContext) error {
	fmt.Println("rolled")
	return nil
}

var reg *activity.Registry

var storage activity.Storage

func initIt() {
	storage, _ = activity.NewMemStorage()
	reg = activity.NewRegistry().
		Add("call1", Call1).Add("rollback1", Rollback1)
}

func TestOneActivityExec(t *testing.T) {

	// initInStartup
	initIt()

	// execute activities
	ctx := activity.ActivityContext{}
	activity.Start(storage, reg, 3).
		Then(Call1, 1)(Rollback1).
		Run(ctx)

}
