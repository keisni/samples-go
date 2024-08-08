package consumer

import (
	"context"
	"fmt"
	"time"

	"go.temporal.io/sdk/workflow"
)

const (
	maxMessage = 2
)

// ChildWorkflow is a Workflow Definition
func ChildWorkflow(ctx workflow.Context, idx int) (int, error) {
	fmt.Printf("Consumer child workflow execution, index:%d\n", idx)

	ao := workflow.LocalActivityOptions{
		StartToCloseTimeout: 10 * time.Second,
	}
	ctx = workflow.WithLocalActivityOptions(ctx, ao)

	var ret int
	err := workflow.ExecuteLocalActivity(ctx, Activity, idx).Get(ctx, &ret)
	if err != nil {
		fmt.Printf("Activity failed. %v\n", err)
		return ret, err
	}
	return ret, nil
}

func Activity(ctx context.Context, idx int) (int, error) {
	fmt.Printf("Activity execution, index:%d\n", idx)

	for i := 0; i < maxMessage; i++ {
		readerCtx, _ := context.WithTimeout(ctx, time.Second*2)
		msg, err := reader.ReadMessage(readerCtx)
		if err != nil {
			return i, nil
		}
		fmt.Printf("Message on %s: %s\n", msg.Topic, string(msg.Value))
	}
	return maxMessage, nil
}
