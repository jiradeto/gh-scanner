package models

import (
	"encoding/json"
	"time"

	"github.com/jinzhu/copier"
	"github.com/jinzhu/gorm/dialects/postgres"
	"github.com/jiradeto/gh-scanner/app/entities"
)

// ScanResult models
type ScanResult struct {
	ID           string
	Status       string
	RepositoryID string
	Findings     *postgres.Jsonb
	QueuedAt     time.Time
	ScanningAt   time.Time
	FinishedAt   time.Time
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// TableName is table name of ScanResult
func (ScanResult) TableName() string {
	return "scan_results"
}

// FromEntity converts entity to model
func (w *ScanResult) FromEntity(e *entities.ScanResult) (*ScanResult, error) {
	var model ScanResult

	err := copier.Copy(&model, &e)
	if err != nil {
		return nil, err
	}

	if e.Findings != nil {
		rawMessage, err := json.Marshal(e.Findings)
		if err != nil {
			return nil, err
		}
		model.Findings = &postgres.Jsonb{
			RawMessage: rawMessage,
		}
	}

	model.Status = e.Status.String()
	return &model, nil
}

// ToEntity converts model to entity
func (w *ScanResult) ToEntity() (*entities.ScanResult, error) {
	var entity entities.ScanResult

	err := copier.Copy(&entity, &w)
	if err != nil {
		return nil, err
	}

	if w.Findings != nil && w.Findings.RawMessage != nil {
		var findings []entities.ScanFinding
		if err := json.Unmarshal(w.Findings.RawMessage, &findings); err != nil {
			return nil, err
		}
		entity.Findings = findings
	}

	entity.Status = new(entities.ScanResultStatus).Parse(w.Status)

	return &entity, nil
}

// ScanResults is an array of ScanResult
type ScanResults []*ScanResult

// FromEntities converts models to entities
func (ws ScanResults) FromEntities(es []*entities.ScanResult) ([]*ScanResult, error) {
	var ms []*ScanResult

	for _, e := range es {
		m, err := new(ScanResult).FromEntity(e)
		if err != nil {
			return nil, err
		}
		ms = append(ms, m)
	}

	return ms, nil
}

// ToEntities converts entities to models
func (ws ScanResults) ToEntities() (entities.ScanResults, error) {
	var es entities.ScanResults

	for _, w := range ws {
		e, err := w.ToEntity()
		if err != nil {
			return nil, err
		}
		es = append(es, e)
	}

	return es, nil
}
