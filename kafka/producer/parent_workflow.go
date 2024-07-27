package producer

import (
	"fmt"
	enumspb "go.temporal.io/api/enums/v1"

	"go.temporal.io/sdk/workflow"
)

// ParentWorkflow is a Workflow Definition
func ParentWorkflow(ctx workflow.Context) error {
	logger := workflow.GetLogger(ctx)

	for i := 0; i < 200; i++ {
		childID := fmt.Sprintf("producer_child_workflow:%d", i)
		cwo := workflow.ChildWorkflowOptions{
			WorkflowID:            childID,
			WorkflowIDReusePolicy: enumspb.WORKFLOW_ID_REUSE_POLICY_TERMINATE_IF_RUNNING,
		}

		ctx = workflow.WithChildOptions(ctx, cwo)
		var result string
		err := workflow.ExecuteChildWorkflow(ctx, ChildWorkflow, i).Get(ctx, &result)
		if err != nil {
			logger.Error("Parent execution received child execution failure.", "Error", err)
			return err
		}
	}
	logger.Info("Producer parent execution completed.")
	return nil
}
