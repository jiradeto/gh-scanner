package repositoryusecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/AlekSi/pointer"
	"github.com/golang/mock/gomock"
	"github.com/jiradeto/gh-scanner/app/domain/entities"
	repositoryrepo "github.com/jiradeto/gh-scanner/app/infrastructure/repos/repository"
	repositoryrepomocks "github.com/jiradeto/gh-scanner/app/infrastructure/repos/repository/mocks"
	"github.com/jiradeto/gh-scanner/app/testutils"
	repositoryusecase "github.com/jiradeto/gh-scanner/app/usecases/repository"
	"github.com/stretchr/testify/assert"
)

func TestFindAllRepositoriesInputValidate(t *testing.T) {
	type fields struct {
		Name            *string
		URL             *string
		Offset          *int
		Limit           *int
		FromCreatedDate *time.Time
		ToCreatedDate   *time.Time
	}
	longText := `sampletextsampletextsampletextsampletextsampletextsampletextsampletext
	sampletextsampletextsampletextsampletextsampletextsampletextsampletextsampletextsamp
	letextsampletextsampletextsampletextsampletextsampletextsampletextsampletextsamplete
	xtsampletextsampletextsampletextsampletextsampletextsampletextsampletextsampletextsa
	mpletextsampletextsampletextsampletextsampletextsampletextsampletextsampletextsample
	textsampletextsampletextsampletextsampletextsampletextsampletextsampletextsampletext
	sampletextsampletextsampletextsampletextsampletextsampletextsampletextsampletextsamp
	tsampletextsampletextsampletextsampletextsampletexttsampletextsampletextsampletextsa`
	now := time.Now()
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "offset is less than one",
			fields: fields{
				Offset: pointer.ToInt(0),
			},
			wantErr: true,
		},
		{
			name: "offset is less than one",
			fields: fields{
				Offset: pointer.ToInt(0),
			},
			wantErr: true,
		},
		{
			name: "offset is greater than max limit",
			fields: fields{
				Offset: pointer.ToInt(300),
			},
			wantErr: true,
		},
		{
			name: "name is longer than max limit",
			fields: fields{
				Name: pointer.ToString(longText),
			},
			wantErr: true,
		},
		{
			name: "url is longer than max limit",
			fields: fields{
				URL: pointer.ToString(longText),
			},
			wantErr: true,
		},
		{
			name: "happy case",
			fields: fields{
				Name:            pointer.ToString("foo_name"),
				URL:             pointer.ToString("foo_url"),
				Offset:          pointer.ToInt(1),
				Limit:           pointer.ToInt(1),
				FromCreatedDate: &now,
				ToCreatedDate:   &now,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &repositoryusecase.FindAllRepositoriesInput{
				Name:            tt.fields.Name,
				URL:             tt.fields.URL,
				Offset:          tt.fields.Offset,
				Limit:           tt.fields.Limit,
				FromCreatedDate: tt.fields.FromCreatedDate,
				ToCreatedDate:   tt.fields.ToCreatedDate,
			}
			if err := f.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_useCase_FindAllRepositories(t *testing.T) {
	type args struct {
		ctx   context.Context
		input repositoryusecase.FindAllRepositoriesInput
	}
	repositories := []*entities.Repository{
		{
			ID:   pointer.ToString("foo_id"),
			Name: pointer.ToString("foo_name"),
			URL:  pointer.ToString("https://github.com/example"),
		},
		{
			ID:   pointer.ToString("bar_id"),
			Name: pointer.ToString("bar_name"),
			URL:  pointer.ToString("https://github.com/example"),
		},
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []*entities.Repository
		wantErr error
	}{
		{
			name: "success",
			fields: fields{
				CreateRepositoryRepo: func(ctrl *gomock.Controller) repositoryrepo.Repo {
					mock := repositoryrepomocks.NewMocks(ctrl)
					mock.EXPECT().FindAllRepositories(gomock.Any(), repositoryrepo.FindAllRepositoriesInput{
						Name: pointer.ToString("foo_name"),
						URL:  pointer.ToString("https://github.com/example"),
					}).Return(repositories, nil)
					return mock
				},
			},
			args: args{
				ctx: context.Background(),
				input: repositoryusecase.FindAllRepositoriesInput{
					Name: pointer.ToString("foo_name"),
					URL:  pointer.ToString("https://github.com/example"),
				},
			},
			want:    repositories,
			wantErr: nil,
		},
		{
			name: "error from repository",
			fields: fields{
				CreateRepositoryRepo: func(ctrl *gomock.Controller) repositoryrepo.Repo {
					mock := repositoryrepomocks.NewMocks(ctrl)
					mock.EXPECT().FindAllRepositories(gomock.Any(), repositoryrepo.FindAllRepositoriesInput{
						Name: pointer.ToString("foo_name"),
						URL:  pointer.ToString("https://github.com/example"),
					}).Return(nil, errors.New("mock error"))
					return mock
				},
			},
			args: args{
				ctx: context.Background(),
				input: repositoryusecase.FindAllRepositoriesInput{
					Name: pointer.ToString("foo_name"),
					URL:  pointer.ToString("https://github.com/example"),
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

			got, err := h.uc.FindAllRepositories(tt.args.ctx, tt.args.input)

			if !testutils.AssertError(t, "useCase.FindAllRepositories()", err, tt.wantErr) {
				return
			}
			assert.Equal(t, tt.want, got, "useCase.FindAllRepositories()")
		})
	}
}
