package leakdetector

import (
	"context"
	"time"

	"github.com/fadygamilm/go-leak-detector/internal/parser"
)

// simply this pkg will track the go routines from the current run to the previous run, and if we found the same goroutine in the previous snapshot run with same status, we will consider it a lean (a stuck go routine) and we will report it.
type LeakDetector struct {
	ctx           context.Context
	PrevSnapshots map[string][]parser.GoroutineStackReport
	LeakThreshold int
	FirstSeen     map[string]time.Time
}

func New(ctx context.Context, leakThreshold int) *LeakDetector {
	return &LeakDetector{
		ctx:           ctx,
		PrevSnapshots: make(map[string][]parser.GoroutineStackReport),
		LeakThreshold: leakThreshold,
		FirstSeen:     make(map[string]time.Time),
	}
}

func (ld *LeakDetector) DetectLeaks(goroutinesReports []parser.GoroutineStackReport) ([]parser.GoroutineStackReport, map[string]time.Duration) {
	leakedRoutines := []parser.GoroutineStackReport{}
	currentSnapshot := make(map[string]parser.GoroutineStackReport, len(goroutinesReports))
	leaksDurations := make(map[string]time.Duration)

	for _, report := range goroutinesReports {
		currentSnapshot[report.Id] = report

		// If this is the first time seeing the goroutine, record the timestamp
		if _, exists := ld.FirstSeen[report.Id]; !exists {
			ld.FirstSeen[report.Id] = time.Now()
		}

		ld.PrevSnapshots[report.Id] = append(ld.PrevSnapshots[report.Id], report)

		if len(ld.PrevSnapshots[report.Id]) >= ld.LeakThreshold {
			isLeak := true // its a leaked routine unless it's proven to be not so
			for i := 1; i < len(ld.PrevSnapshots[report.Id]); i++ {
				// so once i will find a 2 consutive reports with different status, i will consider it not a leak
				if ld.PrevSnapshots[report.Id][i].Status != ld.PrevSnapshots[report.Id][i-1].Status {
					isLeak = false
					break
				}
			}
			if isLeak {
				leakedRoutines = append(leakedRoutines, report)
				leaksDurations[report.Id] = time.Since(ld.FirstSeen[report.Id])
			}
		}
	}

	// i will clean the previous snapshots to keep only the latest snapshot found in the current to avoid memory leaks if the program runs for a long time
	for id := range ld.PrevSnapshots {
		if _, exists := currentSnapshot[id]; !exists {
			delete(ld.PrevSnapshots, id)
			delete(ld.FirstSeen, id)
		}

	}

	return leakedRoutines, leaksDurations
}
