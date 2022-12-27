package entities

import (
	"time"
)

type Severity string

const (
	SeverityLow    Severity = "LOW"
	SeverityMedium Severity = "MED"
	SeverityHigh   Severity = "HIGH"
)

type Begin struct {
	Line int64 `json:"line"`
}

type Positions struct {
	Begin Begin `json:"begin"`
}

type Location struct {
	Path      string    `json:"path"`
	Positions Positions `json:"positions"`
}

type Metadata struct {
	Description string   `json:"description"`
	Severity    Severity `json:"severity"`
}

type ScanFinding struct {
	Type     string   `json:"type"`
	RuleID   string   `json:"ruleId"`
	Location Location `json:"location"`
	Metadata Metadata `json:"metadata"`
}

// ScanResult entity
type ScanResult struct {
	ID           *string          `json:"id"`
	Status       ScanResultStatus `json:"status"`
	RepositoryID *string          `json:"repositoryID"`
	Findings     []ScanFinding    `json:"findings"`
	CreatedAt    *time.Time       `json:"createdAt"`
	UpdatedAt    *time.Time       `json:"updatedAt"`
	QueuedAt     *time.Time       `json:"queuedAt"`
	ScanningAt   *time.Time       `json:"scanningAt"`
	FinishedAt   *time.Time       `json:"finishedAt"`
}

// ScanResults is an array of type of ScanResult
type ScanResults []*ScanResult

func (source ScanResult) MapResponse() map[string]interface{} {
	dest := map[string]interface{}{
		"id":           source.ID,
		"status":       source.Status.String(),
		"RepositoryID": source.RepositoryID,
		"CreatedAt":    source.CreatedAt,
		"UpdatedAt":    source.UpdatedAt,
		"QueuedAt":     source.QueuedAt,
		"ScanningAt":   source.ScanningAt,
		"FinishedAt":   source.FinishedAt,
		"Findings":     source.Findings,
	}

	return dest
}
