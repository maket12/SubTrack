package usecase

import (
	"context"
	"time"

	"github.com/maket12/SubTrack/internal/app/dto"
	"github.com/maket12/SubTrack/internal/app/uc_errors"
	"github.com/maket12/SubTrack/internal/domain/entity"
	"github.com/maket12/SubTrack/internal/domain/port"

	"github.com/google/uuid"
)

type CreateSubscriptionUC struct {
	Subscriptions port.SubscriptionRepository
}

func (uc *CreateSubscriptionUC) Execute(ctx context.Context, in dto.CreateSubscription) (dto.CreateSubscriptionResponse, error) {
	/* ####################
	   #	Validation    #
	   ####################
	*/
	if in.ServiceName == "" {
		return dto.CreateSubscriptionResponse{}, uc_errors.ErrEmptyServiceName
	}
	if in.UserID == "" {
		return dto.CreateSubscriptionResponse{}, uc_errors.ErrEmptyUserID
	}

	/* ####################
	   #	 Parsing      #
	   ####################
	*/
	uid, err := uuid.Parse(in.UserID)
	if err != nil {
		return dto.CreateSubscriptionResponse{}, uc_errors.ErrInvalidUserID
	}

	start, err := time.Parse("02-01-2006", in.StartDate)
	if err != nil {
		return dto.CreateSubscriptionResponse{}, uc_errors.ErrInvalidDate
	}

	var end *time.Time
	if in.EndDate != nil {
		t, err := time.Parse("02-01-2006", *in.EndDate)
		if err != nil {
			return dto.CreateSubscriptionResponse{}, uc_errors.ErrInvalidDate
		}
		end = &t
	}

	/* ####################
	   #	 Request      #
	   ####################
	*/
	sub := &entity.Subscription{
		ServiceName: in.ServiceName,
		Price:       in.Price,
		UserID:      uid,
		StartDate:   start,
		EndDate:     end,
	}

	if err := uc.Subscriptions.Create(ctx, sub); err != nil {
		return dto.CreateSubscriptionResponse{}, uc_errors.Wrap(uc_errors.ErrCreateSubscription, err)
	}

	return dto.CreateSubscriptionResponse{ID: sub.ID}, nil
}
