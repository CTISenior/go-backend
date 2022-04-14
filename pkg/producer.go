package pkg

import (
	"fmt"

	"github.com/Shopify/sarama"
)

//var brokers = []string{"localhost:9091", "localhost:9092", "localhost:9093"}

func Producer() (sarama.SyncProducer, error) {
	config := sarama.NewConfig()
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Return.Successes = true
	producer, err := sarama.NewSyncProducer(Kafka_Brokers, config)

	return producer, err
}

func createMessage(topic, message string) *sarama.ProducerMessage {
	msg := &sarama.ProducerMessage{
		Topic:     topic,
		Partition: -1,
		Value:     sarama.StringEncoder(message),
	}

	return msg
}

var producer, _ = Producer()

func ProduceMessage(topic, message string) {
	msg := createMessage(topic, message)
	producer.SendMessage(msg)
	logMsg := fmt.Sprintf("Kafka producer message - Topic %s -> message %s\n", topic, message)
	fmt.Printf(logMsg + "\n")
	//InfoLogger.Println(logMsg)

}
