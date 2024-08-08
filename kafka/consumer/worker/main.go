package main

import (
	"flag"
	"log"

	"go.temporal.io/sdk/worker"

	"github.com/temporalio/samples-go/kafka/consumer"
)

const (
	hostPort = "192.168.49.2:30880"
)

func main() {
	tEndpoint := flag.String("t_endpoint", "localhost:7233", "temporal endpoint")
	kEndpoint := flag.String("k_endpoint", "localhost:9092", "kafka endpoint")
	kTopic := flag.String("k_topic", "my_topic", "kafka topic")
	flag.Parse()

	if err := consumer.InitReader(*tEndpoint, *kEndpoint, *kTopic); err != nil {
		log.Fatalln("Init failed", err)
	}
	defer consumer.CloseReader()

	w := worker.New(consumer.TC, "producer", worker.Options{})

	w.RegisterWorkflow(consumer.ParentWorkflow)
	w.RegisterWorkflow(consumer.ChildWorkflow)
	w.RegisterActivity(consumer.Activity)

	err := w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("Unable to start worker", err)
	}
}
