package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/djadmin/fort/internal/checks"
)

const (
	colorReset  = "\033[0m"
	colorGreen  = "\033[32m"
	colorRed    = "\033[31m"
	colorYellow = "\033[33m"
	colorBold   = "\033[1m"
	colorDim    = "\033[2m"
)

type jsonReport struct {
	Tool      string       `json:"tool"`
	Version   string       `json:"version"`
	Hostname  string       `json:"hostname"`
	Serial    string       `json:"serial"`
	OSVersion string       `json:"os_version"`
	Timestamp string       `json:"timestamp"`
	Summary   jsonSummary  `json:"summary"`
	Policies  []jsonPolicy `json:"policies"`
}

// jsonPolicy wraps a Result with compliance framework control mappings.
type jsonPolicy struct {
	checks.Result
	Frameworks map[string][]string `json:"frameworks,omitempty"`
}

type jsonSummary struct {
	Total int    `json:"total"`
	Pass  int    `json:"pass"`
	Fail  int    `json:"fail"`
	Warn  int    `json:"warn"`
	Score string `json:"score"`
}

func printHuman(results []checks.Result, hostname, osVer string) {
	pass, fail, warn := tally(results)
	total := len(results)
	sep := strings.Repeat("─", 67)

	fmt.Printf("\n  %sfort v%s%s  —  %s (macOS %s)\n", colorBold, version, colorReset, hostname, osVer)
	fmt.Printf("  %s%s%s\n\n", colorDim, sep, colorReset)

	for _, r := range results {
		printRow(r)
	}

	fmt.Printf("\n  %s%s%s\n", colorDim, sep, colorReset)

	scoreColor := colorGreen + colorBold
	if fail > 0 {
		scoreColor = colorRed + colorBold
	} else if warn > 0 {
		scoreColor = colorYellow + colorBold
	}
	fmt.Printf("  %sScore: %d/%d%s  (%d pass, %d fail, %d warn)\n",
		scoreColor, pass, total, colorReset, pass, fail, warn)

	if anyFixable(results) {
		fmt.Printf("\n  %sRun fort --fix to remediate fixable issues.%s\n", colorDim, colorReset)
	}
	fmt.Println()
}

func printRow(r checks.Result) {
	icon, col := statusIcon(r.Status)
	name := fmt.Sprintf("%-26s", r.Name)
	current := fmt.Sprintf("%-28s", r.Current)

	fmt.Printf("  %s%s%s  %s %s", col, icon, colorReset, name, current)
	if r.Status != checks.StatusPass {
		fmt.Printf("%sexpected: %s%s", colorDim, r.Expected, colorReset)
	}
	if r.Fixed {
		fmt.Printf("  %s✓ fixed%s", colorGreen, colorReset)
	}
	fmt.Println()
}

func statusIcon(s checks.Status) (icon, color string) {
	switch s {
	case checks.StatusPass:
		return "✓", colorGreen
	case checks.StatusFail:
		return "✗", colorRed
	default:
		return "~", colorYellow
	}
}

func tally(results []checks.Result) (pass, fail, warn int) {
	for _, r := range results {
		switch r.Status {
		case checks.StatusPass:
			pass++
		case checks.StatusFail:
			fail++
		case checks.StatusWarn:
			warn++
		}
	}
	return
}

func anyFixable(results []checks.Result) bool {
	for _, r := range results {
		if r.Fixable && r.Status != checks.StatusPass {
			return true
		}
	}
	return false
}

func printJSON(results []checks.Result, hostname, serial, osVer string) {
	pass, fail, warn := tally(results)
	report := jsonReport{
		Tool:      "fort",
		Version:   version,
		Hostname:  hostname,
		Serial:    serial,
		OSVersion: osVer,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Summary: jsonSummary{
			Total: len(results),
			Pass:  pass,
			Fail:  fail,
			Warn:  warn,
			Score: fmt.Sprintf("%d/%d", pass, len(results)),
		},
		Policies: toJSONPolicies(results),
	}
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	_ = enc.Encode(report)
}

func toJSONPolicies(results []checks.Result) []jsonPolicy {
	out := make([]jsonPolicy, len(results))
	for i, r := range results {
		fw := checks.FrameworksFor(r.ID)
		fwMap := make(map[string][]string, len(fw))
		for _, f := range fw {
			fwMap[f.Name] = f.Controls
		}
		out[i] = jsonPolicy{Result: r, Frameworks: fwMap}
	}
	return out
}
