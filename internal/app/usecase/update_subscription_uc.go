package usecase

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/maket12/SubTrack/internal/app/dto"
	"github.com/maket12/SubTrack/internal/app/uc_errors"
	"github.com/maket12/SubTrack/internal/domain/port"

	"github.com/google/uuid"
)

type UpdateSubscriptionUC struct {
	Subscriptions port.SubscriptionRepository
}

func (uc *UpdateSubscriptionUC) Execute(ctx context.Context, in dto.UpdateSubscription) (dto.UpdateSubscriptionResponse, error) {
	/* ####################
	   #	Validation    #
	   ####################
	*/
	if in.ID <= 0 {
		return dto.UpdateSubscriptionResponse{Updated: false}, uc_errors.ErrInvalidSubscriptionID
	}

	if in.ServiceName == nil &&
		in.Price == nil &&
		in.UserID == nil &&
		in.StartDate == nil &&
		in.EndDate == nil {
		return dto.UpdateSubscriptionResponse{Updated: false}, nil
	}

	/* ####################
	   #   Load current   #
	   ####################
	*/
	sub, err := uc.Subscriptions.Get(ctx, in.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return dto.UpdateSubscriptionResponse{Updated: false},
				uc_errors.Wrap(uc_errors.ErrSubscriptionNotFound, err)
		}
		return dto.UpdateSubscriptionResponse{Updated: false},
			uc_errors.Wrap(uc_errors.ErrGetSubscription, err)
	}

	/* ####################
	   #	 Parsing      #
	   ####################
	*/
	if in.ServiceName != nil {
		if *in.ServiceName == "" {
			return dto.UpdateSubscriptionResponse{Updated: false}, uc_errors.ErrEmptyServiceName
		}
		sub.ServiceName = *in.ServiceName
	}

	if in.Price != nil {
		sub.Price = *in.Price
	}

	if in.UserID != nil {
		if *in.UserID == "" {
			return dto.UpdateSubscriptionResponse{Updated: false}, uc_errors.ErrEmptyUserID
		}
		uid, err := uuid.Parse(*in.UserID)
		if err != nil || uid == uuid.Nil {
			return dto.UpdateSubscriptionResponse{Updated: false}, uc_errors.ErrInvalidUserID
		}
		sub.UserID = uid
	}

	if in.StartDate != nil {
		t, err := time.Parse("02-01-2006", *in.StartDate)
		if err != nil {
			return dto.UpdateSubscriptionResponse{Updated: false}, uc_errors.ErrInvalidDate
		}
		sub.StartDate = t
	}

	if in.EndDate != nil {
		if *in.EndDate == "" {
			// empty string end-date as "clear end_date"
			sub.EndDate = nil
		} else {
			t, err := time.Parse("02-01-2006", *in.EndDate)
			if err != nil {
				return dto.UpdateSubscriptionResponse{Updated: false}, uc_errors.ErrInvalidDate
			}
			sub.EndDate = &t
		}
	}

	/* ####################
	   #	 Request      #
	   ####################
	*/
	if err := uc.Subscriptions.Update(ctx, sub); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return dto.UpdateSubscriptionResponse{Updated: false},
				uc_errors.Wrap(uc_errors.ErrSubscriptionNotFound, err)
		}
		return dto.UpdateSubscriptionResponse{Updated: false},
			uc_errors.Wrap(uc_errors.ErrUpdateSubscription, err)
	}

	return dto.UpdateSubscriptionResponse{Updated: true}, nil
}
