**What this changes**
A short summary of the change.

**If this adds a new check**
- [ ] `internal/checks/<name>_darwin.go` implements the `Check` interface
- [ ] Registered in `internal/checks/registry_darwin.go`
- [ ] Framework mappings added in `internal/checks/frameworks.go`
- [ ] Fake result added to `cmd/sample-report/main.go` `randomResults()`
- [ ] Uses only documented macOS interfaces (no private frameworks)
- [ ] Sets `Evidence` via `evidenceCmd()`, never empty

**Checks**
- [ ] `go test ./...` passes
- [ ] `go vet ./...` is clean
- [ ] Tested on macOS (note version + Apple Silicon/Intel below)

Tested on: macOS ___ , ___ chip

**Notes for reviewer**
Anything worth flagging.
