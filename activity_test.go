package activity_test

import (
	"fmt"
	"github.com/lysu/one-activity"
	"testing"
)

func Call1(ctx activity.ActivityContext) error {
	fmt.Println("call1")
	return nil
}

func Rollback1(ctx activity.ActivityContext) error {
	return nil
}

func Call2(ctx activity.ActivityContext) error {
	fmt.Println("call2")
	return nil
}

func Rollback2(ctx activity.ActivityContext) error {
	return nil
}

var reg *activity.Registry

var storage *activity.DBStorage

func initIt() {
	storage = nil // TODO..
	activity.NewRegistry().
		Add("call1", Call1).Add("call2", Call2).
		Add("rollback1", Rollback1).Add("rollback2", Rollback2)
}

func TestOneActivityExec(t *testing.T) {

	// initInStartup
	initIt()

	// execute activities
	ctx := activity.ActivityContext{}
	activity.Start(storage, reg, 3).
		Then(Call1, 1)(Rollback1).
		Then(Call2, 2, "233")(Rollback2, 2, "233").
		Run(ctx)

}
