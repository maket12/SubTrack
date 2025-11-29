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

type UpdateSubscriptionCase struct {
	Name          string
	Input         dto.UpdateSubscription
	Output        dto.UpdateSubscriptionResponse
	GetRepoOutput *entity.Subscription
	WantErr       error
	GetRepoErr    error
	UpdateRepoErr error
}

var UpdateSubscriptionCases = []UpdateSubscriptionCase{
	{
		Name: "invalid sub id",
		Input: dto.UpdateSubscription{
			ID:          0,
			ServiceName: vPtr("YandexMusic"),
			Price:       nil,
			UserID:      nil,
			StartDate:   nil,
			EndDate:     nil,
		},
		WantErr: uc_errors.ErrInvalidSubscriptionID,
	},

	{
		Name: "nothing to update",
		Input: dto.UpdateSubscription{
			ID:          1,
			ServiceName: nil,
			Price:       nil,
			UserID:      nil,
			StartDate:   nil,
			EndDate:     nil,
		},
		Output:     dto.UpdateSubscriptionResponse{Updated: false},
		WantErr:    nil,
		GetRepoErr: nil,
	},

	{
		Name: "not found(get)",
		Input: dto.UpdateSubscription{
			ID:          1,
			ServiceName: vPtr("YandexMusic"),
			Price:       vPtr(200),
			UserID:      nil,
			StartDate:   nil,
			EndDate:     nil,
		},
		WantErr:    uc_errors.ErrSubscriptionNotFound,
		GetRepoErr: sql.ErrNoRows,
	},

	{
		Name: "repository error(get)",
		Input: dto.UpdateSubscription{
			ID:          1,
			ServiceName: vPtr("YandexMusic"),
			Price:       nil,
			UserID:      nil,
			StartDate:   nil,
			EndDate:     nil,
		},
		WantErr:    uc_errors.ErrGetSubscription,
		GetRepoErr: errors.New("db error"),
	},

	{
		Name: "empty service name",
		Input: dto.UpdateSubscription{
			ID:          1,
			ServiceName: vPtr(""),
			Price:       nil,
			UserID:      nil,
			StartDate:   nil,
			EndDate:     nil,
		},
		WantErr: uc_errors.ErrEmptyServiceName,
	},

	{
		Name: "empty user id",
		Input: dto.UpdateSubscription{
			ID:          1,
			ServiceName: vPtr("YandexMusic"),
			Price:       vPtr(1000),
			UserID:      vPtr(""),
			StartDate:   nil,
			EndDate:     nil,
		},
		WantErr: uc_errors.ErrEmptyUserID,
	},

	{
		Name: "invalid user id",
		Input: dto.UpdateSubscription{
			ID:          1,
			ServiceName: vPtr("Yandex"),
			Price:       vPtr(1000),
			UserID:      vPtr("not-a-uuid"),
			StartDate:   nil,
			EndDate:     nil,
		},
		WantErr: uc_errors.ErrInvalidUserID,
	},

	{
		Name: "invalid user id(nil)",
		Input: dto.UpdateSubscription{
			ID:          1,
			ServiceName: vPtr("Yandex"),
			Price:       vPtr(1000),
			UserID:      vPtr(uuid.Nil.String()),
			StartDate:   nil,
			EndDate:     nil,
		},
		WantErr: uc_errors.ErrInvalidUserID,
	},

	{
		Name: "invalid start date",
		Input: dto.UpdateSubscription{
			ID:          1,
			ServiceName: vPtr("Yandex"),
			Price:       vPtr(1000),
			UserID:      vPtr(uuid.New().String()),
			StartDate:   vPtr("10/10/2024"),
			EndDate:     nil,
		},
		WantErr: uc_errors.ErrInvalidDate,
	},

	{
		Name: "invalid end date",
		Input: dto.UpdateSubscription{
			ID:          1,
			ServiceName: vPtr("YandexMusic"),
			Price:       vPtr(1000),
			UserID:      vPtr(uuid.New().String()),
			StartDate:   vPtr("10-10-2024"),
			EndDate:     vPtr("11/11/2024"),
		},
		WantErr: uc_errors.ErrInvalidDate,
	},

	{
		Name: "not found(update)",
		Input: dto.UpdateSubscription{
			ID:          1,
			ServiceName: vPtr("Spotify"),
			Price:       vPtr(975),
			UserID:      nil,
			StartDate:   nil,
			EndDate:     nil,
		},
		GetRepoOutput: &entity.Subscription{
			ID:          1,
			ServiceName: "YandexMusic",
			Price:       500,
			UserID:      uuid.UUID{},
			StartDate:   time.Time{},
			EndDate:     nil,
		},
		WantErr:       uc_errors.ErrSubscriptionNotFound,
		UpdateRepoErr: sql.ErrNoRows,
	},

	{
		Name: "repository error(update)",
		Input: dto.UpdateSubscription{
			ID:          1,
			ServiceName: vPtr("Spotify"),
			Price:       vPtr(975),
			UserID:      nil,
			StartDate:   nil,
			EndDate:     nil,
		},
		GetRepoOutput: &entity.Subscription{
			ID:          1,
			ServiceName: "YandexMusic",
			Price:       500,
			UserID:      uuid.UUID{},
			StartDate:   time.Time{},
			EndDate:     nil,
		},
		WantErr:       uc_errors.ErrUpdateSubscription,
		UpdateRepoErr: errors.New("db error"),
	},

	{
		Name: "empty end date",
		Input: dto.UpdateSubscription{
			ID:          1,
			ServiceName: vPtr("Spotify"),
			Price:       vPtr(975),
			UserID:      nil,
			StartDate:   nil,
			EndDate:     vPtr(""),
		},
		GetRepoOutput: &entity.Subscription{
			ID:          1,
			ServiceName: "YandexMusic",
			Price:       500,
			UserID:      uuid.UUID{},
			StartDate:   time.Time{},
			EndDate:     vPtr(time.Now()),
		},
		Output:        dto.UpdateSubscriptionResponse{Updated: true},
		WantErr:       nil,
		GetRepoErr:    nil,
		UpdateRepoErr: nil,
	},

	{
		Name: "success update",
		Input: dto.UpdateSubscription{
			ID:          1,
			ServiceName: vPtr("Spotify"),
			Price:       vPtr(975),
			UserID:      nil,
			StartDate:   nil,
			EndDate:     vPtr("30-11-2025"),
		},
		GetRepoOutput: &entity.Subscription{
			ID:          1,
			ServiceName: "YandexMusic",
			Price:       500,
			UserID:      uuid.UUID{},
			StartDate:   time.Time{},
			EndDate:     nil,
		},
		Output:        dto.UpdateSubscriptionResponse{Updated: true},
		WantErr:       nil,
		GetRepoErr:    nil,
		UpdateRepoErr: nil,
	},
}

func allNilInput(in dto.UpdateSubscription) bool {
	return in.ServiceName == nil &&
		in.Price == nil &&
		in.UserID == nil &&
		in.StartDate == nil &&
		in.EndDate == nil
}

func isValidationAfterGet(err error) bool {
	return errors.Is(err, uc_errors.ErrEmptyServiceName) ||
		errors.Is(err, uc_errors.ErrEmptyUserID) ||
		errors.Is(err, uc_errors.ErrInvalidUserID) ||
		errors.Is(err, uc_errors.ErrInvalidDate)
}

func TestUpdateSubscriptionUC(t *testing.T) {
	for _, tt := range UpdateSubscriptionCases {
		t.Run(tt.Name, func(t *testing.T) {
			repo := new(mocks.SubscriptionRepository)
			uc := &UpdateSubscriptionUC{Subscriptions: repo}

			allNil := allNilInput(tt.Input)

			needGet := tt.Input.ID > 0 && !allNil

			if needGet {
				var sub *entity.Subscription
				if tt.GetRepoOutput != nil {
					sub = tt.GetRepoOutput
				} else if tt.GetRepoErr == nil {
					sub = &entity.Subscription{}
				}
				repo.
					On("Get", mock.Anything, tt.Input.ID).
					Return(sub, tt.GetRepoErr)
			}

			needUpdate :=
				needGet &&
					tt.GetRepoErr == nil &&
					!isValidationAfterGet(tt.WantErr)

			if needUpdate {
				repo.
					On("Update", mock.Anything, mock.AnythingOfType("*entity.Subscription")).
					Return(tt.UpdateRepoErr)
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
