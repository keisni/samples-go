package producer

import (
	"fmt"
	enumspb "go.temporal.io/api/enums/v1"

	"go.temporal.io/sdk/workflow"
)

const (
	childNum = 200
)

// ParentWorkflow is a Workflow Definition
func ParentWorkflow(ctx workflow.Context) (processed int, err error) {
	logger := workflow.GetLogger(ctx)
	logger.Info("ParentWorkflow begin\n")

	var results []workflow.ChildWorkflowFuture
	for i := 0; i < childNum; i++ {
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
	for _, childResult := range results {
		err := childResult.Get(ctx, nil) // blocks until the child completion
		if err != nil {
			logger.Error("child execution failure.", "Error", err)
			continue
		}
		result += 1
	}
	logger.Info("ParentWorkflow finish.", "result", result)
	return result, nil
}
