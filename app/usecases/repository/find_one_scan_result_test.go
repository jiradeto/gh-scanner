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

func TestFindOneScanResultInputValidate(t *testing.T) {
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
			f := &repositoryusecase.FindOneScanResultInput{
				ID: tt.fields.ID,
			}
			if err := f.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_useCase_FindOneScanResult(t *testing.T) {
	type args struct {
		ctx   context.Context
		input repositoryusecase.FindOneScanResultInput
	}
	mockUUID := "7827f71a-df70-4c16-964f-d65836ec4312"
	scanResult := entities.ScanResult{
		ID:     pointer.ToString(mockUUID),
		Status: entities.ScanResultStatusQueued,
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
				CreateRepositoryRepo: func(ctrl *gomock.Controller) repositoryrepo.Repo {
					mock := repositoryrepomocks.NewMocks(ctrl)
					mock.EXPECT().FindOneScanResult(gomock.Any(), repositoryrepo.FindOneScanResultInput{
						ID: pointer.ToString(mockUUID),
					}).Return(&scanResult, nil)
					return mock
				},
			},
			args: args{
				ctx: context.Background(),
				input: repositoryusecase.FindOneScanResultInput{
					ID: pointer.ToString(mockUUID),
				},
			},
			want:    &scanResult,
			wantErr: nil,
		},
		{
			name: "failed to find scan result",
			fields: fields{
				CreateRepositoryRepo: func(ctrl *gomock.Controller) repositoryrepo.Repo {
					mock := repositoryrepomocks.NewMocks(ctrl)
					mock.EXPECT().FindOneScanResult(gomock.Any(), repositoryrepo.FindOneScanResultInput{
						ID: pointer.ToString(mockUUID),
					}).Return(nil, errors.New("mock error"))
					return mock
				},
			},
			args: args{
				ctx: context.Background(),
				input: repositoryusecase.FindOneScanResultInput{
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

			got, err := h.uc.FindOneScanResult(tt.args.ctx, tt.args.input)

			if !testutils.AssertError(t, "useCase.FindOneScanResult()", err, tt.wantErr) {
				return
			}
			assert.Equal(t, tt.want, got, "useCase.FindOneScanResult()")
		})
	}
}
