package usecase

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/maket12/SubTrack/internal/app/dto"
	"github.com/maket12/SubTrack/internal/app/uc_errors"
	"github.com/maket12/SubTrack/internal/domain/entity"
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
	if in.ID == 0 {
		return dto.UpdateSubscriptionResponse{Updated: false}, uc_errors.ErrInvalidSubscriptionID
	}
	if in.ServiceName == "" {
		return dto.UpdateSubscriptionResponse{Updated: false}, uc_errors.ErrEmptyServiceName
	}
	if in.UserID == "" {
		return dto.UpdateSubscriptionResponse{Updated: false}, uc_errors.ErrEmptyUserID
	}

	/* ####################
	   #	 Parsing      #
	   ####################
	*/
	uid, err := uuid.Parse(in.UserID)
	if err != nil {
		return dto.UpdateSubscriptionResponse{Updated: false}, uc_errors.ErrInvalidUserID
	}

	start, err := time.Parse("02-01-2006", in.StartDate)
	if err != nil {
		return dto.UpdateSubscriptionResponse{Updated: false}, uc_errors.ErrInvalidDate
	}

	var end *time.Time
	if in.EndDate != nil {
		t, err := time.Parse("02-01-2006", *in.EndDate)
		if err != nil {
			return dto.UpdateSubscriptionResponse{Updated: false}, uc_errors.ErrInvalidDate
		}
		end = &t
	}

	/* ####################
	   #	 Request      #
	   ####################
	*/
	sub := &entity.Subscription{
		ID:          in.ID,
		ServiceName: in.ServiceName,
		Price:       in.Price,
		UserID:      uid,
		StartDate:   start,
		EndDate:     end,
	}

	err = uc.Subscriptions.Update(ctx, sub)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return dto.UpdateSubscriptionResponse{Updated: false}, uc_errors.ErrSubscriptionNotFound
		}
		return dto.UpdateSubscriptionResponse{Updated: false}, uc_errors.ErrUpdateSubscription
	}

	return dto.UpdateSubscriptionResponse{Updated: true}, nil
}
