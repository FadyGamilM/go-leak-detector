package monitor

import (
	"context"
	"log"
	"runtime"
	"time"

	leakdetector "github.com/FadyGamilM/go-leak-detector/internal/leak_detector"
	"github.com/FadyGamilM/go-leak-detector/internal/parser"
	"github.com/FadyGamilM/go-leak-detector/pkg/config"
	"github.com/prometheus/client_golang/prometheus"
)

type Metrics struct {
	GoroutineCount    prometheus.Gauge
	LeakCount         prometheus.Gauge
	LeakDuration      prometheus.Histogram
	LeakDetectionTime prometheus.Summary
}
type Monitor struct {
	ctx      context.Context
	parser   *parser.GoroutinesStackParser
	detector *leakdetector.LeakDetector
	config   config.Config
	cancel   context.CancelFunc
	metrics  *Metrics
	registry *prometheus.Registry
}

func New(ctx context.Context, config config.Config) *Monitor {
	ctx, cancel := context.WithCancel(ctx)
	registry := prometheus.NewRegistry()
	metrics := &Metrics{
		GoroutineCount: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "goleak_goroutine_count",
			Help: "Total number of goroutines detected.",
		}),
		LeakCount: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "goleak_leak_count",
			Help: "Number of detected goroutine leaks.",
		}),
		LeakDuration: prometheus.NewHistogram(prometheus.HistogramOpts{
			Name:    "goleak_leak_duration_seconds",
			Help:    "Duration of detected goroutine leaks.",
			Buckets: prometheus.LinearBuckets(0, 10, 10),
		}),
		LeakDetectionTime: prometheus.NewSummary(prometheus.SummaryOpts{
			Name: "goleak_detection_time_seconds",
			Help: "Time taken to detect leaks.",
		}),
	}
	registry.MustRegister(metrics.GoroutineCount, metrics.LeakCount, metrics.LeakDuration, metrics.LeakDetectionTime)
	return &Monitor{
		ctx:      ctx,
		parser:   parser.New(ctx, config.ExcludePatterns),
		detector: leakdetector.New(ctx, config.LeakThreshold),
		config:   config,
		cancel:   cancel,
		metrics:  metrics,
		registry: registry,
	}
}

func (m *Monitor) StartMonitor() {
	ticker := time.NewTicker(m.config.MonitoringIntervalBetweenSnapshots)

	for {
		select {
		case <-ticker.C:
			start := time.Now()

			buf := make([]byte, 1<<m.config.BufferSize)
			n := runtime.Stack(buf, true)
			// log.Println("======> ", string(buf[:n]))

			reports := m.parser.Parse(buf, n)
			m.metrics.GoroutineCount.Set(float64(len(reports)))

			leaks, leakDurations := m.detector.DetectLeaks(reports)
			m.metrics.LeakCount.Set(float64(len(leaks)))

			for _, leak := range leaks {
				if duration, exists := leakDurations[leak.Id]; exists {
					m.metrics.LeakDuration.Observe(duration.Seconds())
				}
			}
			m.metrics.LeakDetectionTime.Observe(time.Since(start).Seconds())

			if len(leaks) > 0 {
				log.Printf("⚠️ Found %d goroutine leaks:\n", len(leaks))
				for _, leak := range leaks {
					log.Printf("Goroutine ID: %s\nStatus: %s\nCreated By: %s (Goroutine ID: %s)\nCreated At: %s:%s\n\n\n",
						leak.Id,
						leak.Status,
						leak.CreatedBy,
						leak.CreatedByGoroutineId,
						leak.CreatedAtFilePath,
						leak.CreatedAtLine,
					)
				}
			}
		case <-m.ctx.Done():
			ticker.Stop()
			return
		}
	}
}

func (m *Monitor) Stop() {
	m.cancel()
}

func (m *Monitor) Registry() *prometheus.Registry {
	return m.registry
}
