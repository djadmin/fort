package main

import (
	"fmt"
	"html/template"
	"os"
	"strings"
	"time"

	"github.com/djadmin/fort/internal/checks"
)

type policyRow struct {
	checks.Result
	Frameworks []checks.FrameworkEntry
}

type reportData struct {
	Version   string
	Hostname  string
	Serial    string
	OSVersion string
	Timestamp string
	Summary   jsonSummary
	Policies  []policyRow
}

func writeReport(results []checks.Result, hostname, serial, osVer, outPath string) error {
	pass, fail, warn := tally(results)
	scoreClass := "s-pass"
	if fail > 0 {
		scoreClass = "s-fail"
	} else if warn > 0 {
		scoreClass = "s-warn"
	}

	policies := make([]policyRow, len(results))
	for i, r := range results {
		policies[i] = policyRow{Result: r, Frameworks: checks.FrameworksFor(r.ID)}
	}

	data := reportData{
		Version:   version,
		Hostname:  hostname,
		Serial:    serial,
		OSVersion: osVer,
		Timestamp: time.Now().UTC().Format("2 Jan 2006, 15:04 UTC"),
		Summary:   jsonSummary{Total: len(results), Pass: pass, Fail: fail, Warn: warn, Score: fmt.Sprintf("%d/%d", pass, len(results))},
		Policies:  policies,
	}

	funcMap := template.FuncMap{
		"scoreClass": func() string { return scoreClass },
		"fwClass": func(name string) string {
			r := strings.NewReplacer(" ", "", ".", "")
			return "fw-" + strings.ToLower(r.Replace(name))
		},
		"fwAbbrev": func(name string) string {
			switch name {
			case "SOC 2":
				return "SOC 2"
			case "ISO 27001":
				return "ISO 27001"
			case "NIST CSF":
				return "NIST CSF"
			case "CIS v8":
				return "CIS v8"
			}
			return name
		},
		"joinControls": func(controls []string) string { return strings.Join(controls, " · ") },
		"isPass":       func(s checks.Status) bool { return s == checks.StatusPass },
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
<style>
*,*::before,*::after{box-sizing:border-box;margin:0;padding:0}
body{font-family:-apple-system,BlinkMacSystemFont,"Segoe UI",system-ui,sans-serif;background:#f0f4f8;color:#1a202c;font-size:14px;line-height:1.5;padding:2rem 1rem;-webkit-font-smoothing:antialiased}
.doc{max-width:960px;margin:0 auto}

/* ── Header ──────────────────────── */
.hd{background:#0f172a;border-radius:12px 12px 0 0;padding:2rem 2.5rem;display:flex;justify-content:space-between;align-items:flex-start}
.hd-brand{font-size:1.625rem;font-weight:800;color:#fff;letter-spacing:-.04em;line-height:1}
.hd-sub{font-size:.75rem;color:#64748b;margin-top:.25rem;letter-spacing:.01em}
.score-block{text-align:right}
.score-num{font-size:3rem;font-weight:800;line-height:1;letter-spacing:-.04em}
.s-pass{color:#4ade80}.s-fail{color:#f87171}.s-warn{color:#fbbf24}
.score-desc{font-size:.6875rem;color:#64748b;text-transform:uppercase;letter-spacing:.08em;margin-top:.375rem}

/* ── Meta band ───────────────────── */
.meta{background:#1e293b;display:grid;grid-template-columns:repeat(4,1fr);border-top:1px solid #334155}
.meta-cell{padding:.875rem 1.75rem;border-right:1px solid #334155}
.meta-cell:last-child{border-right:none}
.meta-k{font-size:.5625rem;text-transform:uppercase;letter-spacing:.1em;color:#64748b;margin-bottom:.2rem}
.meta-v{font-size:.875rem;font-weight:600;color:#e2e8f0;white-space:nowrap;overflow:hidden;text-overflow:ellipsis}

/* ── Stat strip ──────────────────── */
.stats{background:#fff;display:grid;grid-template-columns:repeat(3,1fr);border-bottom:2px solid #f0f4f8}
.stat{padding:1.25rem 2.5rem;border-right:1px solid #f0f4f8;display:flex;align-items:baseline;gap:.625rem}
.stat:last-child{border-right:none}
.stat-n{font-size:2.25rem;font-weight:800;line-height:1;letter-spacing:-.03em}
.n-pass{color:#16a34a}.n-fail{color:#dc2626}.n-warn{color:#d97706}
.stat-l{font-size:.8125rem;color:#64748b}

/* ── Table ───────────────────────── */
.tbl-wrap{background:#fff;border-radius:0 0 12px 12px;overflow:hidden}
table{width:100%;border-collapse:collapse}
thead th{padding:.75rem 1.75rem;text-align:left;font-size:.5625rem;text-transform:uppercase;letter-spacing:.1em;color:#94a3b8;font-weight:600;background:#f8fafc;border-bottom:1px solid #e2e8f0}
tbody tr{border-bottom:1px solid #f8fafc;transition:background .1s}
tbody tr:last-child{border-bottom:none}
tbody tr:hover{background:#fafbfc}
td{padding:.875rem 1.75rem;vertical-align:top}

/* ── Check cell ──────────────────── */
.cn{font-weight:600;color:#0f172a;margin-bottom:.4rem;font-size:.9375rem}
.fw-row{display:flex;flex-wrap:wrap;gap:.3rem}
.fw-pill{display:inline-flex;align-items:stretch;border-radius:4px;overflow:hidden;line-height:1}
.fw-name{font-size:.5rem;font-weight:800;text-transform:uppercase;letter-spacing:.06em;padding:.2rem .35rem;display:flex;align-items:center}
.fw-ctrls{font-size:.625rem;font-family:ui-monospace,"SF Mono",Menlo,monospace;padding:.2rem .4rem;display:flex;align-items:center}
.fw-soc2 .fw-name{background:#6d28d9;color:#fff}
.fw-soc2 .fw-ctrls{background:#ede9fe;color:#5b21b6}
.fw-iso27001 .fw-name{background:#1d4ed8;color:#fff}
.fw-iso27001 .fw-ctrls{background:#dbeafe;color:#1e40af}
.fw-nistcsf .fw-name{background:#047857;color:#fff}
.fw-nistcsf .fw-ctrls{background:#d1fae5;color:#065f46}
.fw-cisv8 .fw-name{background:#b45309;color:#fff}
.fw-cisv8 .fw-ctrls{background:#fef3c7;color:#92400e}

/* ── Status pill ─────────────────── */
.status-td{white-space:nowrap}
.pill{display:inline-flex;align-items:center;gap:.3rem;padding:.3rem .75rem;border-radius:999px;font-size:.6875rem;font-weight:700;text-transform:uppercase;letter-spacing:.05em}
.pill::before{content:'';width:5px;height:5px;border-radius:50%;flex-shrink:0}
.pill.pass{background:#dcfce7;color:#15803d}.pill.pass::before{background:#15803d}
.pill.fail{background:#fee2e2;color:#b91c1c}.pill.fail::before{background:#b91c1c}
.pill.warn{background:#fef9c3;color:#92400e}.pill.warn::before{background:#d97706}
.fixed{display:inline-flex;align-items:center;margin-left:.375rem;padding:.2rem .5rem;border-radius:999px;font-size:.5625rem;font-weight:700;background:#dbeafe;color:#1d4ed8}

/* ── Value cells ─────────────────── */
.val{font-size:.9375rem;color:#374151}
.val-na{color:#cbd5e1;font-size:.875rem}

/* ── Footer ──────────────────────── */
.footer{margin-top:1.5rem;display:flex;justify-content:space-between;align-items:center;font-size:.75rem;color:#94a3b8}
.footer a{color:#64748b;text-decoration:none}
.confidential{font-size:.6875rem;font-weight:600;text-transform:uppercase;letter-spacing:.08em;background:#fef2f2;color:#dc2626;padding:.2rem .6rem;border-radius:4px}

@media print{
  body{background:#fff;padding:0;font-size:12px}
  .doc{max-width:100%}
  .hd{border-radius:0}
  .tbl-wrap{border-radius:0}
  tbody tr:hover{background:none}
  .footer{margin-top:1rem}
}
</style>
</head>
<body>
<div class="doc">

  <div class="hd">
    <div>
      <div class="hd-brand">fort</div>
      <div class="hd-sub">Endpoint Security Assessment Report</div>
    </div>
    <div class="score-block">
      <div class="score-num {{scoreClass}}">{{.Summary.Score}}</div>
      <div class="score-desc">compliance score</div>
    </div>
  </div>

  <div class="meta">
    <div class="meta-cell"><div class="meta-k">Machine</div><div class="meta-v">{{.Hostname}}</div></div>
    <div class="meta-cell"><div class="meta-k">macOS</div><div class="meta-v">{{.OSVersion}}</div></div>
    <div class="meta-cell"><div class="meta-k">Serial</div><div class="meta-v">{{.Serial}}</div></div>
    <div class="meta-cell"><div class="meta-k">Generated</div><div class="meta-v">{{.Timestamp}}</div></div>
  </div>

  <div class="stats">
    <div class="stat"><span class="stat-n n-pass">{{.Summary.Pass}}</span><span class="stat-l">checks passing</span></div>
    <div class="stat"><span class="stat-n n-fail">{{.Summary.Fail}}</span><span class="stat-l">checks failing</span></div>
    <div class="stat"><span class="stat-n n-warn">{{.Summary.Warn}}</span><span class="stat-l">warnings</span></div>
  </div>

  <div class="tbl-wrap">
    <table>
      <thead>
        <tr>
          <th style="width:42%">Check &amp; framework controls</th>
          <th style="width:11%">Status</th>
          <th style="width:21%">Found</th>
          <th style="width:26%">Required</th>
        </tr>
      </thead>
      <tbody>
        {{range .Policies}}
        <tr>
          <td>
            <div class="cn">{{.Name}}</div>
            {{if .Frameworks}}
            <div class="fw-row">
              {{range .Frameworks}}{{$cls := fwClass .Name}}
              <span class="fw-pill {{$cls}}">
                <span class="fw-name">{{.Name | fwAbbrev}}</span>
                <span class="fw-ctrls">{{.Controls | joinControls}}</span>
              </span>
              {{end}}
            </div>
            {{end}}
          </td>
          <td class="status-td">
            <span class="pill {{.Status}}">{{.Status}}</span>
            {{if .Fixed}}<span class="fixed">&#10003; fixed</span>{{end}}
          </td>
          <td class="val">{{.Current}}</td>
          <td>
            {{if isPass .Status}}
            <span class="val-na">—</span>
            {{else}}
            <span class="val">{{.Expected}}</span>
            {{end}}
          </td>
        </tr>
        {{end}}
      </tbody>
    </table>
  </div>

  <div class="footer">
    <span>Generated by <a href="https://github.com/djadmin/fort">fort v{{.Version}}</a> &nbsp;·&nbsp; {{.Timestamp}}</span>
    <span class="confidential">Confidential</span>
  </div>

</div>
</body>
</html>`
