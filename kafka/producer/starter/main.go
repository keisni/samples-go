package main

import (
	"context"
	"flag"
	"log"

	"go.temporal.io/sdk/client"

	"github.com/temporalio/samples-go/kafka/producer"
)

func main() {
	tEndpoint := flag.String("t_endpoint", "localhost:7233", "temporal endpoint")
	flag.Parse()

	// The client is a heavyweight object that should be created once per process.
	c, err := client.Dial(client.Options{
		HostPort: *tEndpoint,
	})
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer c.Close()

	// This workflow ID can be user business logic identifier as well.
	workflowID := "producer-parent-workflow"
	workflowOptions := client.StartWorkflowOptions{
		ID:           workflowID,
		TaskQueue:    "producer",
		CronSchedule: "* * * * *",
	}

	we, err := c.ExecuteWorkflow(context.Background(), workflowOptions, producer.ParentWorkflow)
	if err != nil {
		log.Fatalln("Unable to execute workflow", err)
	}
	log.Println("Started workflow", "WorkflowID", we.GetID(), "RunID", we.GetRunID())
}
