package saga_test

import (
	"github.com/lysu/go-saga"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMarshalLog(t *testing.T) {
	l := &saga.Log{
		Type:    saga.ActionStart,
		SubTxID: "1",
		Params:  []saga.ParamData{},
	}
	sl := l.MustMarshal()
	l2 := saga.MustUnmarshal(sl)
	assert.Equal(t, saga.ActionStart, l2.Type)
}
