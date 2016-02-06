package saga_test

import (
	"github.com/lysu/go-saga"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestKafkaLogAppend(t *testing.T) {

	assert := assert.New(t)

	s, err := saga.NewKafkaStorage(
		[]string{"0.0.0.0:2181"},
		[]string{"0.0.0.0:9092"},
		1,
		1,
		50*time.Millisecond,
	)
	assert.NoError(err)

	err = s.AppendLog("d1", "123456")
	assert.NoError(err)

	data, err := s.Lookup("d1")
	assert.NoError(err)

	t.Log(data)

}
