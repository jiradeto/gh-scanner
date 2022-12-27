package response

import (
	"errors"
	"net/http"

	"github.com/AlekSi/pointer"
	"github.com/gin-gonic/gin"
	"github.com/jiradeto/gh-scanner/app/constants"
	"github.com/jiradeto/gh-scanner/app/environments"
	"github.com/jiradeto/gh-scanner/app/utils/gerrors"
)

const (
	MessageOK string = "OK"
)

// BaseResponseStatus ...
type BaseResponseStatus struct {
	Code     uint     `json:"code"`
	Messages []string `json:"messages"`
	Details  *string  `json:"details,omitempty"`
}

type BaseSuccessResponse struct {
	Status BaseResponseStatus `json:"status"`
	Data   interface{}        `json:"data"`
}

type BaseErrorResponse struct {
	Status BaseResponseStatus `json:"status"`
}

func ResponseSuccess(c *gin.Context, data interface{}) {
	r := new(BaseSuccessResponse)
	r.Status.Code = constants.StatusCodeSuccess
	r.Status.Messages = append(r.Status.Messages, MessageOK)
	r.Data = data
	c.AbortWithStatusJSON(http.StatusOK, *r)
}

func ResponseError(c *gin.Context, err error) {
	status := BaseResponseStatus{
		Code:     constants.StatusCodeGenericInternalError,
		Messages: []string{},
	}
	httpStatusCode := http.StatusInternalServerError
	if len(err.Error()) > 0 && environments.DevMode {
		status.Details = pointer.ToString(err.Error())
	}
	unwrappedErr := errors.Unwrap(err)
	if unwrappedErr == nil {
		unwrappedErr = err
	}
	switch unwrappedErr.(type) {
	case gerrors.InternalError:
		xerr, _ := unwrappedErr.(gerrors.InternalError)
		httpStatusCode = http.StatusInternalServerError
		status.Code = xerr.Code
		status.Messages = append(status.Messages, xerr.Message)
	case gerrors.ParameterError:
		xerr, _ := unwrappedErr.(gerrors.ParameterError)
		httpStatusCode = http.StatusBadRequest
		status.Code = xerr.Code
		status.Messages = append(status.Messages, xerr.Message)
	case gerrors.RecordNotFoundError:
		xerr, _ := unwrappedErr.(gerrors.RecordNotFoundError)
		httpStatusCode = http.StatusNotFound
		status.Code = xerr.Code
		status.Messages = append(status.Messages, xerr.Message)
	case gerrors.ExternalError:
		xerr, _ := unwrappedErr.(gerrors.ExternalError)
		httpStatusCode = xerr.HTTPStatus
		status.Code = xerr.Code
		status.Messages = append(status.Messages, xerr.Message)
	default:
		status.Messages = append(status.Messages, err.Error())
	}

	c.AbortWithStatusJSON(httpStatusCode, BaseErrorResponse{Status: status})
}
