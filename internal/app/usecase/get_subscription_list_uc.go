package usecase

import (
	"context"

	"github.com/maket12/SubTrack/internal/app/dto"
	"github.com/maket12/SubTrack/internal/app/mappers"
	"github.com/maket12/SubTrack/internal/app/uc_errors"
	"github.com/maket12/SubTrack/internal/domain/filter"
	"github.com/maket12/SubTrack/internal/domain/port"

	"github.com/google/uuid"
)

type GetSubscriptionListUC struct {
	Subscriptions port.SubscriptionRepository
}

func (uc *GetSubscriptionListUC) Execute(ctx context.Context, in dto.GetSubscriptionList) (dto.GetSubscriptionListResponse, error) {
	/* ####################
	   #	Validation    #
	   ####################
	*/
	if in.Limit < 0 {
		return dto.GetSubscriptionListResponse{}, uc_errors.ErrInvalidLimit
	}
	if in.Offset < 0 {
		return dto.GetSubscriptionListResponse{}, uc_errors.ErrInvalidOffset
	}

	/* ####################
	   #	 Parsing      #
	   ####################
	*/
	var f filter.ListFilter

	if in.UserID != nil {
		uid, err := uuid.Parse(*in.UserID)
		if err != nil {
			return dto.GetSubscriptionListResponse{}, uc_errors.ErrInvalidUserID
		}
		f.UserID = &uid
	}

	if in.ServiceName != nil {
		f.ServiceName = in.ServiceName
	}

	f.Limit = in.Limit
	f.Offset = in.Offset

	/* ####################
	   #	 Request      #
	   ####################
	*/
	subs, err := uc.Subscriptions.GetList(ctx, f)
	if err != nil {
		return dto.GetSubscriptionListResponse{}, uc_errors.ErrGetSubscriptionList
	}

	/* ####################
	   #	 Mapping      #
	   ####################
	*/
	return mappers.MapIntoGetSubscriptionListDTO(subs), nil
}
