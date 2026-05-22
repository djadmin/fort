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
}

type reportData struct {
	Version    string
	Hostname   string
	Serial     string
	OSVersion  string
	Timestamp  string
	Summary    jsonSummary
	Policies   []policyRow
	ScorePct   string
	RingOffset string // SVG stroke-dashoffset; circumference = 2π×26 = 163.36
	RingColor  string
	ScoreClass string // s-pass | s-warn | s-fail
}

func writeReport(results []checks.Result, hostname, serial, osVer, outPath string) error {
	pass, fail, warn := tally(results)
	total := len(results)

	const ringC = 131.95 // 2π × r(21) — matches SVG viewBox r="21"
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
		policies[i] = policyRow{Result: r, Frameworks: checks.FrameworksFor(r.ID)}
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
		"fwControls": func(fws []checks.FrameworkEntry, name string) string {
			for _, f := range fws {
				if f.Name == name {
					return strings.Join(f.Controls, " · ")
				}
			}
			return "—"
		},
		"soc2": func(fws []checks.FrameworkEntry) string {
			for _, f := range fws {
				if f.Name == "SOC 2" {
					return strings.Join(f.Controls, " · ")
				}
			}
			return ""
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
<link href="https://fonts.googleapis.com/css2?family=DM+Sans:opsz,wght@9..40,400;9..40,500;9..40,600;9..40,700;9..40,800&family=DM+Mono:wght@400;500&display=swap" rel="stylesheet">
<style>
*,*::before,*::after{box-sizing:border-box;margin:0;padding:0}
:root{
  --page:#F5F5F4;
  --surface:#FFFFFF;
  --ink:#1C1917;
  --ink-2:#44403C;
  --ink-3:#78716C;
  --ink-4:#A8A29E;
  --border:#E7E5E4;
  --border-2:#F5F5F4;
  --hd-bg:#0C0F1C;
  --hd-bg-2:#131829;
  --meta-bg:#101422;
  --meta-border:#1C2235;
  --pass:#166534;--pass-bg:#F0FDF4;--pass-bd:#BBF7D0;
  --fail:#991B1B;--fail-bg:#FEF2F2;--fail-bd:#FECACA;
  --warn:#92400E;--warn-bg:#FFFBEB;--warn-bd:#FDE68A;
  --radius:12px;
  --font:'DM Sans',-apple-system,BlinkMacSystemFont,'Segoe UI',system-ui,sans-serif;
  --mono:'DM Mono',ui-monospace,'SF Mono',Menlo,monospace;
}
body{font-family:var(--font);background:var(--page);color:var(--ink);font-size:13px;line-height:1.6;padding:2rem 1rem;-webkit-font-smoothing:antialiased}
.doc{max-width:860px;margin:0 auto}
a{color:inherit;text-decoration:none}

/* ── Header ─────────────────────────────────── */
.hd{
  background:linear-gradient(140deg,var(--hd-bg) 0%,var(--hd-bg-2) 100%);
  border-radius:var(--radius) var(--radius) 0 0;
  padding:1.75rem 2rem;
  display:flex;justify-content:space-between;align-items:center;
  gap:1.5rem;
}
.brand{font-size:1.25rem;font-weight:800;color:#FFF;letter-spacing:-.035em;line-height:1}
.brand-sub{font-size:.6875rem;color:#4B5574;margin-top:.2rem;font-weight:400;letter-spacing:.01em}

.score-group{display:flex;align-items:center;gap:1rem;flex-shrink:0}
.ring-wrap{position:relative;width:52px;height:52px}
.ring-wrap svg{transform:rotate(-90deg);display:block}
.ring-track{fill:none;stroke:#1C2444;stroke-width:4}
.ring-arc{fill:none;stroke-width:4;stroke-linecap:round}
.ring-label{
  position:absolute;inset:0;
  display:flex;flex-direction:column;align-items:center;justify-content:center;
  font-family:var(--font);
}
.ring-frac{font-size:.625rem;font-weight:700;color:#FFF;line-height:1}
.ring-unit{font-size:.4375rem;color:#4B5574;margin-top:.1rem;text-transform:uppercase;letter-spacing:.05em}
.score-text{}
.score-pct{font-size:2.125rem;font-weight:800;letter-spacing:-.04em;line-height:1}
.score-tag{font-size:.5625rem;text-transform:uppercase;letter-spacing:.1em;color:#4B5574;margin-top:.2rem}

/* ── Meta strip ─────────────────────────────── */
.meta{
  background:var(--meta-bg);
  display:grid;grid-template-columns:repeat(4,1fr);
  border-top:1px solid var(--meta-border);
}
.mc{padding:.625rem 1.625rem;border-right:1px solid var(--meta-border)}
.mc:last-child{border-right:none}
.mk{font-size:.5rem;text-transform:uppercase;letter-spacing:.1em;color:#394361;margin-bottom:.15rem}
.mv{font-size:.8125rem;font-weight:500;color:#B8BFCD;white-space:nowrap;overflow:hidden;text-overflow:ellipsis}

/* ── Stats bar ──────────────────────────────── */
.stats{
  background:var(--surface);
  display:flex;
  border-bottom:1px solid var(--border);
}
.sv{
  flex:1;padding:.9rem 1.625rem;
  border-right:1px solid var(--border-2);
  display:flex;align-items:baseline;gap:.5rem;
}
.sv:last-child{border-right:none}
.sn{font-size:1.5rem;font-weight:800;letter-spacing:-.03em;line-height:1}
.sn.p{color:var(--pass)}.sn.f{color:var(--fail)}.sn.w{color:var(--warn)}
.sl{font-size:.8125rem;color:var(--ink-3)}

/* ── Checks table ───────────────────────────── */
.card{background:var(--surface)}
.section-head{
  padding:.5rem 1.625rem;
  font-size:.5rem;font-weight:700;text-transform:uppercase;letter-spacing:.12em;
  color:var(--ink-4);background:#FAFAF9;
  border-bottom:1px solid var(--border);
  border-top:1px solid var(--border);
}
table{width:100%;border-collapse:collapse}
thead th{
  padding:.5rem 1.625rem;
  text-align:left;
  font-size:.5rem;font-weight:700;text-transform:uppercase;letter-spacing:.1em;
  color:var(--ink-4);background:#FAFAF9;
  border-bottom:1px solid var(--border);
}
tbody tr{border-bottom:1px solid var(--border-2)}
tbody tr:last-child{border-bottom:none}
td{padding:.7rem 1.625rem;vertical-align:middle}

/* row status accents */
.accent-td{width:3px;padding:0!important}
tr.row-pass .accent-td{background:transparent}
tr.row-fail .accent-td{background:var(--fail-bd)}
tr.row-warn .accent-td{background:var(--warn-bd)}
tr.row-fail{background:#FFFCFC}
tr.row-warn{background:#FFFDF7}

.cn{font-size:.875rem;font-weight:600;color:var(--ink);line-height:1.3}
.soc2-hint{font-size:.625rem;color:var(--ink-4);font-family:var(--mono);margin-top:.2rem;letter-spacing:.01em}

/* Status pill */
.pill{
  display:inline-flex;align-items:center;gap:.3rem;
  padding:.2rem .6rem;border-radius:999px;
  font-size:.625rem;font-weight:700;text-transform:uppercase;letter-spacing:.05em;
  white-space:nowrap;
}
.pill-icon{font-size:.6rem;line-height:1}
.pill.pass{background:var(--pass-bg);color:var(--pass);border:1px solid var(--pass-bd)}
.pill.fail{background:var(--fail-bg);color:var(--fail);border:1px solid var(--fail-bd)}
.pill.warn{background:var(--warn-bg);color:var(--warn);border:1px solid var(--warn-bd)}

.fixed-chip{
  display:inline-flex;align-items:center;gap:.2rem;
  margin-left:.375rem;padding:.175rem .4rem;border-radius:999px;
  font-size:.5625rem;font-weight:700;
  background:#EFF6FF;color:#1D4ED8;border:1px solid #BFDBFE;
}

.val{font-size:.875rem;color:var(--ink-2)}
.val-muted{color:var(--ink-4)}
.arrow{color:var(--ink-4);margin:0 .3rem;font-size:.8125rem}
.req{font-size:.875rem;color:var(--fail);font-weight:500}
.req-warn{color:var(--warn)}

/* ── Framework reference ────────────────────── */
.fw-ref{}
.fw-table{font-size:.6875rem}
.fw-table thead th{font-size:.5rem;padding:.4rem 1rem}
.fw-table td{padding:.35rem 1rem;border-bottom:1px solid var(--border-2)}
.fw-table tbody tr:last-child td{border-bottom:none}
.fw-dot{width:3px;padding:0!important}
tr.fw-pass .fw-dot{background:transparent}
tr.fw-fail .fw-dot{background:var(--fail-bd)}
tr.fw-warn .fw-dot{background:var(--warn-bd)}
.fw-name{font-weight:500;color:var(--ink-2);font-size:.75rem}
.fw-ctrl{font-family:var(--mono);font-size:.625rem;color:var(--ink-3)}

/* ── Footer ─────────────────────────────────── */
.footer{
  margin-top:1.25rem;padding:0 .25rem;
  display:flex;justify-content:space-between;align-items:center;
  font-size:.6875rem;color:var(--ink-4);
}
.footer a:hover{color:var(--ink-3)}
.confidential{
  font-size:.5rem;font-weight:700;text-transform:uppercase;letter-spacing:.1em;
  background:var(--fail-bg);color:var(--fail);
  padding:.2rem .55rem;border-radius:3px;border:1px solid var(--fail-bd);
}

/* ── Responsive / Print ─────────────────────── */
@media(max-width:640px){
  .meta{grid-template-columns:repeat(2,1fr)}
  .hd{flex-direction:column;gap:1rem;align-items:flex-start}
  .score-group{align-self:flex-end}
}
@media print{
  body{background:#fff;padding:0;font-size:11px}
  .doc{max-width:100%}
  .hd,.meta{print-color-adjust:exact;-webkit-print-color-adjust:exact}
  .hd{border-radius:0}
  tr.row-fail,tr.row-warn{print-color-adjust:exact;-webkit-print-color-adjust:exact}
  .accent-td,.fw-dot{print-color-adjust:exact;-webkit-print-color-adjust:exact}
}
</style>
</head>
<body>
<div class="doc">

<!-- ── Header ──────────────────────────────────── -->
<div class="hd">
  <div>
    <div class="brand">fort</div>
    <div class="brand-sub">Endpoint Security Assessment Report</div>
  </div>
  <div class="score-group">
    <div class="score-text">
      <div class="score-pct {{.ScoreClass}}" style="color:{{.RingColor}}">{{.ScorePct}}<span style="font-size:.45em;font-weight:600;opacity:.45">%</span></div>
      <div class="score-tag">Compliance</div>
    </div>
    <div class="ring-wrap">
      <svg viewBox="0 0 52 52" width="52" height="52">
        <circle class="ring-track" cx="26" cy="26" r="21"/>
        <circle class="ring-arc" cx="26" cy="26" r="21"
          stroke="{{.RingColor}}"
          stroke-dasharray="131.95"
          stroke-dashoffset="{{.RingOffset}}"/>
      </svg>
      <div class="ring-label">
        <span class="ring-frac">{{.Summary.Pass}}/{{.Summary.Total}}</span>
        <span class="ring-unit">checks</span>
      </div>
    </div>
  </div>
</div>

<!-- ── Meta ────────────────────────────────────── -->
<div class="meta">
  <div class="mc"><div class="mk">Machine</div><div class="mv">{{.Hostname}}</div></div>
  <div class="mc"><div class="mk">macOS</div><div class="mv">{{.OSVersion}}</div></div>
  <div class="mc"><div class="mk">Serial</div><div class="mv">{{.Serial}}</div></div>
  <div class="mc"><div class="mk">Generated</div><div class="mv">{{.Timestamp}}</div></div>
</div>

<!-- ── Stats ────────────────────────────────────── -->
<div class="stats card">
  <div class="sv"><span class="sn p">{{.Summary.Pass}}</span><span class="sl">passing</span></div>
  <div class="sv"><span class="sn f">{{.Summary.Fail}}</span><span class="sl">failing</span></div>
  <div class="sv"><span class="sn w">{{.Summary.Warn}}</span><span class="sl">warnings</span></div>
</div>

<!-- ── Check results ────────────────────────────── -->
<div class="card">
  <table>
    <thead>
      <tr>
        <th style="width:3px;padding:0"></th>
        <th style="width:42%">Check</th>
        <th style="width:12%">Status</th>
        <th>Finding</th>
      </tr>
    </thead>
    <tbody>
    {{range .Policies}}
    <tr class="row-{{.Status}}">
      <td class="accent-td"></td>
      <td>
        <div class="cn">{{.Name}}</div>
        {{if soc2 .Frameworks}}<div class="soc2-hint">SOC 2 {{soc2 .Frameworks}}</div>{{end}}
      </td>
      <td>
        <span class="pill {{.Status}}">
          {{if isPass .Status}}<span class="pill-icon">✓</span> Pass
          {{else if eq .Status "fail"}}<span class="pill-icon">✗</span> Fail
          {{else}}<span class="pill-icon">◐</span> Warn
          {{end}}
        </span>
        {{if .Fixed}}<span class="fixed-chip">✓ fixed</span>{{end}}
      </td>
      <td>
        {{if isPass .Status}}
          <span class="val">{{.Current}}</span>
        {{else}}
          <span class="val">{{.Current}}</span><span class="arrow">→</span><span class="{{if eq .Status "fail"}}req{{else}}req-warn{{end}}">{{.Expected}}</span>
        {{end}}
      </td>
    </tr>
    {{end}}
    </tbody>
  </table>
</div>

<!-- ── Framework reference ──────────────────────── -->
<div class="card fw-ref" style="margin-top:1rem;border-radius:var(--radius)">
  <div class="section-head">Framework Control Reference &nbsp;—&nbsp; SOC 2 · ISO 27001 · NIST CSF · CIS v8</div>
  <table class="fw-table">
    <thead>
      <tr>
        <th style="width:3px;padding:0"></th>
        <th style="width:26%">Check</th>
        <th style="width:16%">SOC 2</th>
        <th style="width:22%">ISO 27001</th>
        <th style="width:22%">NIST CSF</th>
        <th style="width:14%">CIS v8</th>
      </tr>
    </thead>
    <tbody>
    {{range .Policies}}
    <tr class="fw-{{.Status}}">
      <td class="fw-dot"></td>
      <td class="fw-name">{{.Name}}</td>
      <td class="fw-ctrl">{{fwControls .Frameworks "SOC 2"}}</td>
      <td class="fw-ctrl">{{fwControls .Frameworks "ISO 27001"}}</td>
      <td class="fw-ctrl">{{fwControls .Frameworks "NIST CSF"}}</td>
      <td class="fw-ctrl">{{fwControls .Frameworks "CIS v8"}}</td>
    </tr>
    {{end}}
    </tbody>
  </table>
</div>

<!-- ── Footer ───────────────────────────────────── -->
<div class="footer">
  <span>Generated by <a href="https://github.com/djadmin/fort">fort v{{.Version}}</a> &nbsp;·&nbsp; {{.Timestamp}}</span>
  <span class="confidential">Confidential</span>
</div>

</div>
</body>
</html>`
