package pkg

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/Shopify/sarama"
)

type Kafka struct {
	Brokers           []string
	topics            []string
	startOffset       int64
	version           string
	ready             chan bool
	group             string
	channelBufferSize int
}

func NewKafka() *Kafka {
	return &Kafka{
		Brokers: Kafka_Brokers,
		topics: []string{
			topics,
		},
		group:             group,
		channelBufferSize: 2,
		ready:             make(chan bool),
		version:           "1.1.1",
	}
}

//var Brokers = []string{"localhost:9091", "localhost:9092", "localhost:9093"}
var topics = "test"
var group = "1"

func (p *Kafka) Init() func() {
	log.Println("Kafka Init")

	version, err := sarama.ParseKafkaVersion(p.version)
	if err != nil {
		log.Fatalf("Error parsing Kafka version: %v", err)
	}
	config := sarama.NewConfig()
	config.Version = version
	config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRange
	config.Consumer.Offsets.Initial = -2
	config.ChannelBufferSize = p.channelBufferSize

	ctx, cancel := context.WithCancel(context.Background())
	client, err := sarama.NewConsumerGroup(p.Brokers, p.group, config)
	if err != nil {
		log.Fatalf("Error creating consumer group client: %v", err)
	}

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer func() {
			wg.Done()
		}()
		for {
			if err := client.Consume(ctx, p.topics, p); err != nil {
				log.Fatalf("Error from consumer: %v", err)
			}
			if ctx.Err() != nil {
				log.Println(ctx.Err())
				return
			}
			p.ready = make(chan bool)
		}
	}()
	<-p.ready
	log.Println("Sarama consumer up and running!...")
	return func() {
		log.Println("kafka close")
		cancel()
		wg.Wait()
		if err = client.Close(); err != nil {
			log.Fatal("Error closing client: %v", err)
		}
	}
}

func (p *Kafka) Setup(sarama.ConsumerGroupSession) error {
	close(p.ready)
	return nil
}

func (p *Kafka) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (k *Kafka) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {

	for message := range claim.Messages() {
		fmt.Println("[topic:%s] [partiton:%d] [offset:%d] [value:%s] [time:%v]",
			message.Topic, message.Partition, message.Offset, string(message.Value), message.Timestamp)

		session.MarkMessage(message, "")
	}
	return nil
}

/*
func main() {
	k := NewKafka()
	c := k.Connect()

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)
	select {
	case <-sigterm:
		fmt.Println("terminating: via signal")
	}
	c()
}
*/
