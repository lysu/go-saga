package saga

import (
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/wvanbergen/kazoo-go"
	"log"
	"strings"
	"time"
)

// Storage uses to support save and lookup saga log.
type Storage interface {

	// AppendLog appends log data into log under given logID
	AppendLog(logID string, data string) error

	// Lookup uses to lookup all log under given logID
	Lookup(logID string) ([]string, error)

	// Close use to close storage and release resources
	Close() error

	// LogIDs returns exists logID
	LogIDs() ([]string, error)

	// Cleanup cleans up all log data in logID
	Cleanup(logID string) error

	// LastLog fetch last log entry with given logID
	LastLog(logID string) (string, error)
}

type memStorage struct {
	data map[string][]string
}

// NewMemStorage creates log storage base on memory.
// This storage use simple `map[string][]string`, just for TestCase used.
// NOT use this in product.
func NewMemStorage() (Storage, error) {
	return &memStorage{
		data: make(map[string][]string),
	}, nil
}

// AppendLog appends log into queue under given logID.
func (s *memStorage) AppendLog(logID string, data string) error {
	logQueue, ok := s.data[logID]
	if !ok {
		logQueue = []string{}
		s.data[logID] = logQueue
	}
	s.data[logID] = append(s.data[logID], data)
	return nil
}

// Lookup lookups log under given logID.
func (s *memStorage) Lookup(logID string) ([]string, error) {
	return s.data[logID], nil
}

// Close uses to close storage and release resources.
func (s *memStorage) Close() error {
	return nil
}

// LogIDs uses to take all Log ID av in current storage
func (s *memStorage) LogIDs() ([]string, error) {
	ids := make([]string, 0, len(s.data))
	for id := range s.data {
		ids = append(ids, id)
	}
	return ids, nil
}

func (s *memStorage) Cleanup(logID string) error {
	delete(s.data, logID)
	return nil
}

func (s *memStorage) LastLog(logID string) (string, error) {
	logData, ok := s.data[logID]
	if !ok {
		return "", fmt.Errorf("LogData %s not found")
	}
	sizeOfLog := len(logData)
	if sizeOfLog == 0 {
		return "", fmt.Errorf("LogData is empty")
	}
	lastLog := logData[sizeOfLog-1]
	return lastLog, nil
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
func NewKafkaStorage(zkAddrs, brokerAddrs []string, partitions, replicas int, returnDuration time.Duration) (Storage, error) {
	conf := kazoo.NewConfig()
	kz, err := kazoo.NewKazoo(zkAddrs, conf)
	if err != nil {
		panic(fmt.Sprintf("Start Zookeeper client failure: %v", err))
	}
	producer, err := sarama.NewSyncProducer(brokerAddrs, nil)
	if err != nil {
		panic(fmt.Sprintf("Start Kafka Storage failure: %v", err))
	}
	consumer, err := sarama.NewConsumer([]string{"localhost:9092"}, nil)
	if err != nil {
		panic(err)
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
		return err
	}
	if !topicExists {
		err = s.kz.CreateTopic(logID, s.partitionNumbers, s.replicaNumbers, map[string]interface{}{})
		if err != nil {
			return err
		}
	}
	msg := &sarama.ProducerMessage{Topic: logID, Value: sarama.StringEncoder(data)} // ?? always new?
	partition, offset, err := s.producer.SendMessage(msg)
	if err != nil {
		log.Printf("FAILED to send message: %s\n", err)
		return err
	}
	log.Printf("> message sent to partition %d at offset %d\n", partition, offset)
	return nil
}

// Lookup lookups log under given logID.
func (s *kafkaStorage) Lookup(logID string) ([]string, error) {
	partitionConsumer, err := s.consumer.ConsumePartition(logID, 0, sarama.OffsetOldest)
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := partitionConsumer.Close(); err != nil {
			log.Fatalln(err)
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
			log.Printf("Consumed message offset %d\n", msg.Offset)
			consumed++
			msgValue := string(msg.Value)
			data = append(data, msgValue)
			timer.Reset(s.consumeReturnDuration)
		case <-timer.C:
			break consumer_loop
		}
	}

	log.Printf("Consumed: %d\n", consumed)
	return data, nil
}

// Close use to close storage and release resources.
func (s *kafkaStorage) Close() error {
	if err1 := s.producer.Close(); err1 != nil {
		log.Println(err1)
	}
	if err2 := s.consumer.Close(); err2 != nil {
		log.Println(err2)
	}
	return nil
}

// LogIDs returns av saga topic in kafka.
func (s *kafkaStorage) LogIDs() ([]string, error) {
	topics, err := s.kz.Topics()
	if err != nil {
		return nil, err
	}
	sagaTopics := make([]string, 0, len(topics))
	for _, topic := range topics {
		if strings.HasPrefix(topic.Name, logPrefix) {
			sagaTopics = append(sagaTopics, topic.Name)
		}
	}
	return sagaTopics, nil
}

// Cleanup cleans log data for given logID
func (s *kafkaStorage) Cleanup(logID string) error {
	return s.kz.DeleteTopic(logID)
}

// LastLog consume last log
func (s *kafkaStorage) LastLog(logID string) (string, error) {
	return "", nil
}
