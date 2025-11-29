package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/maket12/SubTrack/internal/app/dto"
	"github.com/maket12/SubTrack/internal/app/uc_errors"
	"github.com/maket12/SubTrack/internal/domain/entity"
	"github.com/maket12/SubTrack/internal/domain/port/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type GetListCase struct {
	Name       string
	Input      dto.GetSubscriptionList
	RepoOutput []entity.Subscription
	Output     dto.GetSubscriptionListResponse
	WantErr    error
	RepoErr    error
}

var GetListCases = []GetListCase{
	{
		Name: "negative limit",
		Input: dto.GetSubscriptionList{
			UserID:      vPtr(uuid.New().String()),
			ServiceName: vPtr("Yandex"),
			Limit:       -1,
			Offset:      0,
		},
		WantErr: uc_errors.ErrInvalidLimit,
	},

	{
		Name: "negative offset",
		Input: dto.GetSubscriptionList{
			UserID:      vPtr(uuid.New().String()),
			ServiceName: vPtr("Yandex"),
			Limit:       10,
			Offset:      -1,
		},
		WantErr: uc_errors.ErrInvalidOffset,
	},

	{
		Name: "invalid user id",
		Input: dto.GetSubscriptionList{
			UserID:      vPtr("not-a-uuid"),
			ServiceName: vPtr("YandexMusic"),
			Limit:       0,
			Offset:      0,
		},
		WantErr: uc_errors.ErrInvalidUserID,
	},

	{
		Name: "invalid user id",
		Input: dto.GetSubscriptionList{
			UserID:      vPtr(uuid.Nil.String()),
			ServiceName: vPtr("YandexMusic"),
			Limit:       0,
			Offset:      0,
		},
		WantErr: uc_errors.ErrInvalidUserID,
	},

	{
		Name: "repository error",
		Input: dto.GetSubscriptionList{
			UserID:      vPtr(uuid.New().String()),
			ServiceName: vPtr("YandexMusic"),
			Limit:       0,
			Offset:      0,
		},
		WantErr: uc_errors.ErrGetSubscriptionList,
		RepoErr: errors.New("db error"),
	},

	{
		Name: "success get list",
		Input: dto.GetSubscriptionList{
			UserID:      vPtr(uuid.New().String()),
			ServiceName: vPtr("YandexMusic"),
			Limit:       2,
			Offset:      0,
		},
		RepoOutput: []entity.Subscription{
			{
				ID:          1,
				ServiceName: "YandexMusic",
				Price:       450,
				UserID:      uuid.MustParse("79acd47c-cacd-40d7-876b-2f131bdf3014"),
				StartDate:   parseTime("29-11-2025"),
				EndDate:     vPtr(parseTime("29-11-2026")),
			},
			{
				ID:          2,
				ServiceName: "YandexMusic",
				Price:       200,
				UserID:      uuid.MustParse("79acd52c-cacd-40d7-876b-2f131bdf3014"),
				StartDate:   parseTime("29-11-2025"),
				EndDate:     vPtr(parseTime("29-12-2025")),
			},
		},
		Output: dto.GetSubscriptionListResponse{
			Items: []dto.GetSubscriptionResponse{
				{
					ID:          1,
					ServiceName: "YandexMusic",
					Price:       450,
					UserID:      "79acd47c-cacd-40d7-876b-2f131bdf3014",
					StartDate:   "29-11-2025",
					EndDate:     vPtr("29-11-2026"),
				},
				{
					ID:          2,
					ServiceName: "YandexMusic",
					Price:       200,
					UserID:      "79acd52c-cacd-40d7-876b-2f131bdf3014",
					StartDate:   "29-11-2025",
					EndDate:     vPtr("29-12-2025"),
				},
			},
		},
		WantErr: nil,
		RepoErr: nil,
	},
}

func TestGetSubscriptionListUC(t *testing.T) {
	for _, tt := range GetListCases {
		t.Run(tt.Name, func(t *testing.T) {
			repo := new(mocks.SubscriptionRepository)
			uc := &GetSubscriptionListUC{Subscriptions: repo}

			shouldCallRepo :=
				tt.WantErr == nil ||
					errors.Is(tt.WantErr, uc_errors.ErrGetSubscriptionList)

			if shouldCallRepo {
				repo.On("GetList", mock.Anything, mock.Anything).
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
