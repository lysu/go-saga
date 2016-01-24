package saga_test

import (
	"github.com/lysu/go-saga"
	"github.com/stretchr/testify/assert"
	"testing"
)

func T1(ctx saga.SagaContext) {

}

func C1(ctx saga.SagaContext) {

}

func T2(ctx saga.SagaContext) {

}

func C2(ctx saga.SagaContext) {

}

func TestSubTxDefinitions(t *testing.T) {
	txs := saga.SubTxDefinitions{}.
		AddDefinition("A1", T1, C1).
		AddDefinition("A2", T2, C2)
	define, ok := txs.FindDefinition("A1")
	assert.True(t, ok)
	assert.NotNil(t, define.Action)
}

func E() {

}

func TestMissFunc(t *testing.T) {
	func() {
		defer func() {
			if r := recover(); r != nil {
				assert.Equal(t, "First argument must use SagaContext.", r)
				return
			}
			assert.Fail(t, "It must be panic when use E function")
		}()
		saga.SubTxDefinitions{}.AddDefinition("Test", T1, E)
	}()
}
