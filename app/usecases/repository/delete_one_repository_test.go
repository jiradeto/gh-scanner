package repositoryusecase_test

import (
	"context"
	"testing"

	"github.com/AlekSi/pointer"
	"github.com/golang/mock/gomock"
	"github.com/jiradeto/gh-scanner/app/constants"
	"github.com/jiradeto/gh-scanner/app/domain/entities"
	repositoryrepo "github.com/jiradeto/gh-scanner/app/infrastructure/repos/repository"
	repositoryrepomocks "github.com/jiradeto/gh-scanner/app/infrastructure/repos/repository/mocks"
	"github.com/jiradeto/gh-scanner/app/testutils"
	repositoryusecase "github.com/jiradeto/gh-scanner/app/usecases/repository"
	"github.com/jiradeto/gh-scanner/app/utils/gerrors"
	"github.com/pkg/errors"
)

func TestDeleteOneRepositoryInputValidate(t *testing.T) {
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
			f := &repositoryusecase.DeleteOneRepositoryInput{
				ID: tt.fields.ID,
			}
			if err := f.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_useCase_DeleteOneRepository(t *testing.T) {
	type args struct {
		ctx   context.Context
		input repositoryusecase.DeleteOneRepositoryInput
	}
	const errLocation = "repositoryUsecase/DeleteOneRepository: %s"
	mockUUID := "7827f71a-df70-4c16-964f-d65836ec4312"
	repository := entities.Repository{
		ID:   pointer.ToString(mockUUID),
		Name: pointer.ToString("foo_name"),
		URL:  pointer.ToString("https://github.com/example"),
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr error
	}{
		{
			name: "success",
			fields: fields{
				CreateRepositoryRepo: func(ctrl *gomock.Controller) repositoryrepo.Repo {
					mock := repositoryrepomocks.NewMocks(ctrl)
					mock.EXPECT().FindOneRepository(gomock.Any(), repositoryrepo.FindOneRepositoryInput{
						ID: pointer.ToString(mockUUID),
					}).Return(&repository, nil)
					mock.EXPECT().DeleteOneRepository(gomock.Any(), repositoryrepo.DeleteOneRepositoryInput{
						RepositoryEntity: &repository,
					}).Return(nil)
					return mock
				},
			},
			args: args{
				ctx: context.Background(),
				input: repositoryusecase.DeleteOneRepositoryInput{
					ID: pointer.ToString(mockUUID),
				},
			},
			wantErr: nil,
		},
		{
			name: "failed to find one repository",
			fields: fields{
				CreateRepositoryRepo: func(ctrl *gomock.Controller) repositoryrepo.Repo {
					mock := repositoryrepomocks.NewMocks(ctrl)
					mock.EXPECT().FindOneRepository(gomock.Any(), repositoryrepo.FindOneRepositoryInput{
						ID: pointer.ToString(mockUUID),
					}).Return(nil, errors.New("mock error"))
					return mock
				},
			},
			args: args{
				ctx: context.Background(),
				input: repositoryusecase.DeleteOneRepositoryInput{
					ID: pointer.ToString(mockUUID),
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
						ID: pointer.ToString(mockUUID),
					}).Return(nil, nil)
					return mock
				},
			},
			args: args{
				ctx: context.Background(),
				input: repositoryusecase.DeleteOneRepositoryInput{
					ID: pointer.ToString(mockUUID),
				},
			},
			wantErr: gerrors.RecordNotFoundError{
				Code:    constants.StatusCodeEntryNotFound,
				Message: constants.ErrorMessageNotFound,
			}.Wrap(errors.Errorf(errLocation, "not found repository")),
		},
		{
			name: "failed to delete one repository",
			fields: fields{
				CreateRepositoryRepo: func(ctrl *gomock.Controller) repositoryrepo.Repo {
					mock := repositoryrepomocks.NewMocks(ctrl)
					mock.EXPECT().FindOneRepository(gomock.Any(), repositoryrepo.FindOneRepositoryInput{
						ID: pointer.ToString(mockUUID),
					}).Return(&repository, nil)
					mock.EXPECT().DeleteOneRepository(gomock.Any(), repositoryrepo.DeleteOneRepositoryInput{
						RepositoryEntity: &repository,
					}).Return(errors.New("mock error"))
					return mock
				},
			},
			args: args{
				ctx: context.Background(),
				input: repositoryusecase.DeleteOneRepositoryInput{
					ID: pointer.ToString(mockUUID),
				},
			},
			wantErr: errors.New("mock error"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := initTest(t, &tt.fields)
			defer h.done()

			err := h.uc.DeleteOneRepository(tt.args.ctx, tt.args.input)

			if !testutils.AssertError(t, "useCase.DeleteOneRepository()", err, tt.wantErr) {
				return
			}
		})
	}
}
