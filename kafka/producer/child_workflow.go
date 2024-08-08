package producer

import (
	"context"
	"github.com/pkg/errors"
	"github.com/segmentio/kafka-go"
	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/workflow"
	"strconv"
	"time"
)

var (
	TC     client.Client
	writer *kafka.Writer
)

func InitWriter(tEndpoint, kEndpoint, kTopic string) error {
	c, err := client.Dial(client.Options{
		HostPort: tEndpoint,
	})
	if err != nil {
		return errors.Wrap(err, "Unable to create client")
	}
	TC = c

	brokers := []string{kEndpoint} // Kafka 代理地址
	writer = kafka.NewWriter(kafka.WriterConfig{
		Brokers:  brokers,
		Topic:    kTopic,
		Balancer: &kafka.RoundRobin{},
	})
	return nil
}

func CloseWriter() {
	if writer != nil {
		writer.Close()
	}
	if TC != nil {
		TC.Close()
	}
}

// ChildWorkflow is a Workflow Definition
func ChildWorkflow(ctx workflow.Context, idx int) error {
	logger := workflow.GetLogger(ctx)
	logger.Info("Producer child workflow execution", "index", idx)

	ao := workflow.LocalActivityOptions{
		StartToCloseTimeout: 10 * time.Second,
	}
	ctx = workflow.WithLocalActivityOptions(ctx, ao)

	err := workflow.ExecuteLocalActivity(ctx, Activity, idx).Get(ctx, nil)
	if err != nil {
		logger.Error("Activity failed.", "Error", err)
		return err
	}
	return nil
}

func Activity(ctx context.Context, idx int) error {
	logger := activity.GetLogger(ctx)
	logger.Info("Activity execution", "index", idx)

	msg := kafka.Message{
		Value: []byte("temporal message" + strconv.Itoa(idx)),
	}
	return writer.WriteMessages(ctx, msg)
}
