package main

import (
	"context"
	"time"

	"github.com/fadygamilm/go-leak-detector/pkg/config"
	"github.com/fadygamilm/go-leak-detector/pkg/goleak"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	config := config.Config{
		MonitoringIntervalBetweenSnapshots: 5 * time.Second,
		BufferSize:                         20,
		LeakThreshold:                      3,
		ExcludePatterns: []string{
			`net/http\..*Server`,
			`main\.main.*select`,
		},
	}

	monitor := goleak.New(ctx, config)

	go monitor.StartPrometheusServer(ctx, ":9091")

	go monitor.Start()

	// Simulate a leak
	go func() { ch := make(chan struct{}); <-ch }()
	go func() { ch := make(chan struct{}); <-ch }()
	time.Sleep(1 * time.Second)
	select {}
}
