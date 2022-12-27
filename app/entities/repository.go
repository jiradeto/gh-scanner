package entities

import "time"

// Repository entity
type Repository struct {
	CreatedAt *time.Time
	ID        *string
	Name      *string
	URL       *string
}

// Repositories is an array of type Repository
type Repositories []*Repository
