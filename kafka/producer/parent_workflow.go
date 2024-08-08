package producer

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/segmentio/kafka-go"
	"github.com/temporalio/samples-go/kafka/helper"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/workflow"

	enumspb "go.temporal.io/api/enums/v1"
)

var (
	TC        client.Client
	writer    *kafka.Writer
	testCount int
)

func InitWriter(opts *helper.Options) error {
	c, err := client.Dial(client.Options{
		HostPort: opts.TemporalEndpoint,
	})
	if err != nil {
		return errors.Wrap(err, "Unable to create client")
	}
	TC = c

	brokers := []string{opts.KafkaEndpoint}
	writer = kafka.NewWriter(kafka.WriterConfig{
		Brokers:  brokers,
		Topic:    opts.KafkaTopic,
		Balancer: &kafka.RoundRobin{},
	})
	testCount = opts.Count
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

// ParentWorkflow is a Workflow Definition
func ParentWorkflow(ctx workflow.Context) (processed int, err error) {
	fmt.Printf("ParentWorkflow begin\n")

	var results []workflow.ChildWorkflowFuture
	for i := 0; i < testCount; i++ {
		childID := fmt.Sprintf("producer_child_workflow:%d", i)
		cwo := workflow.ChildWorkflowOptions{
			WorkflowID:            childID,
			WorkflowIDReusePolicy: enumspb.WORKFLOW_ID_REUSE_POLICY_TERMINATE_IF_RUNNING,
		}

		ctx = workflow.WithChildOptions(ctx, cwo)
		child := workflow.ExecuteChildWorkflow(ctx, ChildWorkflow, i)
		results = append(results, child)
	}
	// Waits for all child workflows to complete
	result := 0
	for _, child := range results {
		var ret int
		if childErr := child.Get(ctx, &ret); childErr != nil {
			fmt.Printf("child execution failure. %v\n", childErr)
			continue
		}
		result += ret
	}
	fmt.Printf("ParentWorkflow finish. result:%d\n", result)
	return result, nil
}
