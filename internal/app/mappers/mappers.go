package mappers

import (
	"github.com/maket12/SubTrack/internal/app/dto"
	"github.com/maket12/SubTrack/internal/domain/entity"
)

func MapIntoGetSubscriptionDTO(sub *entity.Subscription) dto.GetSubscriptionResponse {
	var end *string
	if sub.EndDate != nil {
		formatted := sub.EndDate.Format("02-01-2006")
		end = &formatted
	}

	return dto.GetSubscriptionResponse{
		ID:          sub.ID,
		ServiceName: sub.ServiceName,
		Price:       sub.Price,
		UserID:      sub.UserID.String(),
		StartDate:   sub.StartDate.Format("02-01-2006"),
		EndDate:     end,
	}
}

func MapIntoGetSubscriptionListDTO(subs []entity.Subscription) dto.GetSubscriptionListResponse {
	items := make([]dto.GetSubscriptionResponse, 0, len(subs))

	for _, s := range subs {
		item := MapIntoGetSubscriptionDTO(&s)
		items = append(items, item)
	}

	return dto.GetSubscriptionListResponse{
		Items: items,
	}
}
