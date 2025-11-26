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
	   #	Validation    #
	   ####################
	*/
	if in.UserID == "" {
		return dto.GetTotalSumResponse{}, uc_errors.ErrEmptyUserID
	}
	if in.StartDate == "" {
		return dto.GetTotalSumResponse{}, uc_errors.ErrEmptyDate
	}

	/* ####################
	   #	 Parsing      #
	   ####################
	*/
	uid, err := uuid.Parse(in.UserID)
	if err != nil {
		return dto.GetTotalSumResponse{}, uc_errors.ErrInvalidUserID
	}

	var sName string
	if in.ServiceName != nil {
		sName = *in.ServiceName
	}

	start, err := time.Parse("02-01-2006", in.StartDate)
	if err != nil {
		return dto.GetTotalSumResponse{}, uc_errors.ErrInvalidDate
	}

	var end *time.Time
	if in.EndDate != nil {
		t, err := time.Parse("02-01-2006", *in.EndDate)
		if err != nil {
			return dto.GetTotalSumResponse{}, uc_errors.ErrInvalidDate
		}
		end = &t
	}

	/* ####################
	   #	 Request      #
	   ####################
	*/
	f := filter.SumFilter{
		UserID:      &uid,
		ServiceName: &sName,
		StartDate:   &start,
		EndDate:     end,
	}

	sum, err := uc.Subscriptions.GetTotalSum(ctx, f)
	if err != nil {
		return dto.GetTotalSumResponse{}, uc_errors.ErrGetTotalSum
	}

	return dto.GetTotalSumResponse{TotalSum: sum}, nil
}
