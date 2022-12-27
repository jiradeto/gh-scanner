package repositoryusecase_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/jiradeto/gh-scanner/app/infrastructure/interfaces/messagequeue"
	messagequeuemocks "github.com/jiradeto/gh-scanner/app/infrastructure/interfaces/messagequeue/mocks"
	repositoryrepo "github.com/jiradeto/gh-scanner/app/infrastructure/repos/repository"
	repositoryrepomocks "github.com/jiradeto/gh-scanner/app/infrastructure/repos/repository/mocks"
	repositoryusecase "github.com/jiradeto/gh-scanner/app/usecases/repository"
)

type fields struct {
	CreateRepositoryRepo     func(ctrl *gomock.Controller) repositoryrepo.Repo
	CreateMessageQueueClient func(ctrl *gomock.Controller) messagequeue.MessageQueueClient
}

var defautlFields = fields{
	CreateRepositoryRepo: func(ctrl *gomock.Controller) repositoryrepo.Repo {
		return repositoryrepomocks.NewMocks(ctrl)
	},
	CreateMessageQueueClient: func(ctrl *gomock.Controller) messagequeue.MessageQueueClient {
		return messagequeuemocks.NewMocks(ctrl)
	},
}

type helper struct {
	uc   repositoryusecase.UseCase
	ctrl *gomock.Controller
	done func()
}

func initTest(t *testing.T, f *fields) *helper {
	ctrl := gomock.NewController(t)

	if f.CreateRepositoryRepo == nil {
		f.CreateRepositoryRepo = defautlFields.CreateRepositoryRepo
	}
	if f.CreateMessageQueueClient == nil {
		f.CreateMessageQueueClient = defautlFields.CreateMessageQueueClient
	}
	return &helper{
		uc: repositoryusecase.New(
			f.CreateRepositoryRepo(ctrl),
			f.CreateMessageQueueClient(ctrl),
		),
		ctrl: ctrl,
		done: func() {
			ctrl.Finish()
		},
	}
}
