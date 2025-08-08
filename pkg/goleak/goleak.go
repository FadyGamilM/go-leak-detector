package goleak

import (
	"context"
	"time"

	"github.com/FadyGamilM/go-leak-detector/internal/monitor"
	"github.com/FadyGamilM/go-leak-detector/pkg/config"
)

type GoLeakMonitor struct {
	monitor *monitor.Monitor
}

// New creates a new leak detector monitor with the given configuration.
func New(ctx context.Context, config config.Config) *GoLeakMonitor {
	if config.MonitoringIntervalBetweenSnapshots == 0 {
		config.MonitoringIntervalBetweenSnapshots = 5 * time.Second
	}
	if config.BufferSize == 0 {
		config.BufferSize = 20
	}
	if config.LeakThreshold == 0 {
		config.LeakThreshold = 3
	}
	return &GoLeakMonitor{
		monitor: monitor.New(ctx, config),
	}
}

func (m *GoLeakMonitor) Start() {
	m.monitor.StartMonitor()
}

func (m *GoLeakMonitor) Stop() {
	m.monitor.Stop()
}
