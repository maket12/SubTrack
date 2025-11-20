package usecase

import (
	"SubTrack/app/dto"
	"SubTrack/app/mappers"
	"SubTrack/app/uc_errors"
	"SubTrack/domain/filter"
	"SubTrack/domain/port"
	"context"

	"github.com/google/uuid"
)

type GetSubscriptionListUC struct {
	Subscriptions port.SubscriptionRepo
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
	subs, err := uc.Subscriptions.GetSubscriptionList(ctx, f)
	if err != nil {
		return dto.GetSubscriptionListResponse{}, uc_errors.ErrGetSubscriptionList
	}

	/* ####################
	   #	 Mapping      #
	   ####################
	*/
	return mappers.MapIntoGetSubscriptionListDTO(subs), nil
}
