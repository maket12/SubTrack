package usecase

import (
	"SubTrack/app/dto"
	"SubTrack/app/mappers"
	"SubTrack/app/uc_errors"
	"SubTrack/domain/port"
	"context"
	"errors"
)

type GetSubscriptionUC struct {
	Subscriptions port.SubscriptionRepo
}

func (uc *GetSubscriptionUC) Execute(ctx context.Context, in dto.GetSubscription) (dto.GetSubscriptionResponse, error) {
	/* ####################
	   #	Validation    #
	   ####################
	*/
	if in.ID == 0 {
		return dto.GetSubscriptionResponse{}, uc_errors.ErrInvalidSubscriptionID
	}

	/* ####################
	   #	 Request      #
	   ####################
	*/
	sub, err := uc.Subscriptions.GetSubscription(ctx, in.ID)
	if err != nil {
		if errors.Is(err, port.ErrSubscriptionNotFound) {
			return dto.GetSubscriptionResponse{}, uc_errors.ErrSubscriptionNotFound
		}
		return dto.GetSubscriptionResponse{}, uc_errors.ErrGetSubscription
	}

	/* ####################
	   #	 Mapping      #
	   ####################
	*/
	return mappers.MapIntoGetSubscriptionDTO(sub), nil
}
