package consumer

import (
	"context"
	"fmt"
	"github.com/segmentio/kafka-go"
	"go.temporal.io/sdk/workflow"
)

const (
	maxMessage = 2
)

var (
	reader *kafka.Reader
)

func InitReader() {
	brokers := []string{"192.168.49.2:32059"} // Kafka 代理地址
	topic := "my-topic"                       // Kafka 主题
	reader = kafka.NewReader(kafka.ReaderConfig{
		Brokers: brokers,
		Topic:   topic,
		GroupID: "my-group", // 消费者组ID
	})
}

func CloseReader() {
	if reader != nil {
		reader.Close()
		reader = nil
	}
}

// ChildWorkflow is a Workflow Definition
func ChildWorkflow(ctx workflow.Context, idx int) (string, error) {
	logger := workflow.GetLogger(ctx)
	logger.Info("Consumer child workflow execution", "index", idx)

	for i := 0; i < maxMessage; i++ {
		dl, ok := ctx.Deadline()
		if !ok {
			return "", nil
		}
		msgCtx, _ := context.WithDeadline(context.Background(), dl)
		msg, err := reader.ReadMessage(msgCtx)
		if err != nil {
			return "", nil
		}

		fmt.Printf("Message on %s: %s\n", msg.Topic, string(msg.Value))
	}
	return "ok", nil
}
