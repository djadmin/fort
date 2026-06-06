package checks

// FrameworkEntry maps one compliance framework to the controls a check satisfies.
type FrameworkEntry struct {
	Name     string   // e.g. "SOC 2", "ISO 27001"
	Controls []string // e.g. ["CC6.1", "CC6.7"]
}

// FrameworksFor returns ordered framework→control mappings for a given check ID.
// Returns nil if no mappings are defined.
func FrameworksFor(id string) []FrameworkEntry {
	return frameworkMappings[id]
}

// frameworkMappings defines every check's coverage across four major frameworks.
// Add new frameworks here — no other code needs to change.
//
// Frameworks covered:
//   SOC 2    — AICPA Trust Services Criteria (CC series)
//   ISO 27001 — ISO/IEC 27001:2022 Annex A controls
//   NIST CSF  — NIST Cybersecurity Framework 2.0
//   CIS v8    — CIS Critical Security Controls v8
var frameworkMappings = map[string][]FrameworkEntry{
	"passwordmgr": {
		{"SOC 2", []string{"CC6.1"}},
		{"ISO 27001", []string{"A.5.17", "A.8.5"}},
		{"NIST CSF", []string{"PR.AC-1", "PR.AC-7"}},
		{"CIS v8", []string{"5.2", "5.4"}},
	},
	"filevault": {
		{"SOC 2", []string{"CC6.1", "CC6.7"}},
		{"ISO 27001", []string{"A.8.3", "A.8.24"}},
		{"NIST CSF", []string{"PR.DS-1", "PR.DS-5"}},
		{"CIS v8", []string{"3.6", "3.11"}},
	},
	"screenlock": {
		{"SOC 2", []string{"CC6.1"}},
		{"ISO 27001", []string{"A.8.1", "A.8.5"}},
		{"NIST CSF", []string{"PR.AC-3", "PR.AC-7"}},
		{"CIS v8", []string{"4.3"}},
	},
	"antivirus": {
		{"SOC 2", []string{"CC6.8", "CC7.1"}},
		{"ISO 27001", []string{"A.8.7"}},
		{"NIST CSF", []string{"DE.CM-4", "PR.PT-3"}},
		{"CIS v8", []string{"10.1", "10.2"}},
	},
	"firewall": {
		{"SOC 2", []string{"CC6.6"}},
		{"ISO 27001", []string{"A.8.20", "A.8.22"}},
		{"NIST CSF", []string{"PR.AC-5", "DE.CM-1"}},
		{"CIS v8", []string{"4.5", "12.1"}},
	},
	"gatekeeper": {
		{"SOC 2", []string{"CC6.8"}},
		{"ISO 27001", []string{"A.8.7", "A.8.19"}},
		{"NIST CSF", []string{"PR.PT-3", "DE.CM-4"}},
		{"CIS v8", []string{"2.6", "10.5"}},
	},
	"sip": {
		{"SOC 2", []string{"CC6.8"}},
		{"ISO 27001", []string{"A.8.7", "A.8.9"}},
		{"NIST CSF", []string{"PR.PT-3", "PR.IP-1"}},
		{"CIS v8", []string{"2.5", "10.5"}},
	},
	"ssh": {
		{"SOC 2", []string{"CC6.6"}},
		{"ISO 27001", []string{"A.8.20", "A.8.21"}},
		{"NIST CSF", []string{"PR.AC-5", "PR.AC-3"}},
		{"CIS v8", []string{"4.8", "12.1"}},
	},
	"localadmin": {
		{"SOC 2", []string{"CC6.3"}},
		{"ISO 27001", []string{"A.5.15", "A.8.2"}},
		{"NIST CSF", []string{"PR.AC-4"}},
		{"CIS v8", []string{"5.4", "5.5"}},
	},
	"guestaccount": {
		{"SOC 2", []string{"CC6.1"}},
		{"ISO 27001", []string{"A.5.15", "A.8.5"}},
		{"NIST CSF", []string{"PR.AC-1"}},
		{"CIS v8", []string{"5.3"}},
	},
	"autologin": {
		{"SOC 2", []string{"CC6.1"}},
		{"ISO 27001", []string{"A.8.5"}},
		{"NIST CSF", []string{"PR.AC-1", "PR.AC-7"}},
		{"CIS v8", []string{"4.3"}},
	},
	"sudo_touchid": {
		{"SOC 2", []string{"CC6.1", "CC6.3"}},
		{"ISO 27001", []string{"A.5.17", "A.8.5"}},
		{"NIST CSF", []string{"PR.AC-7"}},
		{"CIS v8", []string{"6.3", "6.5"}},
	},
	"sharing": {
		{"SOC 2", []string{"CC6.6"}},
		{"ISO 27001", []string{"A.8.20", "A.8.21"}},
		{"NIST CSF", []string{"PR.AC-5", "PR.PT-3"}},
		{"CIS v8", []string{"4.8", "12.1"}},
	},
	"airdrop": {
		{"SOC 2", []string{"CC6.6"}},
		{"ISO 27001", []string{"A.8.20"}},
		{"NIST CSF", []string{"PR.DS-5"}},
		{"CIS v8", []string{"3.14", "12.1"}},
	},
	"osupdates": {
		{"SOC 2", []string{"CC6.8"}},
		{"ISO 27001", []string{"A.8.8"}},
		{"NIST CSF", []string{"PR.IP-12", "ID.RA-1"}},
		{"CIS v8", []string{"7.1", "7.3"}},
	},
	"osversion": {
		{"SOC 2", []string{"CC6.8"}},
		{"ISO 27001", []string{"A.8.8"}},
		{"NIST CSF", []string{"PR.IP-12", "ID.RA-1"}},
		{"CIS v8", []string{"7.4", "7.5"}},
	},
}
