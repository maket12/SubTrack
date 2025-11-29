//go:build integration
// +build integration

package db_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/maket12/SubTrack/internal/adapter/out/db"
	"github.com/maket12/SubTrack/internal/domain/entity"
	"github.com/maket12/SubTrack/internal/domain/filter"
)

func setupDB(t *testing.T) *sqlx.DB {
	dsn := "postgres://test:test@localhost:5432/testdb?sslmode=disable"

	dbx, err := sqlx.Connect("pgx", dsn)
	require.NoError(t, err)

	_, err = dbx.Exec("TRUNCATE subscriptions RESTART IDENTITY CASCADE")
	require.NoError(t, err)

	return dbx
}

func TestPostgres_Create_Get(t *testing.T) {
	dbx := setupDB(t)
	repo := db.NewSubscriptionRepo(dbx)

	sub := &entity.Subscription{
		ServiceName: "Netflix",
		Price:       500,
		UserID:      uuid.New(),
		StartDate:   time.Now(),
	}

	id, err := repo.Create(context.Background(), sub)
	require.NoError(t, err)
	require.True(t, id > 0)

	got, err := repo.Get(context.Background(), id)
	require.NoError(t, err)
	require.Equal(t, "Netflix", got.ServiceName)
	require.Equal(t, 500, got.Price)
}

func TestPostgres_Update(t *testing.T) {
	dbx := setupDB(t)
	repo := db.NewSubscriptionRepo(dbx)

	sub := &entity.Subscription{
		ServiceName: "Spotify",
		Price:       300,
		UserID:      uuid.New(),
		StartDate:   time.Now(),
	}
	id, _ := repo.Create(context.Background(), sub)

	sub.ID = id
	sub.Price = 999

	err := repo.Update(context.Background(), sub)
	require.NoError(t, err)

	got, err := repo.Get(context.Background(), id)
	require.NoError(t, err)
	require.Equal(t, 999, got.Price)
}

func TestPostgres_Delete(t *testing.T) {
	dbx := setupDB(t)
	repo := db.NewSubscriptionRepo(dbx)

	sub := &entity.Subscription{
		ServiceName: "Apple",
		Price:       777,
		UserID:      uuid.New(),
		StartDate:   time.Now(),
	}
	id, _ := repo.Create(context.Background(), sub)

	err := repo.Delete(context.Background(), id)
	require.NoError(t, err)

	_, err = repo.Get(context.Background(), id)
	require.Error(t, err) // sql.ErrNoRows
}

func TestPostgres_GetList(t *testing.T) {
	dbx := setupDB(t)
	repo := db.NewSubscriptionRepo(dbx)

	uid := uuid.New()

	repo.Create(context.Background(), &entity.Subscription{
		ServiceName: "A",
		Price:       100,
		UserID:      uid,
		StartDate:   time.Now(),
	})
	repo.Create(context.Background(), &entity.Subscription{
		ServiceName: "B",
		Price:       200,
		UserID:      uid,
		StartDate:   time.Now(),
	})

	list, err := repo.GetList(context.Background(),
		filter.ListFilter{UserID: &uid, Limit: 10, Offset: 0},
	)
	require.NoError(t, err)
	require.Len(t, list, 2)
}

func TestPostgres_GetTotalSum(t *testing.T) {
	dbx := setupDB(t)
	repo := db.NewSubscriptionRepo(dbx)

	uid := uuid.New()

	repo.Create(context.Background(), &entity.Subscription{
		ServiceName: "A",
		Price:       300,
		UserID:      uid,
		StartDate:   time.Now(),
	})
	repo.Create(context.Background(), &entity.Subscription{
		ServiceName: "B",
		Price:       700,
		UserID:      uid,
		StartDate:   time.Now(),
	})

	sum, err := repo.GetTotalSum(context.Background(),
		filter.SumFilter{UserID: &uid},
	)
	require.NoError(t, err)
	require.Equal(t, 1000, sum)
}
