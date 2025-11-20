package port

import (
	"SubTrack/domain/entity"
	"SubTrack/domain/filter"
	"context"
)

type SubscriptionRepo interface {
	CreateSubscription(ctx context.Context, s *entity.Subscription) error
	GetSubscription(ctx context.Context, id int) (*entity.Subscription, error)
	UpdateSubscription(ctx context.Context, s *entity.Subscription) error
	DeleteSubscription(ctx context.Context, id int) error
	GetSubscriptionList(ctx context.Context, filter filter.ListFilter) ([]entity.Subscription, error)
	GetTotalSum(ctx context.Context, filter filter.SumFilter) (int, error)
}
