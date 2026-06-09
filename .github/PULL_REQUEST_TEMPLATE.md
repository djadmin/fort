**What this changes**
A short summary of the change.

**If this adds or removes a check** (see CONTRIBUTING.md)
- [ ] `internal/checks/<name>_darwin.go` implements the `Check` interface
- [ ] Registered in `internal/checks/registry_darwin.go`
- [ ] Framework mappings added in `internal/checks/frameworks.go`
- [ ] Fake result added to `cmd/sample-report/main.go` `randomResults()` (same `Name()` as the real check)
- [ ] Count bumped in `internal/checks/checks_test.go` (`const want`)
- [ ] Listed in the `landing/index.html` "What it checks" grid
- [ ] `make gen` run and `landing/sample-report.html` committed
- [ ] Uses only documented macOS interfaces (no private frameworks)
- [ ] Sets `Evidence` via `evidenceCmd()`, never empty

**Checks**
- [ ] `go test ./...` passes
- [ ] `go vet ./...` is clean
- [ ] `make check-gen` passes (sample report is up to date)
- [ ] Tested on macOS (note version + Apple Silicon/Intel below)

Tested on: macOS ___ , ___ chip

**Notes for reviewer**
Anything worth flagging.
