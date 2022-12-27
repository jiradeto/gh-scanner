package repositoryusecase_test

import (
	"context"
	"testing"
	"time"

	"github.com/pkg/errors"

	"github.com/AlekSi/pointer"
	"github.com/golang/mock/gomock"
	"github.com/jiradeto/gh-scanner/app/constants"
	"github.com/jiradeto/gh-scanner/app/entities"
	repositoryrepo "github.com/jiradeto/gh-scanner/app/infrastructure/repos/repository"
	repositoryrepomocks "github.com/jiradeto/gh-scanner/app/infrastructure/repos/repository/mocks"
	"github.com/jiradeto/gh-scanner/app/testutils"
	repositoryusecase "github.com/jiradeto/gh-scanner/app/usecases/repository"
	"github.com/jiradeto/gh-scanner/app/utils/gerrors"
	"github.com/stretchr/testify/assert"
)

func TestUpdateOneScanResultInputValidate(t *testing.T) {
	type fields struct {
		ID         *string
		Status     *string
		Findings   []entities.ScanFinding
		QueuedAt   *time.Time
		ScanningAt *time.Time
		FinishedAt *time.Time
	}
	now := time.Now()
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
			name: "ID is empty",
			fields: fields{
				ID: pointer.ToString(""),
			},
			wantErr: true,
		},
		{
			name: "scan result status is empty",
			fields: fields{
				Status: pointer.ToString(""),
			},
			wantErr: true,
		},
		{
			name: "invalid scan result status",
			fields: fields{
				Status: pointer.ToString("foo_status"),
			},
			wantErr: true,
		},
		{
			name: "invalid uuid format",
			fields: fields{
				ID:     pointer.ToString("foo_id"),
				Status: pointer.ToString("foo_status"),
			},
			wantErr: true,
		},
		{
			name: "happy case",
			fields: fields{
				ID:         pointer.ToString("7827f71a-df70-4c16-964f-d65836ec4312"),
				Status:     pointer.ToString(entities.ScanResultStatusInProgress.String()),
				QueuedAt:   &now,
				ScanningAt: &now,
				FinishedAt: &now,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &repositoryusecase.UpdateOneScanResultInput{
				ID:         tt.fields.ID,
				Status:     tt.fields.Status,
				Findings:   tt.fields.Findings,
				QueuedAt:   tt.fields.QueuedAt,
				ScanningAt: tt.fields.ScanningAt,
				FinishedAt: tt.fields.FinishedAt,
			}
			if err := f.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_useCase_UpdateOneScanResult(t *testing.T) {
	type args struct {
		ctx   context.Context
		input repositoryusecase.UpdateOneScanResultInput
	}
	now := time.Now()
	const errLocation = "repositoryUsecase/UpdateOneScanResult: %s"
	mockUUID := "7827f71a-df70-4c16-964f-d65836ec4312"
	scanResult := entities.ScanResult{
		ID:     pointer.ToString(mockUUID),
		Status: entities.ScanResultStatusInProgress,
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
					mock.EXPECT().UpdateOneScanResult(gomock.Any(), repositoryrepo.UpdateOneScanResultInput{
						ScanResultEntity: &scanResult,
					}).Return(nil)
					return mock
				},
			},
			args: args{
				ctx: context.Background(),
				input: repositoryusecase.UpdateOneScanResultInput{
					ID:         pointer.ToString(mockUUID),
					Status:     pointer.ToString(entities.ScanResultStatusInProgress.String()),
					QueuedAt:   &now,
					ScanningAt: &now,
					FinishedAt: &now,
				},
			},
			want:    &scanResult,
			wantErr: nil,
		},
		{
			name: "failed to find one scan result",
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
				input: repositoryusecase.UpdateOneScanResultInput{
					ID:     pointer.ToString(mockUUID),
					Status: pointer.ToString(entities.ScanResultStatusInProgress.String()),
				},
			},
			wantErr: gerrors.InternalError{
				Code:    constants.StatusCodeDatabaseError,
				Message: constants.ErrorMessageDatabaseError,
			}.Wrap(errors.Wrapf(errors.New("mock error"), errLocation, "unable to find scan result")),
		},
		{
			name: "when record not found",
			fields: fields{
				CreateRepositoryRepo: func(ctrl *gomock.Controller) repositoryrepo.Repo {
					mock := repositoryrepomocks.NewMocks(ctrl)
					mock.EXPECT().FindOneScanResult(gomock.Any(), repositoryrepo.FindOneScanResultInput{
						ID: pointer.ToString(mockUUID),
					}).Return(nil, nil)
					return mock
				},
			},
			args: args{
				ctx: context.Background(),
				input: repositoryusecase.UpdateOneScanResultInput{
					ID:     pointer.ToString(mockUUID),
					Status: pointer.ToString(entities.ScanResultStatusInProgress.String()),
				},
			},
			wantErr: gerrors.RecordNotFoundError{
				Code:    constants.StatusCodeEntryNotFound,
				Message: constants.ErrorMessageNotFound,
			}.Wrap(errors.Errorf(errLocation, "not found scan result")),
		},
		{
			name: "failed to update scan result",
			fields: fields{
				CreateRepositoryRepo: func(ctrl *gomock.Controller) repositoryrepo.Repo {
					mock := repositoryrepomocks.NewMocks(ctrl)
					mock.EXPECT().FindOneScanResult(gomock.Any(), repositoryrepo.FindOneScanResultInput{
						ID: pointer.ToString(mockUUID),
					}).Return(&scanResult, nil)
					mock.EXPECT().UpdateOneScanResult(gomock.Any(), repositoryrepo.UpdateOneScanResultInput{
						ScanResultEntity: &scanResult,
					}).Return(errors.New("mock error"))
					return mock
				},
			},
			args: args{
				ctx: context.Background(),
				input: repositoryusecase.UpdateOneScanResultInput{
					ID:     pointer.ToString(mockUUID),
					Status: pointer.ToString(entities.ScanResultStatusInProgress.String()),
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

			got, err := h.uc.UpdateOneScanResult(tt.args.ctx, tt.args.input)

			if !testutils.AssertError(t, "useCase.UpdateOneScanResult()", err, tt.wantErr) {
				return
			}
			assert.Equal(t, tt.want, got, "useCase.UpdateOneScanResult()")
		})
	}
}
