package usecase

import (
	"context"
	"database/sql"
	"errors"

	"github.com/maket12/SubTrack/internal/app/dto"
	"github.com/maket12/SubTrack/internal/app/uc_errors"
	"github.com/maket12/SubTrack/internal/domain/port"
)

type DeleteSubscriptionUC struct {
	Subscriptions port.SubscriptionRepository
}

func (uc *DeleteSubscriptionUC) Execute(ctx context.Context, in dto.DeleteSubscription) (dto.DeleteSubscriptionResponse, error) {
	/* ####################
	   #	Validation    #
	   ####################
	*/
	if in.ID == 0 {
		return dto.DeleteSubscriptionResponse{Deleted: false}, uc_errors.ErrInvalidSubscriptionID
	}

	/* ####################
	   #	 Request      #
	   ####################
	*/
	err := uc.Subscriptions.Delete(ctx, in.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return dto.DeleteSubscriptionResponse{Deleted: false}, uc_errors.ErrSubscriptionNotFound
		}
		return dto.DeleteSubscriptionResponse{Deleted: false}, uc_errors.ErrDeleteSubscription
	}

	return dto.DeleteSubscriptionResponse{Deleted: true}, nil
}
