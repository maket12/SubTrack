package usecase

import (
	"SubTrack/app/dto"
	"SubTrack/app/uc_errors"
	"SubTrack/domain/entity"
	"SubTrack/domain/port"
	"context"
	"time"

	"github.com/google/uuid"
)

type CreateSubscriptionUC struct {
	Subscriptions port.SubscriptionRepo
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

	if err := uc.Subscriptions.CreateSubscription(ctx, sub); err != nil {
		return dto.CreateSubscriptionResponse{}, uc_errors.ErrCreateSubscription
	}

	return dto.CreateSubscriptionResponse{ID: sub.ID}, nil
}
