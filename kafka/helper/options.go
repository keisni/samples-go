package helper

import (
	"flag"
)

type Options struct {
	TemporalEndpoint string
	KafkaEndpoint    string
	KafkaTopic       string
	Count            int
}

func ParseOptions() *Options {
	opts := Options{}

	flag.StringVar(&opts.TemporalEndpoint, "t_endpoint", "localhost:7233", "temporal endpoint")
	flag.StringVar(&opts.KafkaEndpoint, "k_endpoint", "localhost:9092", "kafka endpoint")
	flag.StringVar(&opts.KafkaTopic, "k_topic", "my-topic", "kafka topic")
	flag.IntVar(&opts.Count, "count", 1, "batch count")
	flag.Parse()

	return &opts
}
