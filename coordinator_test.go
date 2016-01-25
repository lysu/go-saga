package saga

import "golang.org/x/net/context"

// This example show how to use init a SEC and add Sub-transaction to it, finally trigger a transfer operation as demo.
// In a saga transaction we deduce `100` from foo at first, then deposit 100 into `bar`.
func ExampleExecutionCoordinator_transfer() {

	// 1. Define sub-transaction method, anonymous method is NOT required, Just define them as normal way.

	DeduceAccount := func(ctx context.Context, account string, amount int) error {
		// Do deduce amount from account
		return nil
	}
	CompensateDeduce := func(ctx context.Context, account string, amount int) error {
		// Compensate deduce amount to account
		return nil
	}
	DepositAccount := func(ctx context.Context, account string, amount int) error {
		// Do deposit amount to account
		return nil
	}
	CompensateDeposit := func(ctx context.Context, account string, amount int) error {
		// Compensate deposit amount from account
		return nil
	}

	// 2. Init SEC as global SINGLETON(this demo not..), and add Sub-transaction definition into SEC.

	storage, _ := NewMemStorage()
	sec := NewSEC(storage)
	sec.AddSubTxDef("deduce", DeduceAccount, CompensateDeduce).
		AddSubTxDef("deposit", DepositAccount, CompensateDeposit)

	// 3. Start a saga to transfer 100 from foo to bar.

	from, to := "foo", "bar"
	amount := 100
	ctx := context.Background()

	var sagaID uint64 = 1
	sec.StartSaga(ctx, sagaID).
		SubTx("deduce", from, amount).
		SubTx("deposit", to, amount).
		EndSaga()

	// 4. done.
}
