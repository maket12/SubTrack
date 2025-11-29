package usecase

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/maket12/SubTrack/internal/app/dto"
	"github.com/maket12/SubTrack/internal/app/uc_errors"
	"github.com/maket12/SubTrack/internal/domain/entity"
	"github.com/maket12/SubTrack/internal/domain/port/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type GetCase struct {
	Name       string
	Input      dto.GetSubscription
	RepoOutput *entity.Subscription
	Output     dto.GetSubscriptionResponse
	WantErr    error
	RepoErr    error
}

var GetCases = []GetCase{
	{
		Name:    "invalid sub id",
		Input:   dto.GetSubscription{ID: 0},
		WantErr: uc_errors.ErrInvalidSubscriptionID,
	},

	{
		Name:    "not found",
		Input:   dto.GetSubscription{ID: 1},
		WantErr: uc_errors.ErrSubscriptionNotFound,
		RepoErr: sql.ErrNoRows,
	},

	{
		Name:    "repository error",
		Input:   dto.GetSubscription{ID: 1},
		WantErr: uc_errors.ErrGetSubscription,
		RepoErr: errors.New("db error"),
	},

	{
		Name:  "success get",
		Input: dto.GetSubscription{ID: 1},
		RepoOutput: &entity.Subscription{
			ID:          1,
			ServiceName: "Netflix",
			Price:       1000,
			UserID:      uuid.MustParse("79acd47c-cacd-40d7-876b-2f131bdf3014"),
			StartDate:   parseTime("01-01-2026"),
			EndDate:     vPtr(parseTime("01-03-2026")),
		},
		Output: dto.GetSubscriptionResponse{
			ID:          1,
			ServiceName: "Netflix",
			Price:       1000,
			UserID:      "79acd47c-cacd-40d7-876b-2f131bdf3014",
			StartDate:   "01-01-2026",
			EndDate:     vPtr("01-03-2026"),
		},
		WantErr: nil,
		RepoErr: nil,
	},
}

func parseTime(t string) time.Time {
	res, _ := time.Parse("02-01-2006", t)
	return res
}

func TestGetSubscriptionUC(t *testing.T) {
	for _, tt := range GetCases {
		t.Run(tt.Name, func(t *testing.T) {
			repo := new(mocks.SubscriptionRepository)
			uc := &GetSubscriptionUC{Subscriptions: repo}

			shouldCallRepo :=
				tt.WantErr == nil ||
					errors.Is(tt.WantErr, uc_errors.ErrGetSubscription) ||
					errors.Is(tt.WantErr, uc_errors.ErrSubscriptionNotFound)

			if shouldCallRepo {
				repo.On("Get", mock.Anything, mock.Anything).
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
