package gerrors

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/jiradeto/gh-scanner/app/constants"
)

// ParameterError Invalid parameter
type ParameterError struct {
	Code            uint
	Message         string
	ValidatorErrors *validator.ValidationErrors
}

// Wrap ...
func (e ParameterError) Wrap(err error) error {
	return fmt.Errorf("%w: %v", e, err)
}

// Error ...
func (e ParameterError) Error() string {
	return e.Message
}

// Is ...
func (e ParameterError) Is(target error) bool {
	_, ok := target.(*ParameterError)
	if !ok {
		return false
	}
	return true
}

// UnprocessableEntity Valid parameter but invalid business and etc.
type UnprocessableEntity struct {
	Code    uint
	Message string
}

// Wrap ...
func (e UnprocessableEntity) Wrap(err error) error {
	return fmt.Errorf("%w: %v", e, err)
}

// Error ...
func (e UnprocessableEntity) Error() string {
	return e.Message
}

// Is ...
func (e UnprocessableEntity) Is(target error) bool {
	_, ok := target.(*UnprocessableEntity)
	if !ok {
		return false
	}
	return true
}

// RecordNotFoundError Cannot find resource.
type RecordNotFoundError struct {
	Code    uint
	Message string
}

// Wrap ...
func (e RecordNotFoundError) Wrap(err error) error {
	return fmt.Errorf("%w: %v", e, err)
}

// Error ...
func (e RecordNotFoundError) Error() string {
	return e.Message
}

// Is ...
func (e RecordNotFoundError) Is(target error) bool {
	_, ok := target.(*RecordNotFoundError)
	if !ok {
		return false
	}
	return true
}

func NewInternalError(e error) InternalError {
	msg := constants.ErrorMessageInternalError
	if e != nil {
		msg = e.Error()
	}
	return InternalError{
		Code:    constants.StatusCodeGenericInternalError,
		Message: msg,
	}
}

// InternalError Database error and etc.
type InternalError struct {
	Code    uint
	Message string
}

// Wrap ...
func (e InternalError) Wrap(err error) error {
	return fmt.Errorf("%w: %v", e, err)
}

// Error ...
func (e InternalError) Error() string {
	return e.Message
}

// Is ...
func (e InternalError) Is(target error) bool {
	_, ok := target.(*InternalError)
	if !ok {
		return false
	}
	return true
}

// ExternalError Database error and etc.
type ExternalError struct {
	HTTPStatus int
	Code       uint
	Message    string
}

// Wrap ...
func (e ExternalError) Wrap(err error) error {
	return fmt.Errorf("%w: %v", e, err)
}

// Error ...
func (e ExternalError) Error() string {
	return e.Message
}

// Is ...
func (e ExternalError) Is(target error) bool {
	_, ok := target.(*ExternalError)
	if !ok {
		return false
	}
	return true
}
