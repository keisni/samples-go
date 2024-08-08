package main

import (
	"flag"
	"log"

	"go.temporal.io/sdk/worker"

	"github.com/temporalio/samples-go/kafka/producer"
)

func main() {
	tEndpoint := flag.String("t_endpoint", "localhost:7233", "temporal endpoint")
	kEndpoint := flag.String("k_endpoint", "localhost:9092", "kafka endpoint")
	kTopic := flag.String("k_topic", "my_topic", "kafka topic")
	flag.Parse()

	if err := producer.InitWriter(*tEndpoint, *kEndpoint, *kTopic); err != nil {
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
