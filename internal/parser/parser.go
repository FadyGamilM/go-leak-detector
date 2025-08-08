package parser

import (
	"context"
	"regexp"
	"strings"
)

type GoroutinesStackParser struct {
	ctx context.Context
}

type GoroutineStackReport struct {
	Id                   string
	Status               string
	CreatedBy            string
	CreatedByGoroutineId string
	CreatedAtFilePath    string
	CreatedAtLine        string
	Stack                string
}

var (
	headerRe    = regexp.MustCompile(`^goroutine (\d+) \[([^\]]+)\]:`)
	createdByRe = regexp.MustCompile(`created by ([^\s]+)(?: in goroutine (\d+))?\n\s+([^\n]+):(\d+)`)
)

func New(ctx context.Context) *GoroutinesStackParser {
	return &GoroutinesStackParser{
		ctx: ctx,
	}
}

func (gsp *GoroutinesStackParser) Parse(buf []byte, lengthOfWrittenBytes int) []GoroutineStackReport {
	stackStr := string(buf[:lengthOfWrittenBytes])
	// indeices will have the start, end+1 index of each found pattern
	indices := headerRe.FindAllStringIndex(stackStr, -1)
	routinesReports := make([]GoroutineStackReport, len(indices))
	for i := range indices {
		start := indices[i][0]
		var end int
		if i+1 < len(indices) {
			end = indices[i+1][0]
		} else {
			end = len(stackStr)
		}
		m := stackStr[start:end]
		// each gorouutine in the stack printed by the runtime ends with two \n\n so i have to trim the suffix spaces before appending the results
		routinesReports[i] = gsp.ParseSingleRoutineStack(strings.TrimSuffix(m, "\n\n"))

	}
	return routinesReports
}

func (gsp *GoroutinesStackParser) ParseSingleRoutineStack(routine string) GoroutineStackReport {
	report := GoroutineStackReport{}

	lines := strings.SplitN(routine, "\n", 2)
	report.Id, report.Status = "", ""
	if m := headerRe.FindStringSubmatch(lines[0]); m != nil {
		report.Id = m[1]
		report.Status = m[2]
	}

	report.CreatedBy, report.CreatedByGoroutineId, report.CreatedAtFilePath, report.CreatedAtLine = "", "", "", ""
	if m := createdByRe.FindStringSubmatch(routine); m != nil {
		report.CreatedBy = m[1]
		report.CreatedByGoroutineId = m[2]
		report.CreatedAtFilePath = m[3]
		report.CreatedAtLine = m[4]
	}
	return report
}
