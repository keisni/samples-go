package main

import (
	"context"
	"log"

	"go.temporal.io/sdk/client"

	"github.com/temporalio/samples-go/kafka/consumer"
)

const (
	hostPort = "192.168.49.2:30880"
)

func main() {
	// The client is a heavyweight object that should be created once per process.
	c, err := client.Dial(client.Options{
		HostPort: hostPort,
	})
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer c.Close()

	// This workflow ID can be user business logic identifier as well.
	workflowID := "consumer-parent-workflow"
	workflowOptions := client.StartWorkflowOptions{
		ID:           workflowID,
		TaskQueue:    "consumer",
		CronSchedule: "*/2 * * * *",
	}

	we, err := c.ExecuteWorkflow(context.Background(), workflowOptions, consumer.ParentWorkflow)
	if err != nil {
		log.Fatalln("Unable to execute workflow", err)
	}
	log.Println("Started workflow", "WorkflowID", we.GetID(), "RunID", we.GetRunID())
}
