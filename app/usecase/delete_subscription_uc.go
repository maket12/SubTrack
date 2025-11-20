package usecase

import (
	"SubTrack/app/dto"
	"SubTrack/app/uc_errors"
	"SubTrack/domain/port"
	"context"
	"errors"
)

type DeleteSubscriptionUC struct {
	Subscriptions port.SubscriptionRepo
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
	err := uc.Subscriptions.DeleteSubscription(ctx, in.ID)
	if err != nil {
		if errors.Is(err, port.ErrSubscriptionNotFound) {
			return dto.DeleteSubscriptionResponse{Deleted: false}, uc_errors.ErrSubscriptionNotFound
		}
		return dto.DeleteSubscriptionResponse{Deleted: false}, uc_errors.ErrDeleteSubscription
	}

	return dto.DeleteSubscriptionResponse{Deleted: true}, nil
}
