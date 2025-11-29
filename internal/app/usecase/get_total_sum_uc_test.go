package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/maket12/SubTrack/internal/app/dto"
	"github.com/maket12/SubTrack/internal/app/uc_errors"
	"github.com/maket12/SubTrack/internal/domain/port/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type GetTotalSumCase struct {
	Name       string
	Input      dto.GetTotalSum
	RepoOutput int
	Output     dto.GetTotalSumResponse
	WantErr    error
	RepoErr    error
}

var GetTotalSumCases = []GetTotalSumCase{
	{
		Name: "invalid user id",
		Input: dto.GetTotalSum{
			UserID:      vPtr("not-a-uuid"),
			ServiceName: vPtr("YandexMusic"),
			StartDate:   vPtr("01-01-2026"),
			EndDate:     vPtr("01-02-2026"),
		},
		WantErr: uc_errors.ErrInvalidUserID,
	},

	{
		Name: "invalid user id(nil)",
		Input: dto.GetTotalSum{
			UserID:      vPtr(uuid.Nil.String()),
			ServiceName: vPtr("YandexMusic"),
			StartDate:   vPtr("01-01-2026"),
			EndDate:     vPtr("01-02-2026"),
		},
		WantErr: uc_errors.ErrInvalidUserID,
	},

	{
		Name: "invalid start date",
		Input: dto.GetTotalSum{
			UserID:      vPtr(uuid.New().String()),
			ServiceName: vPtr("ChatGPT Plus"),
			StartDate:   vPtr("04/04/2025"),
			EndDate:     nil,
		},
		WantErr: uc_errors.ErrInvalidDate,
	},

	{
		Name: "invalid end date",
		Input: dto.GetTotalSum{
			UserID:      vPtr(uuid.New().String()),
			ServiceName: vPtr("Rambler"),
			StartDate:   vPtr("04-04-2025"),
			EndDate:     vPtr("04/05/2025"),
		},
		WantErr: uc_errors.ErrInvalidDate,
	},

	{
		Name: "repository error",
		Input: dto.GetTotalSum{
			UserID:      vPtr(uuid.New().String()),
			ServiceName: vPtr("ChatGPT Pro"),
			StartDate:   vPtr("04-04-2025"),
			EndDate:     vPtr("04-05-2025"),
		},
		WantErr: uc_errors.ErrGetTotalSum,
		RepoErr: errors.New("db error"),
	},

	{
		Name: "success get total sum",
		Input: dto.GetTotalSum{
			UserID:      vPtr(uuid.New().String()),
			ServiceName: vPtr("AxiCorp"),
			StartDate:   vPtr("04-04-2025"),
			EndDate:     vPtr("04-05-2025"),
		},
		RepoOutput: 10250,
		Output:     dto.GetTotalSumResponse{TotalSum: 10250},
		WantErr:    nil,
		RepoErr:    nil,
	},
}

func TestGetTotalSumUC(t *testing.T) {
	for _, tt := range GetTotalSumCases {
		t.Run(tt.Name, func(t *testing.T) {
			repo := new(mocks.SubscriptionRepository)
			uc := &GetTotalSumUC{Subscriptions: repo}

			shouldCallRepo :=
				tt.WantErr == nil ||
					errors.Is(tt.WantErr, uc_errors.ErrGetTotalSum)

			if shouldCallRepo {
				repo.On("GetTotalSum", mock.Anything, mock.Anything).
					Return(tt.RepoOutput, tt.RepoErr)
			}

			resp, err := uc.Execute(context.Background(), tt.Input)

			if tt.WantErr != nil {
				assert.Error(t, err)
				assert.True(t, errors.Is(err, tt.WantErr),
					"expected error '%v' but got '%v'", tt.WantErr, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.Output, resp)
			}

			repo.AssertExpectations(t)
		})
	}
}
