package activity_test

import (
	"testing"
	"github.com/lysu/one-activity"
	"golang.org/x/net/context"
	"fmt"
)

func Call1(ctx context.Context) error {
	fmt.Println("call1")
	return nil
}

func Rollback1(ctx context.Context) error {
	return nil
}

func Call2(ctx context.Context) error {
	fmt.Println("call2")
	return nil
}

func Rollback2(ctx context.Context) error {
	return nil
}

func initIt() {
	activity.NewRegistry().
		Add("call1", Call1).Add("call2", Call2).
		Add("rollback1", Rollback1).Add("rollback2", Rollback2)
}

func TestOneActivityExec(t *testing.T) {

	// initInStartup
	initIt()

	// execute activities
	ctx := context.TODO()
	activity.Start(3).
		Then(Call1, 1)(Rollback1).
		Then(Call2, 2, "233")(Rollback2, 2, "233").
	   Run(ctx)

}
