package repositoryusecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/AlekSi/pointer"
	"github.com/golang/mock/gomock"
	"github.com/jiradeto/gh-scanner/app/entities"
	repositoryrepo "github.com/jiradeto/gh-scanner/app/infrastructure/repos/repository"
	repositoryrepomocks "github.com/jiradeto/gh-scanner/app/infrastructure/repos/repository/mocks"
	"github.com/jiradeto/gh-scanner/app/testutils"
	repositoryusecase "github.com/jiradeto/gh-scanner/app/usecases/repository"
	"github.com/stretchr/testify/assert"
)

func TestFindOneRepositoryInputValidate(t *testing.T) {
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
			f := &repositoryusecase.FindOneRepositoryInput{
				ID: tt.fields.ID,
			}
			if err := f.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_useCase_FindOneRepository(t *testing.T) {
	type args struct {
		ctx   context.Context
		input repositoryusecase.FindOneRepositoryInput
	}
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
		want    *entities.Repository
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
					return mock
				},
			},
			args: args{
				ctx: context.Background(),
				input: repositoryusecase.FindOneRepositoryInput{
					ID: pointer.ToString(mockUUID),
				},
			},
			want:    &repository,
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
				input: repositoryusecase.FindOneRepositoryInput{
					ID: pointer.ToString(mockUUID),
				},
			},
			want:    nil,
			wantErr: errors.New("mock error"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := initTest(t, &tt.fields)
			defer h.done()

			got, err := h.uc.FindOneRepository(tt.args.ctx, tt.args.input)

			if !testutils.AssertError(t, "useCase.FindOneRepository()", err, tt.wantErr) {
				return
			}
			assert.Equal(t, tt.want, got, "useCase.FindOneRepository()")
		})
	}
}
