package usecase

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/maket12/SubTrack/internal/app/dto"
	"github.com/maket12/SubTrack/internal/app/uc_errors"
	"github.com/maket12/SubTrack/internal/domain/port/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"testing"
)

type CreateCase struct {
	Name       string
	Input      dto.CreateSubscription
	RepoOutput int
	Output     dto.CreateSubscriptionResponse
	WantErr    error
	RepoErr    error
}

var CreateCases = []CreateCase{
	{
		Name: "empty service name",
		Input: dto.CreateSubscription{
			ServiceName: "",
			Price:       100,
			UserID:      uuid.New().String(),
			StartDate:   "01-01-2025",
			EndDate:     nil,
		},
		WantErr: uc_errors.ErrEmptyServiceName,
	},

	{
		Name: "empty user id",
		Input: dto.CreateSubscription{
			ServiceName: "AmazonTV",
			Price:       1000,
			UserID:      "",
			StartDate:   "02-04-2025",
			EndDate:     nil,
		},
		WantErr: uc_errors.ErrEmptyUserID,
	},

	{
		Name: "invalid user id",
		Input: dto.CreateSubscription{
			ServiceName: "Spotify",
			Price:       200,
			UserID:      "not-a-uuid",
			StartDate:   "01-01-2025",
			EndDate:     nil,
		},
		WantErr: uc_errors.ErrInvalidUserID,
	},

	{
		Name: "invalid start date",
		Input: dto.CreateSubscription{
			ServiceName: "Whoosh",
			Price:       350,
			UserID:      uuid.New().String(),
			StartDate:   "04/04/2025",
			EndDate:     nil,
		},
		WantErr: uc_errors.ErrInvalidDate,
	},

	{
		Name: "invalid end date",
		Input: dto.CreateSubscription{
			ServiceName: "SberPrime",
			Price:       500,
			UserID:      uuid.New().String(),
			StartDate:   "04-04-2025",
			EndDate:     vPtr("04/05/2025"),
		},
		WantErr: uc_errors.ErrInvalidDate,
	},

	{
		Name: "repository error",
		Input: dto.CreateSubscription{
			ServiceName: "Alfa+",
			Price:       1500,
			UserID:      uuid.New().String(),
			StartDate:   "27-11-2025",
			EndDate:     vPtr("27-11-2026"),
		},
		WantErr: uc_errors.ErrCreateSubscription,
		RepoErr: errors.New("db error"),
	},

	{
		Name: "success create",
		Input: dto.CreateSubscription{
			ServiceName: "Netflix",
			Price:       2000,
			UserID:      uuid.New().String(),
			StartDate:   "01-01-2025",
			EndDate:     vPtr("01-02-2025"),
		},
		RepoOutput: 1,
		Output:     dto.CreateSubscriptionResponse{ID: 1},
		WantErr:    nil,
		RepoErr:    nil,
	},
}

func vPtr[T any](v T) *T {
	return &v
}

func TestCreateSubscriptionUC(t *testing.T) {
	for _, tt := range CreateCases {
		t.Run(tt.Name, func(t *testing.T) {
			repo := new(mocks.SubscriptionRepository)
			uc := &CreateSubscriptionUC{Subscriptions: repo}

			shouldCallRepo :=
				tt.WantErr == nil ||
					errors.Is(tt.WantErr, uc_errors.ErrCreateSubscription)

			if shouldCallRepo {
				repo.On("Create", mock.Anything, mock.Anything).
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
