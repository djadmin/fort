package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/djadmin/fort/internal/checks"
)

var version = "0.1.0"

func main() {
	var (
		jsonOutput bool
		fix        bool
		dryRun     bool
		report     bool
	)

	root := &cobra.Command{
		Use:          "fort",
		Short:        "Endpoint security checker for macOS",
		Long:         "fort — audit, fix, and prove endpoint security. One command.",
		Version:      version,
		SilenceUsage: true,
	}

	root.Flags().BoolVar(&jsonOutput, "json", false, "output structured JSON")
	root.Flags().BoolVar(&fix, "fix", false, "auto-remediate fixable issues")
	root.Flags().BoolVar(&dryRun, "dry-run", false, "show what --fix would change without applying it")
	root.Flags().BoolVar(&report, "report", false, "write an HTML evidence report (fort-report.html)")

	var exitCode int
	root.RunE = func(cmd *cobra.Command, args []string) error {
		code, err := run(jsonOutput, fix, dryRun, report)
		exitCode = code
		return err
	}

	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
	os.Exit(exitCode)
}

// run returns (exitCode, error).
// exitCode: 0 = all pass, 1 = any fail, 2 = any warn and no fail.
// dry-run always returns 0 — it is purely informational.
func run(jsonOutput, fix, dryRun, report bool) (int, error) {
	allChecks := checks.All()
	if len(allChecks) == 0 {
		return 1, fmt.Errorf("no checks available for this platform")
	}

	if dryRun {
		return 0, runDryRun(allChecks)
	}

	results := make([]checks.Result, 0, len(allChecks))
	for _, c := range allChecks {
		r := c.Run()
		if fix && c.Fixable() && r.Status != checks.StatusPass {
			if err := c.Fix(); err != nil {
				fmt.Fprintf(os.Stderr, "  could not fix %s: %v\n", c.Name(), err)
			} else {
				r = c.Run()
				r.Fixed = true
			}
		}
		results = append(results, r)
	}

	h, osVer, serial := hostname(), osVersion(), serialNumber()

	if jsonOutput {
		printJSON(results, h, serial, osVer)
	} else {
		printHuman(results, h, osVer)
	}

	if report {
		ts := time.Now().Format("2006-01-02")
		outPath := fmt.Sprintf("fort-report-%s.html", ts)
		if err := writeReport(results, h, serial, osVer, outPath); err != nil {
			return 1, fmt.Errorf("report: %w", err)
		}
		fmt.Printf("  Report written to %s\n\n", outPath)
	}

	_, fail, warn := tally(results)
	switch {
	case fail > 0:
		return 1, nil
	case warn > 0:
		return 2, nil
	default:
		return 0, nil
	}
}

func runDryRun(allChecks []checks.Check) error {
	sep := strings.Repeat("─", 55)
	fmt.Printf("\n  %sfort v%s — dry run (nothing will be changed)%s\n", colorDim, version, colorReset)
	fmt.Printf("  %s%s%s\n\n", colorDim, sep, colorReset)

	fixable := 0
	for _, c := range allChecks {
		r := c.Run()
		if c.Fixable() && r.Status != checks.StatusPass {
			fixable++
			fmt.Printf("  %s✗%s  %s\n", colorRed, colorReset, c.Name())
			fmt.Printf("      current: %s\n", r.Current)
			fmt.Printf("      would run: %s%s%s\n\n", colorDim, c.FixDescription(), colorReset)
		}
	}

	if fixable == 0 {
		fmt.Printf("  %s✓  Nothing to fix — all checks pass.%s\n\n", colorGreen, colorReset)
		return nil
	}

	fmt.Printf("  %s%s%s\n", colorDim, sep, colorReset)
	fmt.Printf("  %d fixable issue(s) found. Run fort --fix to apply.\n\n", fixable)
	return nil
}

func hostname() string {
	h, err := os.Hostname()
	if err != nil {
		return "unknown"
	}
	return strings.TrimSuffix(h, ".local")
}

func osVersion() string {
	out, err := exec.Command("sw_vers", "-productVersion").Output()
	if err != nil {
		return "unknown"
	}
	return strings.TrimSpace(string(out))
}

func serialNumber() string {
	out, err := exec.Command("ioreg", "-c", "IOPlatformExpertDevice", "-d", "2").Output()
	if err != nil {
		return "unknown"
	}
	for _, line := range strings.Split(string(out), "\n") {
		if strings.Contains(line, "IOPlatformSerialNumber") {
			parts := strings.Split(line, `"`)
			if len(parts) >= 4 {
				return parts[3]
			}
		}
	}
	return "unknown"
}
