package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	survey "github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/terminal"
	"github.com/spf13/cobra"

	"github.com/djadmin/fort/internal/checks"
)

var version = "0.3.1"

func main() {
	var (
		jsonOutput bool
		fix        bool
		dryRun     bool
		report     bool
		yes        bool
		only       string
	)

	root := &cobra.Command{
		Use:          "fort",
		Short:        "Endpoint security checker for macOS",
		Long:         "fort: audit, fix, and prove endpoint security. One command.",
		Version:      version,
		SilenceUsage: true,
	}

	root.Flags().BoolVar(&jsonOutput, "json", false, "output structured JSON")
	root.Flags().BoolVar(&fix, "fix", false, "remediate fixable issues (shows confirmation prompt)")
	root.Flags().BoolVar(&dryRun, "dry-run", false, "show what --fix would change without applying it")
	root.Flags().BoolVar(&report, "report", false, "write an HTML evidence report (fort-report.html)")
	root.Flags().BoolVarP(&yes, "yes", "y", false, "skip confirmation prompt when using --fix")
	root.Flags().StringVar(&only, "only", "", "run only these checks (comma-separated IDs, e.g. filevault,firewall)")

	var exitCode int
	root.RunE = func(cmd *cobra.Command, args []string) error {
		code, err := run(jsonOutput, fix, dryRun, report, yes, only)
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
func run(jsonOutput, fix, dryRun, report, yes bool, only string) (int, error) {
	latestVersion := checkLatestVersion()
	allChecks := checks.All()
	if len(allChecks) == 0 {
		return 1, fmt.Errorf("no checks available for this platform")
	}

	// Filter to requested IDs when --only is set.
	if only != "" {
		wanted := map[string]bool{}
		for _, id := range strings.Split(only, ",") {
			wanted[strings.TrimSpace(id)] = true
		}
		filtered := allChecks[:0]
		for _, c := range allChecks {
			if wanted[c.ID()] {
				filtered = append(filtered, c)
				delete(wanted, c.ID())
			}
		}
		for unknown := range wanted {
			fmt.Fprintf(os.Stderr, "  warning: unknown check ID %q — skipped\n", unknown)
		}
		allChecks = filtered
		if len(allChecks) == 0 {
			return 1, fmt.Errorf("no matching checks for --only %q", only)
		}
	}

	if dryRun {
		return 0, runDryRun(allChecks)
	}

	// First pass: run all checks without fixing anything.
	results := make([]checks.Result, len(allChecks))
	for i, c := range allChecks {
		results[i] = c.Run()
	}

	h, osVer, serial := hostname(), osVersion(), serialNumber()

	// Show current state before any fixes.
	if !jsonOutput {
		printHuman(results, h, osVer)
	}

	// If --fix, collect fixable failures and confirm before applying.
	if fix {
		type pending struct {
			check  checks.Check
			result checks.Result
			idx    int
		}
		var fixable []pending
		for i, c := range allChecks {
			if c.Fixable() && results[i].Status != checks.StatusPass {
				fixable = append(fixable, pending{c, results[i], i})
			}
		}

		if len(fixable) == 0 {
			fmt.Printf("  %s✓  Nothing to fix.%s\n\n", colorGreen, colorReset)
		} else {
			// Show what will be changed and ask for confirmation unless --yes.
			if !yes && !jsonOutput {
				sep := strings.Repeat("─", 67)
				noun := "fix"
				if len(fixable) != 1 {
					noun = "fixes"
				}
				fmt.Printf("\n  %s%s%s\n", colorDim, sep, colorReset)
				fmt.Printf("  %s%d %s to apply:%s\n\n", colorBold, len(fixable), noun, colorReset)
				for i, p := range fixable {
					fmt.Printf("  %s%d.%s %s%s%s\n",
						colorDim, i+1, colorReset,
						colorBold, p.check.Name(), colorReset)
					fmt.Printf("     %scurrent:%s %-16s %s→%s %s%s%s\n",
						colorDim, colorReset, p.result.Current,
						colorDim, colorReset,
						colorGreen, p.result.Expected, colorReset)
					if desc := p.check.FixDescription(); desc != "" {
						fmt.Printf("     %s$ %s%s\n", colorDim, desc, colorReset)
					}
					fmt.Println()
				}
				// Build labelled options — all selected by default.
				opts := make([]string, len(fixable))
				for i, p := range fixable {
					opts[i] = fmt.Sprintf("%-26s  %s → %s", p.check.Name(), p.result.Current, p.result.Expected)
				}

				var chosen []string
				prompt := &survey.MultiSelect{
					Message: "Select fixes to apply",
					Options: opts,
					Default: opts,
				}
				err := survey.AskOne(prompt, &chosen, survey.WithIcons(func(icons *survey.IconSet) {
					icons.MarkedOption.Text = "✓"
					icons.MarkedOption.Format = "green+b"
					icons.UnmarkedOption.Text = "○"
					icons.UnmarkedOption.Format = "default+d"
					icons.SelectFocus.Text = "❯"
					icons.SelectFocus.Format = "green+b"
				}))
				if err == terminal.InterruptErr || len(chosen) == 0 {
					fmt.Printf("\n  %sNo changes made.%s\n\n", colorDim, colorReset)
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

				// Keep only what the user selected.
				chosenSet := map[string]bool{}
				for _, c := range chosen {
					chosenSet[c] = true
				}
				filtered := fixable[:0]
				for i, p := range fixable {
					if chosenSet[opts[i]] {
						filtered = append(filtered, p)
					}
				}
				fixable = filtered
				fmt.Println()
			}

			// Apply fixes.
			for _, p := range fixable {
				if err := p.check.Fix(); err != nil {
					fmt.Fprintf(os.Stderr, "  could not fix %s: %v\n", p.check.Name(), err)
				} else {
					r := p.check.Run()
					r.Fixed = true
					results[p.idx] = r
				}
			}

			// Re-display results after fixes.
			if !jsonOutput {
				printHuman(results, h, osVer)
			}
		}
	}

	if jsonOutput {
		printJSON(results, h, serial, osVer)
	}

	if report {
		ts := time.Now().Format("2006-01-02")
		outPath := fmt.Sprintf("fort-report-%s.html", ts)
		if err := writeReport(results, h, serial, osVer, outPath); err != nil {
			return 1, fmt.Errorf("report: %w", err)
		}
		fmt.Printf("  Report written to %s\n\n", outPath)
	}

	if latest := <-latestVersion; latest != "" && latest != version {
		fmt.Printf("  %s↑  fort v%s is available (you have v%s) — brew upgrade djadmin/tap/fort%s\n\n",
			colorDim, latest, version, colorReset)
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
	sep := strings.Repeat("─", 67)
	fmt.Printf("\n  %sfort v%s · dry run (nothing will be changed)%s\n", colorDim, version, colorReset)
	fmt.Printf("  %s%s%s\n\n", colorDim, sep, colorReset)

	fixable := 0
	for _, c := range allChecks {
		r := c.Run()
		if c.Fixable() && r.Status != checks.StatusPass {
			fixable++
			fmt.Printf("  %s✗%s  %s\n", colorRed, colorReset, c.Name())
			fmt.Printf("      current:   %s\n", r.Current)
			fmt.Printf("      would run: %s%s%s\n\n", colorDim, c.FixDescription(), colorReset)
		}
	}

	if fixable == 0 {
		fmt.Printf("  %s✓  Nothing to fix, all checks pass.%s\n\n", colorGreen, colorReset)
		return nil
	}

	fmt.Printf("  %s%s%s\n", colorDim, sep, colorReset)
	fmt.Printf("  %d fixable issue(s). Run fort --fix to apply.\n\n", fixable)
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

const updateCheckInterval = 24 * time.Hour

// checkLatestVersion fetches the latest GitHub release tag at most once per day.
// Returns empty string on any error, if already checked today, or if FORT_NO_UPDATE_CHECK is set.
func checkLatestVersion() <-chan string {
	ch := make(chan string, 1)
	if os.Getenv("FORT_NO_UPDATE_CHECK") != "" {
		ch <- ""
		return ch
	}

	cacheFile := updateCheckCacheFile()
	if info, err := os.Stat(cacheFile); err == nil {
		if time.Since(info.ModTime()) < updateCheckInterval {
			ch <- ""
			return ch
		}
	}

	go func() {
		client := &http.Client{Timeout: 3 * time.Second}
		req, err := http.NewRequest("GET", "https://api.github.com/repos/djadmin/fort/releases/latest", nil)
		if err != nil {
			ch <- ""
			return
		}
		req.Header.Set("Accept", "application/vnd.github+json")
		resp, err := client.Do(req)
		if err != nil || resp.StatusCode != 200 {
			ch <- ""
			return
		}
		defer resp.Body.Close()
		var release struct {
			TagName string `json:"tag_name"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
			ch <- ""
			return
		}
		// Touch the cache file so we don't check again for 24h.
		_ = os.MkdirAll(filepath.Dir(cacheFile), 0700)
		_ = os.WriteFile(cacheFile, nil, 0600)
		ch <- strings.TrimPrefix(release.TagName, "v")
	}()
	return ch
}

func updateCheckCacheFile() string {
	dir, err := os.UserCacheDir()
	if err != nil {
		dir = os.TempDir()
	}
	return filepath.Join(dir, "fort", "update-check")
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
