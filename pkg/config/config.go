package config

import (
	"time"

	"github.com/fadygamilm/go-leak-detector/internal/parser"
)

// Config holds configuration for the leak detector.
type Config struct {
	MonitoringIntervalBetweenSnapshots time.Duration
	BufferSize                         int                                 // Buffer size for stack traces
	LeakThreshold                      int                                 // Number of snapshots a goroutine must persist to be considered a leak
	ExcludePatterns                    []string                            // Regex patterns to exclude goroutines (e.g., runtime goroutines)
	Callback                           func([]parser.GoroutineStackReport) // Optional callback for leak events
}
