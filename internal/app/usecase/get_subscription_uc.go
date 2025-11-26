package usecase

import (
	"context"
	"database/sql"
	"errors"

	"github.com/maket12/SubTrack/internal/app/dto"
	"github.com/maket12/SubTrack/internal/app/mappers"
	"github.com/maket12/SubTrack/internal/app/uc_errors"
	"github.com/maket12/SubTrack/internal/domain/port"
)

type GetSubscriptionUC struct {
	Subscriptions port.SubscriptionRepository
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
	sub, err := uc.Subscriptions.Get(ctx, in.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
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
