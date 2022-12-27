package repositoryusecase_test

import (
	"context"
	"testing"

	"github.com/AlekSi/pointer"
	"github.com/golang/mock/gomock"
	"github.com/jiradeto/gh-scanner/app/constants"
	"github.com/jiradeto/gh-scanner/app/domain/entities"
	"github.com/jiradeto/gh-scanner/app/infrastructure/interfaces/messagequeue"
	messagequeuemocks "github.com/jiradeto/gh-scanner/app/infrastructure/interfaces/messagequeue/mocks"
	repositoryrepo "github.com/jiradeto/gh-scanner/app/infrastructure/repos/repository"
	repositoryrepomocks "github.com/jiradeto/gh-scanner/app/infrastructure/repos/repository/mocks"
	"github.com/jiradeto/gh-scanner/app/testutils"
	repositoryusecase "github.com/jiradeto/gh-scanner/app/usecases/repository"
	"github.com/jiradeto/gh-scanner/app/utils/gerrors"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestStartScannerInputValidate(t *testing.T) {
	type fields struct {
		ID *string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name:    "missing required fields",
			fields:  fields{},
			wantErr: true,
		},
		{
			name: "empty ID",
			fields: fields{
				ID: pointer.ToString(""),
			},
			wantErr: true,
		},
		{
			name: "invalid uuid format",
			fields: fields{
				ID: pointer.ToString("foo_uuid"),
			},
			wantErr: true,
		},
		{
			name: "happy case",
			fields: fields{
				ID: pointer.ToString("7827f71a-df70-4c16-964f-d65836ec4312"),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &repositoryusecase.StartScannerInput{
				ID: tt.fields.ID,
			}
			if err := f.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_useCase_StartScanner(t *testing.T) {
	type args struct {
		ctx   context.Context
		input repositoryusecase.StartScannerInput
	}
	const errLocation = "repositoryUsecase/StartScanner: %s"
	mockRepositoryID := "7827f71a-df70-4c16-964f-d65836ec4312"
	mockScanResultID := "34a1300e-285e-424e-9c9d-14d2b3dbc6e9"
	repository := entities.Repository{
		ID:   pointer.ToString(mockRepositoryID),
		Name: pointer.ToString("foo_name"),
		URL:  pointer.ToString("https://github.com/example"),
	}
	scanResult := entities.ScanResult{
		ID: pointer.ToString(mockScanResultID),
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *entities.ScanResult
		wantErr error
	}{
		{
			name: "success",
			fields: fields{
				CreateMessageQueueClient: func(ctrl *gomock.Controller) messagequeue.MessageQueueClient {
					mock := messagequeuemocks.NewMocks(ctrl)
					mock.EXPECT().PublishMessage(&messagequeue.StartScannerMessage{
						ResultId: *scanResult.ID,
						URL:      *repository.URL,
					}).Return(nil)
					return mock
				},
				CreateRepositoryRepo: func(ctrl *gomock.Controller) repositoryrepo.Repo {
					mock := repositoryrepomocks.NewMocks(ctrl)
					mock.EXPECT().FindOneRepository(gomock.Any(), repositoryrepo.FindOneRepositoryInput{
						ID: pointer.ToString(mockRepositoryID),
					}).Return(&repository, nil)
					mock.EXPECT().CreateOneScanResult(gomock.Any(), repositoryrepo.CreateOneScanResultInput{
						ScanResultEntity: &entities.ScanResult{
							RepositoryID: repository.ID,
							Status:       entities.ScanResultStatusQueued,
						},
					}).Return(&scanResult, nil)
					return mock
				},
			},
			args: args{
				ctx: context.Background(),
				input: repositoryusecase.StartScannerInput{
					ID: pointer.ToString(mockRepositoryID),
				},
			},
			want:    &scanResult,
			wantErr: nil,
		},
		{
			name: "failed to create one scan result",
			fields: fields{
				CreateRepositoryRepo: func(ctrl *gomock.Controller) repositoryrepo.Repo {
					mock := repositoryrepomocks.NewMocks(ctrl)
					mock.EXPECT().FindOneRepository(gomock.Any(), repositoryrepo.FindOneRepositoryInput{
						ID: pointer.ToString(mockRepositoryID),
					}).Return(&repository, nil)
					mock.EXPECT().CreateOneScanResult(gomock.Any(), repositoryrepo.CreateOneScanResultInput{
						ScanResultEntity: &entities.ScanResult{
							RepositoryID: repository.ID,
							Status:       entities.ScanResultStatusQueued,
						},
					}).Return(nil, errors.New("mock error"))
					return mock
				},
			},
			args: args{
				ctx: context.Background(),
				input: repositoryusecase.StartScannerInput{
					ID: pointer.ToString(mockRepositoryID),
				},
			},
			want: nil,
			wantErr: gerrors.InternalError{
				Code:    constants.StatusCodeDatabaseError,
				Message: constants.ErrorMessageDatabaseError,
			}.Wrap(errors.Wrapf(errors.New("mock error"), errLocation, "unable to create scan result")),
		},
		{
			name: "error to send to message queue",
			fields: fields{
				CreateMessageQueueClient: func(ctrl *gomock.Controller) messagequeue.MessageQueueClient {
					mock := messagequeuemocks.NewMocks(ctrl)
					mock.EXPECT().PublishMessage(&messagequeue.StartScannerMessage{
						ResultId: *scanResult.ID,
						URL:      *repository.URL,
					}).Return(errors.New("mock error"))
					return mock
				},
				CreateRepositoryRepo: func(ctrl *gomock.Controller) repositoryrepo.Repo {
					mock := repositoryrepomocks.NewMocks(ctrl)
					mock.EXPECT().FindOneRepository(gomock.Any(), repositoryrepo.FindOneRepositoryInput{
						ID: pointer.ToString(mockRepositoryID),
					}).Return(&repository, nil)
					mock.EXPECT().CreateOneScanResult(gomock.Any(), repositoryrepo.CreateOneScanResultInput{
						ScanResultEntity: &entities.ScanResult{
							RepositoryID: repository.ID,
							Status:       entities.ScanResultStatusQueued,
						},
					}).Return(&scanResult, nil)
					return mock
				},
			},
			args: args{
				ctx: context.Background(),
				input: repositoryusecase.StartScannerInput{
					ID: pointer.ToString(mockRepositoryID),
				},
			},
			want: nil,
			wantErr: gerrors.InternalError{
				Code:    constants.StatusCodeDatabaseError,
				Message: constants.ErrorMessageDatabaseError,
			}.Wrap(errors.Wrapf(errors.New("mock error"), errLocation, "failed to publish message queue")),
		},

		{
			name: "failed to find one repository",
			fields: fields{
				CreateRepositoryRepo: func(ctrl *gomock.Controller) repositoryrepo.Repo {
					mock := repositoryrepomocks.NewMocks(ctrl)
					mock.EXPECT().FindOneRepository(gomock.Any(), repositoryrepo.FindOneRepositoryInput{
						ID: pointer.ToString(mockRepositoryID),
					}).Return(nil, errors.New("mock error"))
					return mock
				},
			},
			args: args{
				ctx: context.Background(),
				input: repositoryusecase.StartScannerInput{
					ID: pointer.ToString(mockRepositoryID),
				},
			},
			wantErr: gerrors.InternalError{
				Code:    constants.StatusCodeDatabaseError,
				Message: constants.ErrorMessageDatabaseError,
			}.Wrap(errors.Wrapf(errors.New("mock error"), errLocation, "unable to find repository")),
		},
		{
			name: "when record not found",
			fields: fields{
				CreateRepositoryRepo: func(ctrl *gomock.Controller) repositoryrepo.Repo {
					mock := repositoryrepomocks.NewMocks(ctrl)
					mock.EXPECT().FindOneRepository(gomock.Any(), repositoryrepo.FindOneRepositoryInput{
						ID: pointer.ToString(mockRepositoryID),
					}).Return(nil, nil)
					return mock
				},
			},
			args: args{
				ctx: context.Background(),
				input: repositoryusecase.StartScannerInput{
					ID: pointer.ToString(mockRepositoryID),
				},
			},
			wantErr: gerrors.RecordNotFoundError{
				Code:    constants.StatusCodeEntryNotFound,
				Message: constants.ErrorMessageNotFound,
			}.Wrap(errors.Errorf(errLocation, "not found repository")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := initTest(t, &tt.fields)
			defer h.done()

			got, err := h.uc.StartScanner(tt.args.ctx, tt.args.input)

			if !testutils.AssertError(t, "useCase.StartScanner()", err, tt.wantErr) {
				return
			}
			assert.Equal(t, tt.want, got, "useCase.StartScanner()")
		})
	}
}
