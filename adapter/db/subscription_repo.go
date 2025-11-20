package db

import (
	"SubTrack/domain/entity"
	"SubTrack/domain/filter"
	"SubTrack/domain/port"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"

	"github.com/jmoiron/sqlx"
)

type SubscriptionRepo struct {
	db     *sqlx.DB
	logger *slog.Logger
}

func NewSubscriptionRepo(db *sqlx.DB, log *slog.Logger) *SubscriptionRepo {
	return &SubscriptionRepo{
		db:     db,
		logger: log,
	}
}

func (r *SubscriptionRepo) CreateSubscription(ctx context.Context, s *entity.Subscription) error {
	query := `
		INSERT INTO subscriptions
			(service_name, price, user_id, start_date, end_date)
		VALUES 
		    ($1, $2, $3, $4, $5)
		RETURNING id
	`

	var id int
	err := r.db.QueryRowContext(
		ctx,
		query,
		s.ServiceName,
		s.Price,
		s.UserID,
		s.StartDate,
		s.EndDate,
	).Scan(&id)

	if err != nil {
		r.logger.ErrorContext(ctx,
			"failed to create subscription",
			"error", err,
			"service_name", s.ServiceName,
			"user_id", s.UserID,
		)
		return fmt.Errorf("%w", port.ErrCreateSubscription)
	}

	s.ID = id

	r.logger.InfoContext(ctx,
		"added subscription",
		"id", s.ID,
		"service_name", s.ServiceName,
		"user_id", s.UserID,
	)

	return nil
}

func (r *SubscriptionRepo) GetSubscription(ctx context.Context, id int) (*entity.Subscription, error) {
	query := `
		SELECT id, service_name, price, user_id, start_date, end_date
		FROM subscriptions
		WHERE id = $1
	`

	var sub entity.Subscription

	err := r.db.GetContext(ctx, &sub, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			r.logger.WarnContext(ctx,
				"subscription not found",
				"id", id,
			)
			return nil, port.ErrSubscriptionNotFound
		}

		r.logger.ErrorContext(ctx,
			"failed to get subscription",
			"error", err,
			"id", id,
		)
		return nil, fmt.Errorf("%w: %v", port.ErrGetSubscription, err)
	}

	return &sub, nil
}

func (r *SubscriptionRepo) UpdateSubscription(ctx context.Context, s *entity.Subscription) error {
	query := `
        UPDATE subscriptions
        SET 
            service_name = $1,
            price        = $2,
            user_id      = $3,
            start_date   = $4,
            end_date     = $5
        WHERE id = $6
    `

	res, err := r.db.ExecContext(
		ctx,
		query,
		s.ServiceName,
		s.Price,
		s.UserID,
		s.StartDate,
		s.EndDate,
		s.ID,
	)

	if err != nil {
		r.logger.ErrorContext(ctx,
			"failed to update subscription",
			"error", err,
			"id", s.ID,
			"service_name", s.ServiceName,
			"user_id", s.UserID,
		)
		return fmt.Errorf("%w: %v", port.ErrUpdateSubscription, err)
	}

	rows, err := res.RowsAffected()

	if err != nil {
		r.logger.ErrorContext(ctx,
			"failed to count rows while updating",
			"id", s.ID)
		return fmt.Errorf("%w: %v", port.ErrUpdateSubscription, err)
	}
	if rows == 0 {
		r.logger.WarnContext(ctx,
			"subscription not found",
			"id", s.ID,
		)
		return port.ErrSubscriptionNotFound
	}

	r.logger.InfoContext(ctx,
		"subscription updated",
		"id", s.ID,
		"service_name", s.ServiceName,
		"user_id", s.UserID,
	)

	return nil
}

func (r *SubscriptionRepo) DeleteSubscription(ctx context.Context, id int) error {
	query := `
		DELETE FROM subscriptions
		WHERE id = $1
	`

	res, err := r.db.ExecContext(ctx, query, id)

	if err != nil {
		r.logger.ErrorContext(ctx,
			"failed to delete subscription",
			"id", id,
		)
		return fmt.Errorf("%w: %v", port.ErrDeleteSubscription, err)
	}

	rows, err := res.RowsAffected()

	if err != nil {
		r.logger.ErrorContext(ctx,
			"failed to count rows while deleting",
			"id", id)
		return fmt.Errorf("%w: %v", port.ErrDeleteSubscription, err)
	}
	if rows == 0 {
		return port.ErrSubscriptionNotFound
	}

	r.logger.InfoContext(ctx,
		"subscription deleted",
		"id", id,
	)

	return nil
}

func (r *SubscriptionRepo) GetSubscriptionList(ctx context.Context, filter filter.ListFilter) ([]entity.Subscription, error) {
	query := `
		SELECT id, service_name, price, user_id, start_date, end_date
		FROM subscriptions
		WHERE 1=1
	`

	var args []any
	argNum := 1

	if filter.UserID != nil {
		query += fmt.Sprintf(" AND user_id = $%d", argNum)
		args = append(args, *filter.UserID)
		argNum++
	}

	if filter.ServiceName != nil {
		query += fmt.Sprintf(" AND service_name = $%d", argNum)
		args = append(args, *filter.ServiceName)
		argNum++
	}

	query += fmt.Sprintf(" LIMIT $%d", argNum)
	args = append(args, filter.Limit)
	argNum++

	query += fmt.Sprintf(" OFFSET $%d", argNum)
	args = append(args, filter.Offset)

	var subs []entity.Subscription
	if err := r.db.SelectContext(ctx, &subs, query, args...); err != nil {
		r.logger.ErrorContext(ctx,
			"failed to get subscription list",
			"service_name", filter.ServiceName,
			"user_id", filter.UserID,
			"limit", filter.Limit,
			"offset", filter.Offset,
		)
		return nil, fmt.Errorf("%w: %v", port.ErrGetSubscriptionList, err)
	}

	r.logger.InfoContext(ctx, "successfully got sub list",
		"count", len(subs),
	)

	return subs, nil
}

func (r *SubscriptionRepo) GetTotalSum(ctx context.Context, f filter.SumFilter) (int, error) {
	query := `
		SELECT SUM(price) 
		FROM subscriptions
		WHERE 1=1
	`

	var args []any
	argNum := 1

	if f.UserID != nil {
		query += fmt.Sprintf(" AND user_id = $%d", argNum)
		args = append(args, *f.UserID)
		argNum++
	}

	if f.ServiceName != nil {
		query += fmt.Sprintf(" AND service_name = $%d", argNum)
		args = append(args, *f.ServiceName)
		argNum++
	}

	var sum int
	if err := r.db.GetContext(ctx, &sum, query, args...); err != nil {
		r.logger.ErrorContext(ctx,
			"failed to get total sum",
			"service_name", f.ServiceName,
			"user_id", f.UserID,
			"start_date", f.StartDate,
			"end_date", f.EndDate,
		)
		return 0, fmt.Errorf("%w: %v", port.ErrGetTotalSum, err)
	}

	r.logger.InfoContext(ctx, "successfully got total sum",
		"total sum", sum,
	)

	return sum, nil
}
