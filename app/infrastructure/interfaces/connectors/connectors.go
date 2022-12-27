package connectors

import (
	"strings"
	"sync"

	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
)

var (
	mut sync.Mutex
)

// DatabaseConfig is a database config
type DatabaseConfig struct {
	Host      string `validate:"required,min=1"`
	Port      string `validate:"required,min=1"`
	User      string `validate:"required,min=1"`
	Password  string `validate:"required,min=1"`
	DB        string `validate:"required,min=1"`
	DBLogMode bool
}

// Validate is a function to validate function input
func (db *DatabaseConfig) Validate() error {
	v := validator.New()
	if err := v.Struct(db); err != nil {
		ve, _ := err.(validator.ValidationErrors)
		var errFields []string
		for _, fe := range ve {
			errFields = append(errFields, fe.Field())
		}
		return errors.Errorf("%s is not defined", strings.Join(errFields, ","))
	}
	return nil
}
