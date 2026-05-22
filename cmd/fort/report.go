package main

import (
	"fmt"
	"html/template"
	"math"
	"os"
	"strings"
	"time"

	"github.com/djadmin/fort/internal/checks"
)

type policyRow struct {
	checks.Result
	Frameworks []checks.FrameworkEntry
	UIHint     string // human-readable remediation hint shown in the fail/warn card
}

type reportData struct {
	Version    string
	Hostname   string
	Serial     string
	OSVersion  string
	Timestamp  string
	Summary    jsonSummary
	Policies   []policyRow
	ScorePct   string // "80"
	RingOffset string // "32.67" — SVG stroke-dashoffset for the progress ring
	RingColor  string // "#22c55e" — colour of ring + score text
	ScoreClass string // "s-pass" | "s-warn" | "s-fail" — CSS class for score element
}

// uiHints provides human-readable remediation guidance shown in fail/warn cards.
// These are intentionally friendlier than FixDescription() which shows raw commands.
var uiHints = map[string]string{
	"passwordmgr":  "Install 1Password, Bitwarden, or another dedicated password manager.",
	"filevault":    "Enable in System Settings → Privacy & Security → FileVault. Requires a restart.",
	"screenlock":   "Enable in System Settings → Lock Screen → Require password.",
	"antivirus":    "Install CrowdStrike Falcon, Malwarebytes, SentinelOne, or another AV/EDR agent.",
	"firewall":     "Enable in System Settings → Network → Firewall.",
	"gatekeeper":   "Enable in System Settings → Privacy & Security → App Store and identified developers.",
	"sip":          "Boot into Recovery Mode (hold ⌘R on startup) and run csrutil enable.",
	"ssh":          "Disable in System Settings → Sharing → Remote Login.",
	"localadmin":   "Create a standard user account for daily use. Reserve the admin account for elevated tasks only.",
	"guestaccount": "Disable in System Settings → Users & Groups → Guest User.",
	"autologin":    "Disable in System Settings → Lock Screen → Disable automatic login.",
	"sharing":      "Disable all sharing in System Settings → General → Sharing.",
	"airdrop":      "Set to Contacts Only in Finder → AirDrop or System Settings → General → AirDrop & Handoff.",
	"osupdates":    "Enable in System Settings → General → Software Update → Automatic Updates.",
	"osversion":    "Apply pending updates via System Settings → General → Software Update.",
}

func writeReport(results []checks.Result, hostname, serial, osVer, outPath string) error {
	pass, fail, warn := tally(results)
	total := len(results)

	// ── Score ring ────────────────────────────────────────────────────────────
	// SVG circle r=26 → circumference = 2π×26 = 163.36
	const ringC = 163.36
	var pct int
	var ringOffset float64
	if total > 0 {
		ratio := float64(pass) / float64(total)
		pct = int(math.Round(ratio * 100))
		ringOffset = ringC * (1.0 - ratio)
	}
	var ringColor, scoreClass string
	switch {
	case pct >= 90:
		ringColor, scoreClass = "#22c55e", "s-pass"
	case pct >= 70:
		ringColor, scoreClass = "#eab308", "s-warn"
	case pct >= 50:
		ringColor, scoreClass = "#f97316", "s-warn"
	default:
		ringColor, scoreClass = "#ef4444", "s-fail"
	}

	policies := make([]policyRow, len(results))
	for i, r := range results {
		policies[i] = policyRow{
			Result:     r,
			Frameworks: checks.FrameworksFor(r.ID),
			UIHint:     uiHints[r.ID],
		}
	}

	data := reportData{
		Version:    version,
		Hostname:   hostname,
		Serial:     serial,
		OSVersion:  osVer,
		Timestamp:  time.Now().UTC().Format("2 Jan 2006, 15:04 UTC"),
		Summary:    jsonSummary{Total: total, Pass: pass, Fail: fail, Warn: warn, Score: fmt.Sprintf("%d/%d", pass, total)},
		Policies:   policies,
		ScorePct:   fmt.Sprintf("%d", pct),
		RingOffset: fmt.Sprintf("%.2f", ringOffset),
		RingColor:  ringColor,
		ScoreClass: scoreClass,
	}

	funcMap := template.FuncMap{
		"isPass": func(s checks.Status) bool { return s == checks.StatusPass },
		// frameworkControls returns "CC6.1 · CC6.7" for a given framework name, or "—"
		"frameworkControls": func(fws []checks.FrameworkEntry, name string) string {
			for _, f := range fws {
				if f.Name == name {
					return strings.Join(f.Controls, " · ")
				}
			}
			return "—"
		},
	}

	tmpl, err := template.New("report").Funcs(funcMap).Parse(reportTmpl)
	if err != nil {
		return fmt.Errorf("parse template: %w", err)
	}
	f, err := os.Create(outPath)
	if err != nil {
		return fmt.Errorf("create file: %w", err)
	}
	defer f.Close()
	return tmpl.Execute(f, data)
}

const reportTmpl = `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width,initial-scale=1">
<title>Security Assessment · {{.Hostname}}</title>
<link href="https://fonts.googleapis.com/css2?family=DM+Sans:opsz,wght@9..40,300;9..40,400;9..40,500;9..40,600;9..40,700;9..40,800&family=DM+Mono:wght@400;500&display=swap" rel="stylesheet">
<style>
*,*::before,*::after{box-sizing:border-box;margin:0;padding:0}
:root{
  --bg:#f4f4f3;--surface:#fff;--ink:#18181b;--ink-2:#3f3f46;--ink-3:#71717a;--ink-4:#a1a1aa;
  --border:#e4e4e7;--border-light:#f4f4f5;
  --green:#15803d;--green-soft:#22c55e;--green-bg:#f0fdf4;--green-border:#bbf7d0;
  --red:#dc2626;--red-bg:#fef2f2;--red-border:#fecaca;
  --amber:#d97706;--amber-bg:#fffbeb;--amber-border:#fde68a;
  --radius:10px;
  --font:'DM Sans',-apple-system,system-ui,sans-serif;
  --mono:'DM Mono',ui-monospace,'SF Mono',monospace;
}
body{font-family:var(--font);background:var(--bg);color:var(--ink);font-size:14px;line-height:1.55;padding:2.5rem 1rem;-webkit-font-smoothing:antialiased}
.doc{max-width:880px;margin:0 auto}

/* ── Header ── */
.hd{background:linear-gradient(135deg,#0c0f1a 0%,#151b2e 100%);border-radius:var(--radius) var(--radius) 0 0;padding:2rem 2.5rem;display:flex;justify-content:space-between;align-items:center}
.brand{font-size:1.375rem;font-weight:700;color:#fff;letter-spacing:-.03em;line-height:1}
.brand-tag{font-size:.6875rem;color:#4b5574;margin-top:.3rem;letter-spacing:.01em;font-weight:400}
.score-cluster{display:flex;align-items:center;gap:1.25rem}
.ring-wrap{position:relative;width:56px;height:56px;flex-shrink:0}
.ring-wrap svg{transform:rotate(-90deg);display:block}
.ring-track{fill:none;stroke:#1e2540;stroke-width:4.5}
.ring-fill{fill:none;stroke-width:4.5;stroke-linecap:round}
.ring-text{position:absolute;inset:0;display:flex;align-items:center;justify-content:center;font-size:.8125rem;font-weight:700;color:#fff;letter-spacing:-.01em}
.score-right{text-align:right}
.score-pct{font-size:2rem;font-weight:800;line-height:1;letter-spacing:-.04em}
.score-label{font-size:.5625rem;color:#4b5574;text-transform:uppercase;letter-spacing:.1em;margin-top:.2rem}

/* ── Meta band ── */
.meta{background:#141829;display:grid;grid-template-columns:repeat(4,1fr);border-top:1px solid #1e2540}
.meta-cell{padding:.625rem 1.5rem;border-right:1px solid #1e2540}
.meta-cell:last-child{border-right:none}
.meta-k{font-size:.5rem;text-transform:uppercase;letter-spacing:.1em;color:#4b5574;margin-bottom:.1rem}
.meta-v{font-size:.8125rem;font-weight:500;color:#c8cee0;white-space:nowrap;overflow:hidden;text-overflow:ellipsis}

/* ── Stats ── */
.stats{background:var(--surface);display:grid;grid-template-columns:repeat(3,1fr);border-bottom:1px solid var(--border)}
.stat{padding:1rem 1.5rem;border-right:1px solid var(--border-light);display:flex;align-items:baseline;gap:.5rem}
.stat:last-child{border-right:none}
.stat-n{font-size:1.5rem;font-weight:700;line-height:1;letter-spacing:-.02em}
.n-pass{color:var(--green)}.n-fail{color:var(--red)}.n-warn{color:var(--amber)}
.stat-l{font-size:.8125rem;color:var(--ink-3)}

/* ── Table ── */
.tbl-wrap{background:var(--surface);overflow:hidden}
table{width:100%;border-collapse:collapse}
thead th{padding:.625rem 1.5rem;text-align:left;font-size:.5625rem;text-transform:uppercase;letter-spacing:.1em;color:var(--ink-4);font-weight:600;background:#fafafa;border-bottom:1px solid var(--border)}
tbody tr{border-bottom:1px solid var(--border-light);transition:background .1s}
tbody tr:last-child{border-bottom:none}
tbody tr:hover{background:#fafbfc}
td{padding:.75rem 1.5rem;vertical-align:top}
.check-name{font-weight:600;color:var(--ink);font-size:.875rem;line-height:1.3}
.val{font-size:.875rem;color:var(--ink-2)}

/* ── Status pill ── */
.pill{display:inline-flex;align-items:center;gap:.3rem;padding:.225rem .6rem;border-radius:999px;font-size:.6875rem;font-weight:600;text-transform:uppercase;letter-spacing:.04em}
.pill::before{content:'';width:5px;height:5px;border-radius:50%;flex-shrink:0}
.pill.pass{background:var(--green-bg);color:var(--green)}.pill.pass::before{background:var(--green)}
.pill.fail{background:var(--red-bg);color:var(--red)}.pill.fail::before{background:var(--red)}
.pill.warn{background:var(--amber-bg);color:var(--amber)}.pill.warn::before{background:var(--amber)}
.fixed{display:inline-flex;align-items:center;margin-left:.375rem;padding:.2rem .5rem;border-radius:999px;font-size:.5625rem;font-weight:700;background:#dbeafe;color:#1d4ed8;vertical-align:middle}

/* ── Fail / warn card ── */
.remedy-card{margin-top:.5rem;padding:.5rem .75rem;border-radius:6px;font-size:.75rem;line-height:1.55}
.remedy-card.fail{background:var(--red-bg);border:1px solid var(--red-border)}
.remedy-card.warn{background:var(--amber-bg);border:1px solid var(--amber-border)}
.remedy-found{font-weight:500}
.remedy-card.fail .remedy-found{color:var(--red)}
.remedy-card.warn .remedy-found{color:var(--amber)}
.remedy-found strong{font-weight:700}
.remedy-hint{color:var(--ink-2);margin-top:.25rem}
.remedy-hint code{font-family:var(--mono);font-size:.6875rem;background:#fff;padding:.1rem .3rem;border-radius:3px;border:1px solid var(--border)}

/* ── Compliance section ── */
.compliance-section{background:var(--surface);border-radius:0 0 var(--radius) var(--radius);border-top:1px solid var(--border)}
.compliance-toggle{display:flex;align-items:center;justify-content:space-between;padding:.75rem 1.5rem;cursor:pointer;user-select:none;transition:background .15s}
.compliance-toggle:hover{background:var(--border-light)}
.ct-left{display:flex;align-items:center;gap:.5rem}
.ct-label{font-size:.6875rem;font-weight:600;color:var(--ink-3);text-transform:uppercase;letter-spacing:.06em}
.ct-frameworks{font-size:.625rem;color:var(--ink-4);background:var(--border-light);padding:.125rem .5rem;border-radius:999px;border:1px solid var(--border)}
.ct-chevron{width:14px;height:14px;color:var(--ink-4);transition:transform .2s ease}
.ct-chevron.open{transform:rotate(180deg)}
.compliance-body{display:none;padding:0 1.5rem 1rem}
.compliance-body.open{display:block}
.cm-table{width:100%;border-collapse:collapse;font-size:.75rem}
.cm-table thead th{padding:.4rem .625rem;font-size:.5rem;text-transform:uppercase;letter-spacing:.08em;color:var(--ink-4);font-weight:600;background:var(--border-light);border-bottom:1px solid var(--border);text-align:left}
.cm-table td{padding:.35rem .625rem;border-bottom:1px solid var(--border-light);vertical-align:middle}
.cm-table tbody tr:last-child td{border-bottom:none}
.cm-name{font-weight:500;color:var(--ink-2)}
.cm-ctrl{font-family:var(--mono);font-size:.625rem;color:var(--ink-3)}
.cm-dot{display:inline-block;width:6px;height:6px;border-radius:50%;vertical-align:middle}
.cm-dot.pass{background:var(--green)}.cm-dot.fail{background:var(--red)}.cm-dot.warn{background:var(--amber)}

/* ── Footer ── */
.footer{margin-top:1.25rem;display:flex;justify-content:space-between;align-items:center;font-size:.6875rem;color:var(--ink-4)}
.footer a{color:var(--ink-3);text-decoration:none}
.footer a:hover{text-decoration:underline}
.confidential{font-size:.5625rem;font-weight:700;text-transform:uppercase;letter-spacing:.08em;background:var(--red-bg);color:var(--red);padding:.175rem .5rem;border-radius:3px;border:1px solid var(--red-border)}

/* ── Responsive ── */
@media(max-width:680px){
  .meta{grid-template-columns:repeat(2,1fr)}
  .hd{flex-direction:column;gap:1.25rem;align-items:flex-start}
  .score-cluster{align-self:flex-end}
}
@media print{
  body{background:#fff;padding:0;font-size:12px}
  .doc{max-width:100%}
  .hd,.meta{print-color-adjust:exact;-webkit-print-color-adjust:exact}
  .hd{border-radius:0}
  .compliance-body{display:block!important}
  .compliance-toggle{display:none}
  tbody tr:hover{background:none}
}
</style>
</head>
<body>
<div class="doc">

  <div class="hd">
    <div>
      <div class="brand">fort</div>
      <div class="brand-tag">Endpoint Security Assessment</div>
    </div>
    <div class="score-cluster">
      <div class="score-right">
        <div class="score-pct {{.ScoreClass}}" style="color:{{.RingColor}}">{{.ScorePct}}<span style="font-size:.5em;font-weight:600;opacity:.5">%</span></div>
        <div class="score-label">Compliance</div>
      </div>
      <div class="ring-wrap">
        <svg viewBox="0 0 64 64" width="56" height="56">
          <circle class="ring-track" cx="32" cy="32" r="26"/>
          <circle class="ring-fill" cx="32" cy="32" r="26"
            stroke="{{.RingColor}}"
            stroke-dasharray="163.36"
            stroke-dashoffset="{{.RingOffset}}"/>
        </svg>
        <div class="ring-text">{{.Summary.Pass}}/{{.Summary.Total}}</div>
      </div>
    </div>
  </div>

  <div class="meta">
    <div class="meta-cell"><div class="meta-k">Machine</div><div class="meta-v">{{.Hostname}}</div></div>
    <div class="meta-cell"><div class="meta-k">macOS</div><div class="meta-v">{{.OSVersion}}</div></div>
    <div class="meta-cell"><div class="meta-k">Serial</div><div class="meta-v">{{.Serial}}</div></div>
    <div class="meta-cell"><div class="meta-k">Generated</div><div class="meta-v">{{.Timestamp}}</div></div>
  </div>

  <div class="stats">
    <div class="stat"><span class="stat-n n-pass">{{.Summary.Pass}}</span><span class="stat-l">passing</span></div>
    <div class="stat"><span class="stat-n n-fail">{{.Summary.Fail}}</span><span class="stat-l">failing</span></div>
    <div class="stat"><span class="stat-n n-warn">{{.Summary.Warn}}</span><span class="stat-l">warnings</span></div>
  </div>

  <div class="tbl-wrap">
    <table>
      <thead><tr>
        <th style="width:40%">Check</th>
        <th style="width:15%">Status</th>
        <th>Value</th>
      </tr></thead>
      <tbody>
        {{range .Policies}}
        <tr>
          <td><div class="check-name">{{.Name}}</div></td>
          <td>
            <span class="pill {{.Status}}">{{.Status}}</span>
            {{if .Fixed}}<span class="fixed">&#10003; fixed</span>{{end}}
            {{if not (isPass .Status)}}
            <div class="remedy-card {{.Status}}">
              <div class="remedy-found">Found <strong>{{.Current}}</strong> — expected <strong>{{.Expected}}</strong></div>
              {{if .UIHint}}
              <div class="remedy-hint">{{.UIHint}}{{if .Fixable}} &nbsp;<code>sudo fort --fix</code>{{end}}</div>
              {{end}}
            </div>
            {{end}}
          </td>
          <td class="val">{{.Current}}</td>
        </tr>
        {{end}}
      </tbody>
    </table>
  </div>

  <div class="compliance-section">
    <div class="compliance-toggle" onclick="toggleCompliance()">
      <div class="ct-left">
        <span class="ct-label">Compliance Framework Mapping</span>
        <span class="ct-frameworks">SOC 2 · ISO 27001 · NIST CSF · CIS v8</span>
      </div>
      <svg class="ct-chevron" id="chevron" viewBox="0 0 16 16" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round"><path d="M4 6l4 4 4-4"/></svg>
    </div>
    <div class="compliance-body" id="complianceBody">
      <table class="cm-table">
        <thead><tr>
          <th style="width:4%"></th>
          <th style="width:22%">Check</th>
          <th style="width:18%">SOC 2</th>
          <th style="width:20%">ISO 27001</th>
          <th style="width:20%">NIST CSF</th>
          <th style="width:16%">CIS v8</th>
        </tr></thead>
        <tbody>
          {{range .Policies}}
          <tr>
            <td><span class="cm-dot {{.Status}}"></span></td>
            <td class="cm-name">{{.Name}}</td>
            <td class="cm-ctrl">{{frameworkControls .Frameworks "SOC 2"}}</td>
            <td class="cm-ctrl">{{frameworkControls .Frameworks "ISO 27001"}}</td>
            <td class="cm-ctrl">{{frameworkControls .Frameworks "NIST CSF"}}</td>
            <td class="cm-ctrl">{{frameworkControls .Frameworks "CIS v8"}}</td>
          </tr>
          {{end}}
        </tbody>
      </table>
    </div>
  </div>

  <div class="footer">
    <span>Generated by <a href="https://github.com/djadmin/fort">fort v{{.Version}}</a> &nbsp;·&nbsp; {{.Timestamp}}</span>
    <span class="confidential">Confidential</span>
  </div>

</div>
<script>
function toggleCompliance() {
  document.getElementById('complianceBody').classList.toggle('open');
  document.getElementById('chevron').classList.toggle('open');
}
</script>
</body>
</html>`
