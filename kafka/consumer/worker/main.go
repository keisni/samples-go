package main

import (
	"fmt"
	"log"

	"go.temporal.io/sdk/worker"

	"github.com/temporalio/samples-go/kafka/consumer"
	"github.com/temporalio/samples-go/kafka/helper"
)

const (
	hostPort = "192.168.49.2:30880"
)

func main() {
	opts := helper.ParseOptions()
	fmt.Printf("opts: %v\n", opts)

	if err := consumer.InitReader(opts); err != nil {
		log.Fatalln("Init failed", err)
	}
	defer consumer.CloseReader()

	w := worker.New(consumer.TC, "consumer", worker.Options{})

	w.RegisterWorkflow(consumer.ParentWorkflow)
	w.RegisterWorkflow(consumer.ChildWorkflow)
	w.RegisterActivity(consumer.Activity)

	err := w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("Unable to start worker", err)
	}
}
