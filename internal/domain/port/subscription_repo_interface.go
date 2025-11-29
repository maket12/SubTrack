package port

import (
	"context"

	"github.com/maket12/SubTrack/internal/domain/entity"
	"github.com/maket12/SubTrack/internal/domain/filter"
)

type SubscriptionRepository interface {
	Create(ctx context.Context, s *entity.Subscription) (int, error)
	Get(ctx context.Context, id int) (*entity.Subscription, error)
	Update(ctx context.Context, s *entity.Subscription) error
	Delete(ctx context.Context, id int) error
	GetList(ctx context.Context, filter filter.ListFilter) ([]entity.Subscription, error)
	GetTotalSum(ctx context.Context, filter filter.SumFilter) (int, error)
}
