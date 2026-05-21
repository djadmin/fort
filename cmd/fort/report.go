package main

import (
	"fmt"
	"html/template"
	"os"
	"time"

	"github.com/djadmin/fort/internal/checks"
)

type reportData struct {
	Version   string
	Hostname  string
	Serial    string
	OSVersion string
	Timestamp string
	Summary   jsonSummary
	Policies  []checks.Result
}

func writeReport(results []checks.Result, hostname, serial, osVer, outPath string) error {
	pass, fail, warn := tally(results)
	scoreClass := "pass"
	if fail > 0 {
		scoreClass = "fail"
	} else if warn > 0 {
		scoreClass = "warn"
	}

	data := reportData{
		Version:   version,
		Hostname:  hostname,
		Serial:    serial,
		OSVersion: osVer,
		Timestamp: time.Now().UTC().Format("January 2, 2006 at 15:04 UTC"),
		Summary: jsonSummary{
			Total: len(results),
			Pass:  pass,
			Fail:  fail,
			Warn:  warn,
			Score: fmt.Sprintf("%d/%d", pass, len(results)),
		},
		Policies: results,
	}

	funcMap := template.FuncMap{
		"scoreClass": func() string { return scoreClass },
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
<title>fort — Security Assessment · {{.Hostname}}</title>
<style>
*{box-sizing:border-box;margin:0;padding:0}
body{font-family:-apple-system,BlinkMacSystemFont,"Segoe UI",system-ui,sans-serif;background:#f1f5f9;color:#0f172a;padding:2rem 1rem;font-size:15px;line-height:1.5}
.wrap{max-width:760px;margin:0 auto}
.card{background:white;border-radius:10px;overflow:hidden;box-shadow:0 1px 3px rgba(0,0,0,.07),0 0 0 1px rgba(0,0,0,.04);margin-bottom:1rem}
.header{background:#0f172a;padding:1.5rem 2rem;display:flex;justify-content:space-between;align-items:center}
.brand{color:white;font-size:1.25rem;font-weight:700;letter-spacing:-.025em}
.brand-sub{color:#64748b;font-size:.8125rem;margin-top:.125rem}
.score-box{text-align:right}
.score-num{font-size:2.25rem;font-weight:700;line-height:1}
.score-num.pass{color:#4ade80}
.score-num.fail{color:#f87171}
.score-num.warn{color:#fbbf24}
.score-sub{color:#64748b;font-size:.75rem;margin-top:.25rem}
.meta{display:grid;grid-template-columns:repeat(auto-fill,minmax(160px,1fr));border-bottom:1px solid #f1f5f9}
.meta-item{padding:1rem 2rem;border-right:1px solid #f1f5f9}
.meta-item:last-child{border-right:none}
.meta-item label{display:block;font-size:.6875rem;text-transform:uppercase;letter-spacing:.06em;color:#94a3b8;margin-bottom:.125rem}
.meta-item span{font-size:.9375rem;font-weight:500}
table{width:100%;border-collapse:collapse}
thead tr{border-bottom:2px solid #f1f5f9}
th{text-align:left;font-size:.6875rem;text-transform:uppercase;letter-spacing:.06em;color:#94a3b8;padding:.75rem 2rem;font-weight:500}
td{padding:.9rem 2rem;border-bottom:1px solid #f8fafc;vertical-align:middle}
tr:last-child td{border-bottom:none}
.check-name{font-weight:500}
.badge{display:inline-flex;align-items:center;padding:.2rem .55rem;border-radius:4px;font-size:.6875rem;font-weight:600;text-transform:uppercase;letter-spacing:.04em}
.badge.pass{background:#dcfce7;color:#15803d}
.badge.fail{background:#fee2e2;color:#b91c1c}
.badge.warn{background:#fef9c3;color:#92400e}
.fixed-tag{display:inline-flex;align-items:center;margin-left:.375rem;padding:.2rem .45rem;border-radius:4px;font-size:.6875rem;font-weight:600;background:#dbeafe;color:#1d4ed8}
.val{color:#374151}
.footer{text-align:center;color:#94a3b8;font-size:.75rem;padding:1.5rem}
.footer a{color:#64748b;text-decoration:none}
@media print{
  body{background:white;padding:0;font-size:13px}
  .wrap{max-width:100%}
  .card{box-shadow:none;border:1px solid #e2e8f0;border-radius:6px;margin-bottom:.75rem}
}
</style>
</head>
<body>
<div class="wrap">
  <div class="card">
    <div class="header">
      <div>
        <div class="brand">fort</div>
        <div class="brand-sub">Endpoint Security Assessment</div>
      </div>
      <div class="score-box">
        <div class="score-num {{scoreClass}}">{{.Summary.Score}}</div>
        <div class="score-sub">{{.Summary.Pass}} pass &nbsp;·&nbsp; {{.Summary.Fail}} fail &nbsp;·&nbsp; {{.Summary.Warn}} warn</div>
      </div>
    </div>

    <div class="meta">
      <div class="meta-item"><label>Machine</label><span>{{.Hostname}}</span></div>
      <div class="meta-item"><label>macOS</label><span>{{.OSVersion}}</span></div>
      <div class="meta-item"><label>Serial</label><span>{{.Serial}}</span></div>
      <div class="meta-item"><label>Generated</label><span>{{.Timestamp}}</span></div>
    </div>

    <table>
      <thead>
        <tr><th>Check</th><th>Status</th><th>Found</th><th>Required</th></tr>
      </thead>
      <tbody>
        {{range .Policies}}
        <tr>
          <td class="check-name">{{.Name}}</td>
          <td>
            <span class="badge {{.Status}}">{{.Status}}</span>{{if .Fixed}}<span class="fixed-tag">&#10003; fixed</span>{{end}}
          </td>
          <td class="val">{{.Current}}</td>
          <td class="val">{{.Expected}}</td>
        </tr>
        {{end}}
      </tbody>
    </table>
  </div>

  <div class="footer">
    Generated by <a href="https://github.com/djadmin/fort">fort v{{.Version}}</a> on {{.Timestamp}}<br>
    This report contains security posture information — handle appropriately.
  </div>
</div>
</body>
</html>`
