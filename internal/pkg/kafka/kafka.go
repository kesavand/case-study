package kafka

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Shopify/sarama"
)

//kafkahndlr
type Producer struct {
	producer sarama.SyncProducer
}

//Consumer logic for sarama library
type Consumer struct {
	saramaCg   sarama.ConsumerGroup
	HandleFunc func(msg string) error
	ID         string
	topics     string
}

//Options to configure consumption and Produceing Options
const (
	clientID        = "TEST"
	WriteTimeout    = 2
	ConsumerGroup   = "test-cg"
	kafkaversion    = "2.0.0"
	kafkaKey        = "test-key"
	maxReadAttempts = 3
	readWait        = 1
)

//kafkaInterface
type kafkaProducer interface {
	Produce(ctx context.Context, topic string, payload string) (err error)
}

type kafkaConsumer interface {
	Read(ctx context.Context, topics string, errChn chan error, handleMessage func(msg string) error) error
}

//NewKafkaHndlr constructor
func NewKafkaProducer(broker string) (kafkaProducer, error) {
	//Producer Configuration
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true                 // writes fail messages to che fail channel
	config.ClientID = clientID                           // application specific id for logging
	config.Net.WriteTimeout = time.Second * WriteTimeout // how long to wait for transmit

	producer, err := sarama.NewSyncProducer([]string{broker}, config)
	if err != nil {
		return nil, err
	}

	kp := &Producer{
		producer: producer,
	}

	return kp, nil
}

//NewKafkaHndlr constructor
func NewKafkaConsumer(broker string) (kafkaConsumer, error) {
	var err error
	// Consumer configuration
	conf := sarama.NewConfig()
	conf.Consumer.Return.Errors = true
	conf.Consumer.Offsets.Initial = sarama.OffsetOldest

	conf.Version, err = sarama.ParseKafkaVersion(kafkaversion)
	if err != nil {
		return nil, err
	}

	cg, err := sarama.NewConsumerGroup([]string{broker}, ConsumerGroup, conf)
	if err != nil {
		return nil, err
	}

	kc := &Consumer{
		saramaCg: cg,
	}

	return kc, nil
}

//Close closes the kafka client
func (kp Producer) Close() error {
	if kp.producer != nil {
		err := kp.producer.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

//Close closes the kafka client
func (kc Consumer) Close() error {
	err := kc.saramaCg.Close()
	if err == nil {
		fmt.Printf("Closed kafka consumer")
	} else {
		return err
	}

	return nil
}

//Produce message to the Kafka stream
func (kh Producer) Produce(ctx context.Context, topic string, payload string) (err error) {
	msg := &sarama.ProducerMessage{
		Key:   sarama.StringEncoder(kafkaKey),
		Topic: topic,
		Value: sarama.StringEncoder(payload),
	}

	part, off, err := kh.producer.SendMessage(msg)
	if err != nil {
		return err
	}
	fmt.Printf("Produceed msg '%s',  to partition %d topic %s with %d offset\n", payload, part, msg.Topic, off)

	return nil
}

//Read kafka msgs
func (kh *Consumer) Read(ctx context.Context, topics string, errrChn chan error, handleMessage func(msg string) error) error {
	consumer := &Consumer{
		HandleFunc: handleMessage,
		//	ID:         strconv.Itoa(0),
	}

	go func() {
		for err := range kh.saramaCg.Errors() {
			errrChn <- err
		}
	}()

	go func() {
		for {
			if err := kh.saramaCg.Consume(ctx, strings.Split(topics, ","), consumer); err != nil {
				fmt.Printf("error in reading kafka message %v\n, err", err)
			}
			// check if context was canceled, signaling that the consumer should stop
			if ctx.Err() != nil {
				fmt.Printf("will exit consumer since context was stopped")
				return
			}
		}
	}()
	fmt.Printf("Consumer started successfully")
	return nil
}

func (consumer *Consumer) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
func (consumer *Consumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages().
func (consumer *Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	var failedAttempts = 0
	for msg := range claim.Messages() {
		log.Printf("Message claimed: value = %s, timestamp = %v, topic = %s", string(msg.Value), msg.Timestamp, msg.Topic)
		err := consumer.HandleFunc(string(msg.Value))
		if err != nil {
			failedAttempts++
			if failedAttempts < maxReadAttempts {
				fmt.Printf("read failed for  %d attempts.", failedAttempts)
				session.MarkMessage(msg, "")
				return err
			} else {
				time.Sleep(time.Second * readWait)
				fmt.Printf("Read failed retry")
			}

		}

		session.MarkMessage(msg, "")
	}

	return nil
}
