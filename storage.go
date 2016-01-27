package saga

import (
	"fmt"
	"github.com/Shopify/sarama"
	"log"
)

// Storage uses to support save and lookup saga log.
type Storage interface {

	// AppendLog appends log data into log under given logID
	AppendLog(logID string, data string) error

	// Lookup uses to lookup all log under given logID
	Lookup(logID string) ([]string, error)

	// Close use to close storage and release resources
	Close() error
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

// Close use to close storage and release resources.
func (s *memStorage) Close() error {
	return nil
}

type kafkaStorage struct {
	producer sarama.SyncProducer
	consumer sarama.Consumer
}

// NewKafkaStorage creates log storage base on Kafka.
func NewKafkaStorage(addrs []string) (Storage, error) {
	producer, err := sarama.NewSyncProducer(addrs, nil)
	if err != nil {
		panic(fmt.Sprintf("Start Kafka Storage failure: %v", err))
	}
	consumer, err := sarama.NewConsumer([]string{"localhost:9092"}, nil)
	if err != nil {
		panic(err)
	}
	return &kafkaStorage{
		producer: producer,
		consumer: consumer,
	}, nil
}

// AppendLog appends log into queue under given logID.
func (s *kafkaStorage) AppendLog(logID string, data string) error {
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

	data := []string{}
	consumed := 0
	for {
		select {
		case msg := <-partitionConsumer.Messages():
			log.Printf("Consumed message offset %d\n", msg.Offset)
			consumed++
			msgValue := string(msg.Value)
			data = append(data, msgValue)
			if msgValue == "" { // ???
				break
			}
		}
	}

	//log.Printf("Consumed: %d\n", consumed)
	//return data, nil
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
