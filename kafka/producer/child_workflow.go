package producer

import (
	"context"
	"github.com/segmentio/kafka-go"
	"go.temporal.io/sdk/workflow"
	"strconv"
)

var (
	writer *kafka.Writer
)

func InitWriter() {
	brokers := []string{"192.168.49.2:32059"} // Kafka 代理地址
	topic := "my-topic"                       // Kafka 主题

	writer = kafka.NewWriter(kafka.WriterConfig{
		Brokers:  brokers,
		Topic:    topic,
		Balancer: &kafka.RoundRobin{},
	})
}

func CloseWriter() {
	if writer != nil {
		writer.Close()
	}
}

// ChildWorkflow is a Workflow Definition
func ChildWorkflow(ctx workflow.Context, idx int) (string, error) {
	logger := workflow.GetLogger(ctx)
	logger.Info("Producer child workflow execution", "index", idx)

	msg := kafka.Message{
		Value: []byte("temporal message" + strconv.Itoa(idx)),
	}
	dl, ok := ctx.Deadline()
	if !ok {
		return "", nil
	}
	msgCtx, _ := context.WithDeadline(context.Background(), dl)
	writer.WriteMessages(msgCtx, msg)
	return "ok", nil
}
