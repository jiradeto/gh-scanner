package models

import (
	"time"

	"github.com/jinzhu/copier"
	"github.com/jiradeto/gh-scanner/app/domain/entities"
)

// Repository models
type Repository struct {
	ID        string
	Name      string
	URL       string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// TableName is table name of Repository
func (Repository) TableName() string {
	return "repositories"
}

// FromEntity converts entity to model
func (w *Repository) FromEntity(e *entities.Repository) (*Repository, error) {
	var model Repository

	err := copier.Copy(&model, &e)
	if err != nil {
		return nil, err
	}

	return &model, nil
}

// ToEntity converts model to entity
func (w *Repository) ToEntity() (*entities.Repository, error) {
	var entity entities.Repository

	err := copier.Copy(&entity, &w)
	if err != nil {
		return nil, err
	}

	return &entity, nil
}

// Repositories is an array of Repository
type Repositories []*Repository

// FromEntities converts models to entities
func (ws Repositories) FromEntities(es []*entities.Repository) ([]*Repository, error) {
	var ms []*Repository

	for _, e := range es {
		m, err := new(Repository).FromEntity(e)
		if err != nil {
			return nil, err
		}
		ms = append(ms, m)
	}

	return ms, nil
}

// ToEntities converts entities to models
func (ws Repositories) ToEntities() (entities.Repositories, error) {
	var es entities.Repositories

	for _, w := range ws {
		e, err := w.ToEntity()
		if err != nil {
			return nil, err
		}
		es = append(es, e)
	}

	return es, nil
}
