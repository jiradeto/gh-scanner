package entities

// ScanResultStatus is enum type for scan result status
type ScanResultStatus int

// ScanResultStatus enum constants
const (
	ScanResultStatusQueued ScanResultStatus = iota + 1
	ScanResultStatusInProgress
	ScanResultStatusSusccess
	ScanResultStatusFailure
)

var ScanResultStatusIndexMapper = map[string]ScanResultStatus{
	"queued":      ScanResultStatusQueued,
	"in_progress": ScanResultStatusInProgress,
	"success":     ScanResultStatusSusccess,
	"failure":     ScanResultStatusFailure,
}

var ScanResultStatusStringMapper = map[ScanResultStatus]string{
	ScanResultStatusQueued:     "queued",
	ScanResultStatusInProgress: "in_progress",
	ScanResultStatusSusccess:   "success",
	ScanResultStatusFailure:    "failure",
}

// Parse converts string to ScanResultStatus
func (c ScanResultStatus) Parse(ScanResultStatus string) ScanResultStatus {
	return ScanResultStatusIndexMapper[ScanResultStatus]
}

func (c ScanResultStatus) String() string {
	return ScanResultStatusStringMapper[c]
}

// Is is a function to check whether ScanResultStatus is equal to expected ScanResultStatus
func (c ScanResultStatus) Is(expected ScanResultStatus) bool {
	return c.String() == expected.String()
}
