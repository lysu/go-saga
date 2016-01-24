package saga

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMarshalLog(t *testing.T) {
	l := &Log{
		Type:    ActionStart,
		SubTxID: "1",
		Params:  []ParamData{},
	}
	sl := l.mustMarshal()
	l2 := mustUnmarshalLog(sl)
	assert.Equal(t, ActionStart, l2.Type)
}
