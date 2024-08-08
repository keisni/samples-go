package main

import (
	"fmt"
	"log"

	"go.temporal.io/sdk/worker"

	"github.com/temporalio/samples-go/kafka/helper"
	"github.com/temporalio/samples-go/kafka/producer"
)

func main() {
	opts := helper.ParseOptions()
	fmt.Printf("opts: %v\n", opts)

	if err := producer.InitWriter(opts); err != nil {
		log.Fatalln("Init failed", err)
	}
	defer producer.CloseWriter()

	w := worker.New(producer.TC, "producer", worker.Options{})

	w.RegisterWorkflow(producer.ParentWorkflow)
	w.RegisterWorkflow(producer.ChildWorkflow)
	w.RegisterActivity(producer.Activity)

	err := w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("Unable to start worker", err)
	}
}
