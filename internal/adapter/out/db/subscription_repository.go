package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/maket12/SubTrack/internal/domain/entity"
	"github.com/maket12/SubTrack/internal/domain/filter"

	"github.com/jmoiron/sqlx"
)

type SubscriptionRepository struct {
	db *sqlx.DB
}

func NewSubscriptionRepo(db *sqlx.DB) *SubscriptionRepository {
	return &SubscriptionRepository{
		db: db,
	}
}

func (r *SubscriptionRepository) Create(ctx context.Context, s *entity.Subscription) (int, error) {
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
		return 0, fmt.Errorf("failed to create subscription using db: %w", err)
	}

	return id, nil
}

func (r *SubscriptionRepository) Get(ctx context.Context, id int) (*entity.Subscription, error) {
	query := `
		SELECT id, service_name, price, user_id, start_date, end_date
		FROM subscriptions
		WHERE id = $1
	`

	var sub entity.Subscription

	err := r.db.GetContext(ctx, &sub, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
		return nil, fmt.Errorf("failed to get subscription using db: %w", err)
	}

	return &sub, nil
}

func (r *SubscriptionRepository) Update(ctx context.Context, s *entity.Subscription) error {
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
		return fmt.Errorf("failed to update subscription using db: %w", err)
	}

	rows, err := res.RowsAffected()

	if err != nil {
		return fmt.Errorf("failed to count rows: %w", err)
	}
	if rows == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *SubscriptionRepository) Delete(ctx context.Context, id int) error {
	query := `
		DELETE FROM subscriptions
		WHERE id = $1
	`

	res, err := r.db.ExecContext(ctx, query, id)

	if err != nil {
		return fmt.Errorf("failed to delete subscription using db: %w", err)
	}

	rows, err := res.RowsAffected()

	if err != nil {
		return fmt.Errorf("failed to delete subscription using db: %w", err)
	}
	if rows == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *SubscriptionRepository) GetList(ctx context.Context, f filter.ListFilter) ([]entity.Subscription, error) {
	query := `
		SELECT id, service_name, price, user_id, start_date, end_date
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

	query += fmt.Sprintf(" LIMIT $%d", argNum)
	args = append(args, f.Limit)
	argNum++

	query += fmt.Sprintf(" OFFSET $%d", argNum)
	args = append(args, f.Offset)

	var subs []entity.Subscription
	if err := r.db.SelectContext(ctx, &subs, query, args...); err != nil {
		return nil, fmt.Errorf("failed to get list of subscriptions using db: %w", err)
	}

	return subs, nil
}

func (r *SubscriptionRepository) GetTotalSum(ctx context.Context, f filter.SumFilter) (int, error) {
	query := `
		SELECT COALESCE(SUM(price), 0)
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

	if f.StartDate != nil {
		query += fmt.Sprintf(" AND start_date >= $%d", argNum)
		args = append(args, *f.StartDate)
		argNum++
	}

	if f.EndDate != nil {
		query += fmt.Sprintf(" AND start_date <= $%d", argNum)
		args = append(args, *f.EndDate)
	}

	var sum int
	if err := r.db.GetContext(ctx, &sum, query, args...); err != nil {
		return 0, fmt.Errorf("failed to get total sum using db: %w", err)
	}

	return sum, nil
}
