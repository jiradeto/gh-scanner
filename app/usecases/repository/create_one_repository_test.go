package repositoryusecase_test

import (
	"context"
	"testing"

	"github.com/AlekSi/pointer"
	"github.com/golang/mock/gomock"
	"github.com/jiradeto/gh-scanner/app/domain/entities"
	repositoryrepo "github.com/jiradeto/gh-scanner/app/infrastructure/repos/repository"
	repositoryrepomocks "github.com/jiradeto/gh-scanner/app/infrastructure/repos/repository/mocks"
	"github.com/jiradeto/gh-scanner/app/testutils"
	repositoryusecase "github.com/jiradeto/gh-scanner/app/usecases/repository"
	"github.com/stretchr/testify/assert"
)

func TestCreateOneRepositoryInputValidate(t *testing.T) {
	type fields struct {
		Name *string
		URL  *string
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
			name: "missing name required fields",
			fields: fields{
				Name: pointer.ToString("foo"),
			},
			wantErr: true,
		},
		{
			name: "when name is empty",
			fields: fields{
				Name: pointer.ToString(""),
			},
			wantErr: true,
		},
		{
			name: "when provide invalid url",
			fields: fields{
				Name: pointer.ToString("foo_name"),
				URL:  pointer.ToString("foo_url"),
			},
			wantErr: true,
		},
		{
			name: "when given url is not github",
			fields: fields{
				Name: pointer.ToString("foo_name"),
				URL:  pointer.ToString("https://example.com/foo"),
			},
			wantErr: true,
		},
		{
			name: "success case",
			fields: fields{
				Name: pointer.ToString("foo_name"),
				URL:  pointer.ToString("https://github.com/foo"),
			},
			wantErr: false,
		},
		{
			name: "success case on non-https",
			fields: fields{
				Name: pointer.ToString("foo_name"),
				URL:  pointer.ToString("http://github.com/foo"),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &repositoryusecase.CreateOneRepositoryInput{
				Name: tt.fields.Name,
				URL:  tt.fields.URL,
			}
			if err := f.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_useCase_CreateOneRepository(t *testing.T) {
	const errLocation = "[card issuance usecase/create or update card issuance] %s"
	type args struct {
		ctx   context.Context
		input repositoryusecase.CreateOneRepositoryInput
	}
	repository := entities.Repository{
		ID:   pointer.ToString("foo_program_id"),
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
					mock.EXPECT().CreateOneRepository(gomock.Any(), repositoryrepo.CreateOneRepositoryInput{
						RepositoryEntity: &entities.Repository{
							Name: pointer.ToString("foo_name"),
							URL:  pointer.ToString("https://github.com/example"),
						},
					}).Return(&repository, nil)
					return mock
				},
			},
			args: args{
				ctx: context.Background(),
				input: repositoryusecase.CreateOneRepositoryInput{
					Name: pointer.ToString("foo_name"),
					URL:  pointer.ToString("https://github.com/example"),
				},
			},
			want:    &repository,
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := initTest(t, &tt.fields)
			defer h.done()

			got, err := h.uc.CreateOneRepository(tt.args.ctx, tt.args.input)

			if !testutils.AssertError(t, "useCase.CreateOneRepository()", err, tt.wantErr) {
				return
			}
			assert.Equal(t, tt.want, got, "useCase.CreateOneRepository()")
		})
	}
}
