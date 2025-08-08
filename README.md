## Overview
Go Leak Detector continuously monitors your application's goroutines and identifies potential memory leaks by tracking goroutines that remain in the same state across multiple snapshots. It provides real-time detection with detailed reporting and Prometheus metrics integration.

## Features
- Real-time Detection: Continuously monitors goroutines with configurable intervals
- Prometheus Integration: Built-in metrics for monitoring and alerting
- Configurable Thresholds: Customizable leak detection sensitivity
- Pattern Exclusion: Filter out known long-running goroutines (servers, background workers)
- Detailed Reporting: Rich context including goroutine ID, status, creation location
- Production Ready: Minimal overhead with efficient stack trace parsing

## Installation
To install Go Leak Detector, use the following command:
```bash
go get github.com/FadyGamilM/go-leak-detector@v1.0.1
``` 

## Usage
```go
package main

import (
    "context"
    "time"

    "github.com/FadyGamilM/go-leak-detector/pkg/config"
    "github.com/FadyGamilM/go-leak-detector/pkg/goleak"
)

func main() {
    ctx := context.Background()
    
    cfg := config.Config{
        MonitoringIntervalBetweenSnapshots: 5 * time.Second,
        BufferSize:                         20,
        LeakThreshold:                      3,
        // Optional: patterns to exclude known long-running goroutines
        // such as HTTP servers or background workers
        // Adjust these patterns based on your application structure
        ExcludePatterns: []string{
            `net/http\..*Server`,
            `your-app/internal/background`,
        },
    }

    monitor := goleak.New(ctx, cfg)
    
    // Start Prometheus metrics server (Optional) 
    go monitor.StartPrometheusServer(ctx, ":9091")
    
    // Start leak detection
    go monitor.Start()
    
    // Your application code here
}
```

## Configuration

| Parameter                        | Description                                 | Default      |
|-----------------------------------|---------------------------------------------|--------------|
| MonitoringIntervalBetweenSnapshots| Time between goroutine snapshots            | 5s           |
| BufferSize                        | Stack trace buffer size (2^n bytes)         | 20 (1MB)     |
| LeakThreshold                     | Snapshots before flagging as leak           | 3            |
| ExcludePatterns                   | Regex patterns to exclude from detection    | `[]`         |

## Metrics

The detector exposes the following Prometheus metrics at `/metrics`:

| Metric Name                        | Description                              |
|-------------------------------------|------------------------------------------|
| `goleak_goroutine_count`            | Total number of tracked goroutines       |
| `goleak_leak_count`                 | Number of detected leaks                 |
| `goleak_leak_duration_seconds`      | Histogram of leak durations              |
| `goleak_detection_time_seconds`     | Detection processing time                |


## Example of the main.go output logs: 
```shell
 go run cmd/example/main.go
2025/08/09 00:15:39 ⚠️ Found 4 goroutine leaks:
2025/08/09 00:15:39 Goroutine ID: 1
Status: select (no cases)
Created By:  (Goroutine ID: )
Created At: :


2025/08/09 00:15:39 Goroutine ID: 27
Status: chan receive
Created By: main.main (Goroutine ID: 1)
Created At: /Users/fady/projects/sandbox/go-leak-detector/cmd/example/main.go:32


2025/08/09 00:15:39 Goroutine ID: 28
Status: chan receive
Created By: main.main (Goroutine ID: 1)
Created At: /Users/fady/projects/sandbox/go-leak-detector/cmd/example/main.go:33


2025/08/09 00:15:39 Goroutine ID: 29
Status: chan receive
Created By: github.com/fadygamilm/go-leak-detector/pkg/goleak.(*GoLeakMonitor).StartPrometheusServer (Goroutine ID: 25)
Created At: /Users/fady/projects/sandbox/go-leak-detector/pkg/goleak/prometheus.go:15


2025/08/09 00:15:44 ⚠️ Found 4 goroutine leaks:
2025/08/09 00:15:44 Goroutine ID: 1
Status: select (no cases)
Created By:  (Goroutine ID: )
Created At: :


2025/08/09 00:15:44 Goroutine ID: 27
Status: chan receive
Created By: main.main (Goroutine ID: 1)
Created At: /Users/fady/projects/sandbox/go-leak-detector/cmd/example/main.go:32


2025/08/09 00:15:44 Goroutine ID: 28
Status: chan receive
Created By: main.main (Goroutine ID: 1)
Created At: /Users/fady/projects/sandbox/go-leak-detector/cmd/example/main.go:33


2025/08/09 00:15:44 Goroutine ID: 29
Status: chan receive
Created By: github.com/fadygamilm/go-leak-detector/pkg/goleak.(*GoLeakMonitor).StartPrometheusServer (Goroutine ID: 25)
Created At: /Users/fady/projects/sandbox/go-leak-detector/pkg/goleak/prometheus.go:15


2025/08/09 00:15:49 ⚠️ Found 4 goroutine leaks:
2025/08/09 00:15:49 Goroutine ID: 1
Status: select (no cases)
Created By:  (Goroutine ID: )
Created At: :


2025/08/09 00:15:49 Goroutine ID: 27
Status: chan receive
Created By: main.main (Goroutine ID: 1)
Created At: /Users/fady/projects/sandbox/go-leak-detector/cmd/example/main.go:32


2025/08/09 00:15:49 Goroutine ID: 28
Status: chan receive
Created By: main.main (Goroutine ID: 1)
Created At: /Users/fady/projects/sandbox/go-leak-detector/cmd/example/main.go:33


2025/08/09 00:15:49 Goroutine ID: 29
Status: chan receive
Created By: github.com/fadygamilm/go-leak-detector/pkg/goleak.(*GoLeakMonitor).StartPrometheusServer (Goroutine ID: 25)
Created At: /Users/fady/projects/sandbox/go-leak-detector/pkg/goleak/prometheus.go:15
```