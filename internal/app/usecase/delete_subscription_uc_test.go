package usecase

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/maket12/SubTrack/internal/app/dto"
	"github.com/maket12/SubTrack/internal/app/uc_errors"
	"github.com/maket12/SubTrack/internal/domain/port/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type DeleteCase struct {
	Name    string
	Input   dto.DeleteSubscription
	Output  dto.DeleteSubscriptionResponse
	WantErr error
	RepoErr error
}

var DeleteCases = []DeleteCase{
	{
		Name:    "invalid sub id",
		Input:   dto.DeleteSubscription{ID: 0},
		WantErr: uc_errors.ErrInvalidSubscriptionID,
	},

	{
		Name:    "not found",
		Input:   dto.DeleteSubscription{ID: 1},
		WantErr: uc_errors.ErrSubscriptionNotFound,
		RepoErr: sql.ErrNoRows,
	},

	{
		Name:    "repository error",
		Input:   dto.DeleteSubscription{ID: 1},
		WantErr: uc_errors.ErrDeleteSubscription,
		RepoErr: errors.New("db error"),
	},

	{
		Name:    "success delete",
		Input:   dto.DeleteSubscription{ID: 1},
		Output:  dto.DeleteSubscriptionResponse{Deleted: true},
		WantErr: nil,
		RepoErr: nil,
	},
}

func TestDeleteSubscriptionUC(t *testing.T) {
	for _, tt := range DeleteCases {
		t.Run(tt.Name, func(t *testing.T) {
			repo := new(mocks.SubscriptionRepository)
			uc := &DeleteSubscriptionUC{Subscriptions: repo}

			shouldCallRepo :=
				tt.WantErr == nil ||
					errors.Is(tt.WantErr, uc_errors.ErrDeleteSubscription) ||
					errors.Is(tt.WantErr, uc_errors.ErrSubscriptionNotFound)

			if shouldCallRepo {
				repo.On("Delete", mock.Anything, mock.Anything).
					Return(tt.RepoErr)
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
