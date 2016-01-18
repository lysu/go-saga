package saga_test

import (
	"fmt"
	"github.com/lysu/go-saga"
	"github.com/stretchr/testify/assert"
	"testing"
)

type FailureMode int

const (
	OK = iota
	DeduceFail
	DepositFail
)

var (
	reg      *saga.Registry
	storage  saga.Storage
	memDB    map[string]int
	testMode FailureMode
)

func initIt(mode FailureMode) {
	storage, _ = saga.NewMemStorage()
	reg = saga.NewRegistry().
		Add("DeduceAccount", DeduceAccount).Add("CompensateDeduce", CompensateDeduce).
		Add("DepositAccount", DepositAccount).Add("CompensateDeposit", CompensateDeposit)
	memDB = map[string]int{
		"foo": 200,
		"bar": 0,
	}
	testMode = mode
}

func DeduceAccount(ctx saga.SagaContext, account string, amount int) error {
	if testMode == DeduceFail {
		return fmt.Errorf("Deduce failure")
	}
	memDB[account] = (memDB[account] - amount)
	return nil
}

func CompensateDeduce(ctx saga.SagaContext, account string, amount int) error {
	memDB[account] = (memDB[account] + amount)
	return nil
}

func DepositAccount(ctx saga.SagaContext, account string, amount int) error {
	if testMode == DepositFail {
		return fmt.Errorf("Deposit failure")
	}
	memDB[account] = (memDB[account] + amount)
	return nil
}

func CompensateDeposit(ctx saga.SagaContext, account string, amount int) error {
	memDB[account] = (memDB[account] - amount)
	return nil
}

func TestSagaWithBothSuccess(t *testing.T) {

	// initInStartup
	initIt(OK)

	from, to := "foo", "bar"
	amount := 100

	// execute activities
	ctx := saga.SagaContext{}
	saga.Def(storage, reg).
		Then(DeduceAccount, from, amount)(CompensateDeduce, from, amount).
		Then(DepositAccount, to, amount)(CompensateDeposit, to, amount).
		Run(ctx)

	// assert
	assert.Equal(t, 100, memDB[from])
	assert.Equal(t, 100, memDB[to])

}

func TestSagaWithDepositFail(t *testing.T) {

	// initInStartup
	initIt(DepositFail)

	from, to := "foo", "bar"
	amount := 100

	// execute activities
	ctx := saga.SagaContext{}
	saga.Def(storage, reg).
		Then(DeduceAccount, from, amount)(CompensateDeduce, from, amount).
		Then(DepositAccount, to, amount)(CompensateDeposit, to, amount).
		Run(ctx)

	// assert
	assert.Equal(t, 200, memDB[from])
	assert.Equal(t, 0, memDB[to])

}
