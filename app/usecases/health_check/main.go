package healthcheckusecase

import (
	"context"
)

// UseCase is a interface for healthcheckusecase
type UseCase interface {
	CheckServiceVersion(context context.Context) (*CheckServiceVersionResponse, error)
}

type useCase struct {
}

// New is a constructor method of UseCase
func New() UseCase {
	return &useCase{}
}
