package parser

import (
	"context"
	"regexp"
	"strings"
)

type GoroutinesStackParser struct {
	ctx               context.Context
	execludedPatterns []*regexp.Regexp
}

type GoroutineStackReport struct {
	Id                   string
	Status               string
	CreatedBy            string
	CreatedByGoroutineId string
	CreatedAtFilePath    string
	CreatedAtLine        string
	FullStackString      string
}

var (
	headerRe               = regexp.MustCompile(`(?m)^goroutine (\d+) \[([^\]]+)\]:`)
	createdByRe            = regexp.MustCompile(`created by ([^\s]+)(?: in goroutine (\d+))?\n\s+([^\n]+):(\d+)`)
	defaultExcludePatterns = []string{
		`^runtime\.`,
		`^syscall\.`,
		`\bsyscall\b`,
		`github\.com/FadyGamilM/go-leak-detector/internal/monitor`,
		`github\.com/prometheus/client_golang`,
	}
)

func New(ctx context.Context, excludePatterns []string) *GoroutinesStackParser {
	var compiledPatterns []*regexp.Regexp
	allPatterns := append(defaultExcludePatterns, excludePatterns...)
	for _, pattern := range allPatterns {
		if re, err := regexp.Compile(pattern); err == nil {
			compiledPatterns = append(compiledPatterns, re)
		}
	}
	return &GoroutinesStackParser{
		ctx:               ctx,
		execludedPatterns: compiledPatterns,
	}
}

func (gsp *GoroutinesStackParser) Parse(buf []byte, lengthOfWrittenBytes int) []GoroutineStackReport {
	stackStr := string(buf[:lengthOfWrittenBytes])
	// log.Println("the stack trace is: ", stackStr)
	// indeices will have the start, end+1 index of each found pattern
	indices := headerRe.FindAllStringIndex(stackStr, -1)
	routinesReports := []GoroutineStackReport{}
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
		trimmedM := strings.TrimSuffix(m, "\n\n")

		if !gsp.shouldInclude(trimmedM) {
			continue
		}

		report := gsp.ParseSingleRoutineStack(trimmedM)
		if report.Id != "" {
			routinesReports = append(routinesReports, report)
		}
	}
	return routinesReports
}

func (gsp *GoroutinesStackParser) shouldInclude(routine string) bool {
	for _, re := range gsp.execludedPatterns {
		if re.MatchString(routine) {
			return false
		}
	}

	return true
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
