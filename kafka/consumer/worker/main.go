package main

import (
	"log"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"

	"github.com/temporalio/samples-go/kafka/consumer"
)

const (
	hostPort = "192.168.49.2:30880"
)

func main() {
	// The client and worker are heavyweight objects that should be created once per process.
	c, err := client.Dial(client.Options{
		HostPort: hostPort,
	})
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer c.Close()
	consumer.InitReader()
	defer consumer.CloseReader()

	w := worker.New(c, "producer", worker.Options{})

	w.RegisterWorkflow(consumer.ParentWorkflow)
	w.RegisterWorkflow(consumer.ChildWorkflow)

	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("Unable to start worker", err)
	}
}
