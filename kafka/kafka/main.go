package main

import (
	"context"
	"fmt"
	"github.com/segmentio/kafka-go"
	"time"
)

// ChildWorkflow is a Workflow Definition
func main() {
	brokers := []string{"192.168.49.2:32059"} // Kafka 代理地址
	topic := "my-topic"                       // Kafka 主题

	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers:  brokers,
		Topic:    topic,
		Balancer: &kafka.RoundRobin{},
	})
	defer writer.Close()

	msg := kafka.Message{
		Value: []byte("temporal message"),
	}
	msgCtx, _ := context.WithTimeout(context.Background(), time.Second*10)
	err := writer.WriteMessages(msgCtx, msg)
	if err != nil {
		fmt.Printf("Failed to send message :%v", err)
		return
	}
	fmt.Printf("send message ok")
}
