package producer

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/segmentio/kafka-go"
	"go.temporal.io/sdk/workflow"
)

// ChildWorkflow is a Workflow Definition
func ChildWorkflow(ctx workflow.Context, idx int) (int, error) {
	fmt.Printf("Producer child workflow execution, index:%d\n", idx)

	ao := workflow.LocalActivityOptions{
		StartToCloseTimeout: 10 * time.Second,
	}
	ctx = workflow.WithLocalActivityOptions(ctx, ao)

	var ret int
	err := workflow.ExecuteLocalActivity(ctx, Activity, idx).Get(ctx, &ret)
	if err != nil {
		fmt.Printf("Activity failed. %v\n", err)
		return 0, err
	}
	return 1, nil
}

func Activity(ctx context.Context, idx int) (ret int, err error) {
	fmt.Printf("Activity execution, index:%d\n", idx)

	msg := kafka.Message{
		Value: []byte("temporal message" + strconv.Itoa(idx)),
	}
	if err = writer.WriteMessages(ctx, msg); err != nil {
		return
	}
	ret = 1
	return
}
