package entities

type ErrorStatus struct {
	Code    uint     `json:"code" form:"error_code"`
	Message string   `json:"message" form:"error_message"`
	Details []string `json:"details,omitempty"`
	Err     error    `json:"-"`
}

// Error is implement function from error interface to get error string
func (e *ErrorStatus) Error() string {
	if e.Err != nil {
		return e.Err.Error()
	}
	return ""
}
