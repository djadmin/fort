package checks

// Status is the outcome of a single security check.
type Status string

const (
	StatusPass Status = "pass"
	StatusFail Status = "fail"
	StatusWarn Status = "warn"
)

// Result holds the structured output of one check run.
type Result struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Status   Status `json:"status"`
	Current  string `json:"current"`
	Expected string `json:"expected"`
	Fixable  bool   `json:"fixable"`
	Fixed    bool   `json:"fixed"`
	// Evidence is a terminal transcript of the commands run, suitable for auditors.
	Evidence string `json:"evidence,omitempty"`
}

// Check is the interface every security policy must implement.
type Check interface {
	ID() string
	Name() string
	Run() Result
	Fixable() bool
	Fix() error
	// FixDescription returns a human-readable description of what Fix() will do.
	// Returns "" for non-fixable checks.
	FixDescription() string
}
