package consumer

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"github.com/segmentio/kafka-go"
	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/workflow"
	"time"
)

const (
	maxMessage = 2
)

var (
	TC     client.Client
	reader *kafka.Reader
)

func InitReader(tEndpoint, kEndpoint, kTopic string) error {
	c, err := client.Dial(client.Options{
		HostPort: tEndpoint,
	})
	if err != nil {
		return errors.Wrap(err, "Unable to create client")
	}
	TC = c
	brokers := []string{kEndpoint} // Kafka 代理地址
	reader = kafka.NewReader(kafka.ReaderConfig{
		Brokers: brokers,
		Topic:   kTopic,
		GroupID: "my-group", // 消费者组ID
	})
	return nil
}

func CloseReader() {
	if reader != nil {
		reader.Close()
		reader = nil
	}
	if TC != nil {
		TC.Close()
	}
}

// ChildWorkflow is a Workflow Definition
func ChildWorkflow(ctx workflow.Context, idx int) (int, error) {
	logger := workflow.GetLogger(ctx)
	logger.Info("Consumer child workflow execution", "index", idx)

	ao := workflow.LocalActivityOptions{
		StartToCloseTimeout: 10 * time.Second,
	}
	ctx = workflow.WithLocalActivityOptions(ctx, ao)

	ret := 0
	err := workflow.ExecuteLocalActivity(ctx, Activity, idx).Get(ctx, ret)
	if err != nil {
		logger.Error("Activity failed.", "Error", err)
		return ret, err
	}
	return ret, nil
}

func Activity(ctx context.Context, idx int) (int, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("Activity execution", "index", idx)

	ret := 0
	for i := 0; i < maxMessage; i++ {
		msg, err := reader.ReadMessage(ctx)
		if err != nil {
			return ret, nil
		}
		fmt.Printf("Message on %s: %s\n", msg.Topic, string(msg.Value))
	}
	return ret, nil
}
