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
2025/08/09 01:21:55 ⚠️ Found 4 goroutine leaks:
2025/08/09 01:21:55 Goroutine ID: 1
Status: select (no cases)
Created By:  (Goroutine ID: )
Created At: :


2025/08/09 01:21:55 Goroutine ID: 44
Status: chan receive
Created By: main.main (Goroutine ID: 1)
Created At: /Users/fady/projects/sandbox/go-leak-detector/cmd/example/main.go:32


2025/08/09 01:21:55 Goroutine ID: 45
Status: chan receive
Created By: main.main (Goroutine ID: 1)
Created At: /Users/fady/projects/sandbox/go-leak-detector/cmd/example/main.go:33


2025/08/09 01:21:55 Goroutine ID: 50
Status: chan receive
Created By: github.com/FadyGamilM/go-leak-detector/pkg/goleak.(*GoLeakMonitor).StartPrometheusServer (Goroutine ID: 42)
Created At: /Users/fady/projects/sandbox/go-leak-detector/pkg/goleak/prometheus.go:15


2025/08/09 01:22:00 ⚠️ Found 4 goroutine leaks:
2025/08/09 01:22:00 Goroutine ID: 1
Status: select (no cases)
Created By:  (Goroutine ID: )
Created At: :


2025/08/09 01:22:00 Goroutine ID: 44
Status: chan receive
Created By: main.main (Goroutine ID: 1)
Created At: /Users/fady/projects/sandbox/go-leak-detector/cmd/example/main.go:32


2025/08/09 01:22:00 Goroutine ID: 45
Status: chan receive
Created By: main.main (Goroutine ID: 1)
Created At: /Users/fady/projects/sandbox/go-leak-detector/cmd/example/main.go:33


2025/08/09 01:22:00 Goroutine ID: 50
Status: chan receive
Created By: github.com/FadyGamilM/go-leak-detector/pkg/goleak.(*GoLeakMonitor).StartPrometheusServer (Goroutine ID: 42)
Created At: /Users/fady/projects/sandbox/go-leak-detector/pkg/goleak/prometheus.go:15
```

## Example of the `metrics/` endpoint setuped to be used for prometheus :
```shell
 curl http://localhost:9091/metrics
# HELP goleak_detection_time_seconds Time taken to detect leaks.
# TYPE goleak_detection_time_seconds summary
goleak_detection_time_seconds_sum 0.007834749
goleak_detection_time_seconds_count 6
# HELP goleak_goroutine_count Total number of goroutines detected.
# TYPE goleak_goroutine_count gauge
goleak_goroutine_count 4
# HELP goleak_leak_count Number of detected goroutine leaks.
# TYPE goleak_leak_count gauge
goleak_leak_count 4
# HELP goleak_leak_duration_seconds Duration of detected goroutine leaks.
# TYPE goleak_leak_duration_seconds histogram
goleak_leak_duration_seconds_bucket{le="0"} 0
goleak_leak_duration_seconds_bucket{le="10"} 0
goleak_leak_duration_seconds_bucket{le="20"} 12
goleak_leak_duration_seconds_bucket{le="30"} 16
goleak_leak_duration_seconds_bucket{le="40"} 16
goleak_leak_duration_seconds_bucket{le="50"} 16
goleak_leak_duration_seconds_bucket{le="60"} 16
goleak_leak_duration_seconds_bucket{le="70"} 16
goleak_leak_duration_seconds_bucket{le="80"} 16
goleak_leak_duration_seconds_bucket{le="90"} 16
goleak_leak_duration_seconds_bucket{le="+Inf"} 16
goleak_leak_duration_seconds_sum 279.985597206
goleak_leak_duration_seconds_count 16
```