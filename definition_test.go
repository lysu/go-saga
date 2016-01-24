package saga

import (
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
	"testing"
)

func T1(ctx context.Context) {

}

func C1(ctx context.Context) {

}

func T2(ctx context.Context) {

}

func C2(ctx context.Context) {

}

func TestSubTxDefinitions(t *testing.T) {
	txs := subTxDefinitions{}.
		addDefinition("A1", T1, C1).
		addDefinition("A2", T2, C2)
	define, ok := txs.findDefinition("A1")
	assert.True(t, ok)
	assert.NotNil(t, define.action)
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
		subTxDefinitions{}.addDefinition("Test", T1, E)
	}()
}
