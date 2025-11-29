package usecase

import (
	"context"
	"time"

	"github.com/maket12/SubTrack/internal/app/dto"
	"github.com/maket12/SubTrack/internal/app/uc_errors"
	"github.com/maket12/SubTrack/internal/domain/filter"
	"github.com/maket12/SubTrack/internal/domain/port"

	"github.com/google/uuid"
)

type GetTotalSumUC struct {
	Subscriptions port.SubscriptionRepository
}

func (uc *GetTotalSumUC) Execute(ctx context.Context, in dto.GetTotalSum) (dto.GetTotalSumResponse, error) {
	/* ####################
	   #	 Parsing      #
	   ####################
	*/
	var uidPtr *uuid.UUID
	if in.UserID != nil {
		uid, err := uuid.Parse(*in.UserID)
		if err != nil {
			return dto.GetTotalSumResponse{}, uc_errors.ErrInvalidUserID
		}
		if uid == uuid.Nil {
			return dto.GetTotalSumResponse{}, uc_errors.ErrInvalidUserID
		}
		uidPtr = &uid
	}

	var serviceNamePtr *string
	if in.ServiceName != nil && *in.ServiceName != "" {
		s := *in.ServiceName
		serviceNamePtr = &s
	}

	var startPtr *time.Time
	if in.StartDate != nil {
		start, err := time.Parse("02-01-2006", *in.StartDate)
		if err != nil {
			return dto.GetTotalSumResponse{}, uc_errors.ErrInvalidDate
		}
		startPtr = &start
	}

	var endPtr *time.Time
	if in.EndDate != nil && *in.EndDate != "" {
		t, err := time.Parse("02-01-2006", *in.EndDate)
		if err != nil {
			return dto.GetTotalSumResponse{}, uc_errors.ErrInvalidDate
		}
		endPtr = &t
	}

	/* ####################
	   #	 Request      #
	   ####################
	*/
	f := filter.SumFilter{
		UserID:      uidPtr,
		ServiceName: serviceNamePtr,
		StartDate:   startPtr,
		EndDate:     endPtr,
	}

	sum, err := uc.Subscriptions.GetTotalSum(ctx, f)
	if err != nil {
		return dto.GetTotalSumResponse{}, uc_errors.Wrap(uc_errors.ErrGetTotalSum, err)
	}

	return dto.GetTotalSumResponse{TotalSum: sum}, nil
}
