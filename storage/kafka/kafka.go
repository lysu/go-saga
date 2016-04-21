package kafka

import (
	"github.com/Shopify/sarama"
	"github.com/juju/errors"
	"github.com/lysu/go-saga"
	"github.com/lysu/go-saga/storage"
	"github.com/lysu/kazoo-go"
	"strings"
	"sync"
	"time"
)

var storageInstance storage.Storage
var kafkaInit sync.Once

func init() {
	saga.StorageProvider = func(cfg storage.StorageConfig) storage.Storage {
		kafkaInit.Do(func() {
			var err error
			storageInstance, err = newKafkaStorage(
				cfg.Kafka.ZkAddrs, cfg.Kafka.BrokerAddrs, cfg.Kafka.Partitions,
				cfg.Kafka.Replicas, cfg.Kafka.ReturnDuration,
			)
			if err != nil {
				panic(err)
			}
		})
		return storageInstance
	}
}

type kafkaStorage struct {
	producer              sarama.SyncProducer
	consumer              sarama.Consumer
	kz                    *kazoo.Kazoo
	partitionNumbers      int
	replicaNumbers        int
	consumeReturnDuration time.Duration
}

// NewKafkaStorage creates log storage base on Kafka.
func newKafkaStorage(zkAddrs, brokerAddrs []string, partitions, replicas int, returnDuration time.Duration) (storage.Storage, error) {
	conf := kazoo.NewConfig()
	kz, err := kazoo.NewKazoo(zkAddrs, conf)
	if err != nil {
		return nil, errors.Annotate(err, "Start Zookeeper client failure")
	}
	producer, err := sarama.NewSyncProducer(brokerAddrs, nil)
	if err != nil {
		return nil, errors.Annotatef(err, "Start Kafka Storage failure: %v", brokerAddrs)
	}
	consumer, err := sarama.NewConsumer(brokerAddrs, nil)
	if err != nil {
		return nil, errors.Annotatef(err, "Create Consumer failure: %v", brokerAddrs)
	}
	return &kafkaStorage{
		producer:              producer,
		consumer:              consumer,
		kz:                    kz,
		partitionNumbers:      partitions,
		replicaNumbers:        replicas,
		consumeReturnDuration: returnDuration,
	}, nil
}

// AppendLog appends log into queue under given logID.
func (s *kafkaStorage) AppendLog(logID string, data string) error {
	topicExists, err := s.kz.ExistsTopic(logID)
	if err != nil {
		return errors.Annotatef(err, "for %s", logID)
	}
	if !topicExists {
		err = s.kz.CreateTopic(logID, s.partitionNumbers, s.replicaNumbers, map[string]interface{}{})
		if err != nil {
			return errors.Annotatef(err, "for topic %s", logID)
		}
	}
	msg := &sarama.ProducerMessage{Topic: logID, Value: sarama.StringEncoder(data)} // ?? always new?
	partition, offset, err := s.producer.SendMessage(msg)
	if err != nil {
		return errors.Annotatef(err, " failure send %s", data)
	}
	saga.Logger.Printf("> message sent to partition %d at offset %d\n", partition, offset)
	return nil
}

// Lookup lookups log under given logID.
func (s *kafkaStorage) Lookup(logID string) ([]string, error) {
	partitionConsumer, err := s.consumer.ConsumePartition(logID, 0, sarama.OffsetOldest)
	if err != nil {
		return nil, errors.Annotatef(err, "Consume topic %s failured", logID)
	}

	defer func() {
		if err := partitionConsumer.Close(); err != nil {
			saga.Logger.Printf("[WARNING]Close consumer failure %v", err)
		}
	}()

	timer := time.NewTimer(s.consumeReturnDuration)
	defer timer.Stop()
	data := []string{}
	consumed := 0
consumer_loop:
	for {
		select {
		case msg := <-partitionConsumer.Messages():
			saga.Logger.Printf("Consumed message offset %d\n", msg.Offset)
			consumed++
			msgValue := string(msg.Value)
			data = append(data, msgValue)
			timer.Reset(s.consumeReturnDuration)
		case <-timer.C:
			break consumer_loop
		}
	}

	saga.Logger.Printf("Consumed: %d\n", consumed)
	return data, nil
}

// Close use to close storage and release resources.
func (s *kafkaStorage) Close() error {
	if err1 := s.producer.Close(); err1 != nil {
		return errors.Annotate(err1, "Close producer failure")
	}
	if err2 := s.consumer.Close(); err2 != nil {
		return errors.Annotate(err2, "Close consumer failure")
	}
	return nil
}

// LogIDs returns av saga topic in kafka.
func (s *kafkaStorage) LogIDs() ([]string, error) {
	topics, err := s.kz.Topics()
	if err != nil {
		return nil, errors.Annotate(err, "Get topic info failure")
	}
	sagaTopics := make([]string, 0, len(topics))
	for _, topic := range topics {
		if strings.HasPrefix(topic.Name, saga.LogPrefix) {
			sagaTopics = append(sagaTopics, topic.Name)
		}
	}
	return sagaTopics, nil
}

// Cleanup cleans log data for given logID
func (s *kafkaStorage) Cleanup(logID string) error {
	err := s.kz.DeleteTopic(logID)
	if err != nil {
		return errors.Annotatef(err, "Delete topic %s failure", logID)
	}
	return nil
}

// LastLog consume last log
func (s *kafkaStorage) LastLog(logID string) (string, error) {
	return "", nil
}
