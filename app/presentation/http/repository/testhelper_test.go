package repositoryhttp_test

import (
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
)

type fields struct {
}

var defaultFields = fields{}

type helper struct {
	ctrl *gomock.Controller
	done func()
}

func initTest(t *testing.T, f *fields) *helper {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)

	return &helper{

		ctrl: ctrl,
		done: func() {
			ctrl.Finish()
		},
	}
}
