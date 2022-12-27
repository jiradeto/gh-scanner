package testutils

import (
	"errors"
	"reflect"
	"strings"
	"testing"
)

// AssertError is a function for asserting if actual error and want error are equal
func AssertError(t *testing.T, prefix string, actualErr error, wantErr error) bool {

	if actualErr == nil {
		if wantErr == nil {
			return true
		}

		t.Errorf("%s error = %v, wantErr %v", prefix, actualErr, wantErr)
		return false
	}
	if errors.Unwrap(actualErr) != nil {
		if errors.Unwrap(wantErr) != nil {
			actualUnwrapErr := errors.Unwrap(actualErr)
			wantUnwrapErr := errors.Unwrap(wantErr)
			if reflect.TypeOf(actualUnwrapErr) != reflect.TypeOf(wantUnwrapErr) {
				t.Errorf("%s error type = %v, wantErr %v", prefix, reflect.TypeOf(actualUnwrapErr), reflect.TypeOf(wantUnwrapErr))
				return false
			}
		} else {
			t.Errorf("%s error type = %v, wantErr %v", prefix, reflect.TypeOf(actualErr), reflect.TypeOf(wantErr))
			return false
		}
	} else {
		if errors.Unwrap(wantErr) == nil {
			if reflect.TypeOf(actualErr) != reflect.TypeOf(wantErr) {
				t.Errorf("%s error type = %v, wantErr %v", prefix, reflect.TypeOf(actualErr), reflect.TypeOf(wantErr))
				return false
			}
		} else {
			t.Errorf("%s error type = %v, wantErr %v", prefix, reflect.TypeOf(actualErr), reflect.TypeOf(wantErr))
			return false
		}
	}
	if !strings.Contains(actualErr.Error(), wantErr.Error()) {
		t.Errorf("\n%s \nerror = %v \nwantErr %v", prefix, actualErr, wantErr)
		return false
	}
	return true
}
