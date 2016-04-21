package saga_test

import (
	"github.com/lysu/go-saga"
	_ "github.com/lysu/go-saga/storage/kafka"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func KafkaLogAppend(t *testing.T) {

	assert := assert.New(t)
	saga.StorageConfig.Kafka.ZkAddrs = []string{"0.0.0.0:2181"}
	saga.StorageConfig.Kafka.BrokerAddrs = []string{"0.0.0.0:9092"}
	saga.StorageConfig.Kafka.Partitions = 1
	saga.StorageConfig.Kafka.Replicas = 1
	saga.StorageConfig.Kafka.ReturnDuration = 50 * time.Millisecond

	err := saga.LogStorage().AppendLog("d1", "123456")
	assert.NoError(err)

	data, err := saga.LogStorage().Lookup("d1")
	assert.NoError(err)

	t.Log(data)

}
