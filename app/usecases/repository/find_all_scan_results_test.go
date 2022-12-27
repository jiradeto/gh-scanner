package repositoryusecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/AlekSi/pointer"
	"github.com/golang/mock/gomock"
	"github.com/jiradeto/gh-scanner/app/entities"
	repositoryrepo "github.com/jiradeto/gh-scanner/app/infrastructure/repos/repository"
	repositoryrepomocks "github.com/jiradeto/gh-scanner/app/infrastructure/repos/repository/mocks"
	"github.com/jiradeto/gh-scanner/app/testutils"
	repositoryusecase "github.com/jiradeto/gh-scanner/app/usecases/repository"
	"github.com/stretchr/testify/assert"
)

func TestFindAllScanResultsInputValidate(t *testing.T) {
	type fields struct {
		RepositoryID    *string
		Limit           *int
		FromCreatedDate *time.Time
		ToCreatedDate   *time.Time
	}
	now := time.Now()
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "repository id is not valid uuid",
			fields: fields{
				RepositoryID:    pointer.ToString("foo_id"),
				Limit:           pointer.ToInt(1),
				FromCreatedDate: &now,
				ToCreatedDate:   &now,
			},
			wantErr: true,
		},
		{
			name: "happy case",
			fields: fields{
				RepositoryID:    pointer.ToString("7827f71a-df70-4c16-964f-d65836ec4312"),
				Limit:           pointer.ToInt(1),
				FromCreatedDate: &now,
				ToCreatedDate:   &now,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &repositoryusecase.FindAllScanResultsInput{
				RepositoryID:    tt.fields.RepositoryID,
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

func Test_useCase_FindAllScanResults(t *testing.T) {
	type args struct {
		ctx   context.Context
		input repositoryusecase.FindAllScanResultsInput
	}
	mockUUID := "7827f71a-df70-4c16-964f-d65836ec4312"
	scanResults := []*entities.ScanResult{
		{
			ID:     pointer.ToString("foo_id"),
			Status: entities.ScanResultStatusQueued,
		},
		{
			ID:     pointer.ToString("bar_id"),
			Status: entities.ScanResultStatusQueued,
		},
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []*entities.ScanResult
		wantErr error
	}{
		{
			name: "success",
			fields: fields{
				CreateRepositoryRepo: func(ctrl *gomock.Controller) repositoryrepo.Repo {
					mock := repositoryrepomocks.NewMocks(ctrl)
					mock.EXPECT().FindAllScanResults(gomock.Any(), repositoryrepo.FindAllScanResultsInput{
						RepositoryID: pointer.ToString(mockUUID),
					}).Return(scanResults, nil)
					return mock
				},
			},
			args: args{
				ctx: context.Background(),
				input: repositoryusecase.FindAllScanResultsInput{
					RepositoryID: pointer.ToString(mockUUID),
				},
			},
			want:    scanResults,
			wantErr: nil,
		},
		{
			name: "error from repository",
			fields: fields{
				CreateRepositoryRepo: func(ctrl *gomock.Controller) repositoryrepo.Repo {
					mock := repositoryrepomocks.NewMocks(ctrl)
					mock.EXPECT().FindAllScanResults(gomock.Any(), repositoryrepo.FindAllScanResultsInput{
						RepositoryID: pointer.ToString(mockUUID),
					}).Return(nil, errors.New("mock error"))
					return mock
				},
			},
			args: args{
				ctx: context.Background(),
				input: repositoryusecase.FindAllScanResultsInput{
					RepositoryID: pointer.ToString(mockUUID),
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

			got, err := h.uc.FindAllScanResults(tt.args.ctx, tt.args.input)

			if !testutils.AssertError(t, "useCase.FindAllScanResults()", err, tt.wantErr) {
				return
			}
			assert.Equal(t, tt.want, got, "useCase.FindAllScanResults()")
		})
	}
}
