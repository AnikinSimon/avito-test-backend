// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0
// source: reception.sql

package db

import (
	"context"
	"time"

	entity "github.com/AnikinSimon/avito-test-backend/internal/models/entity"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

const addProductToReception = `-- name: AddProductToReception :one
INSERT INTO products (id, type, reception_id) VALUES
    ($1, $2, $3)
    RETURNING id, date_time, type, reception_id
`

type AddProductToReceptionParams struct {
	ID          uuid.UUID
	Type        entity.ProductType
	ReceptionID uuid.UUID
}

func (q *Queries) AddProductToReception(ctx context.Context, arg AddProductToReceptionParams) (Product, error) {
	row := q.db.QueryRowContext(ctx, addProductToReception, arg.ID, arg.Type, arg.ReceptionID)
	var i Product
	err := row.Scan(
		&i.ID,
		&i.DateTime,
		&i.Type,
		&i.ReceptionID,
	)
	return i, err
}

const createReception = `-- name: CreateReception :one
INSERT INTO receptions (id, date_time, pvz_id) VALUES
    ($1, $2, $3)
    RETURNING id, date_time, pvz_id, status
`

type CreateReceptionParams struct {
	ID       uuid.UUID
	DateTime time.Time
	PvzID    uuid.UUID
}

func (q *Queries) CreateReception(ctx context.Context, arg CreateReceptionParams) (Reception, error) {
	row := q.db.QueryRowContext(ctx, createReception, arg.ID, arg.DateTime, arg.PvzID)
	var i Reception
	err := row.Scan(
		&i.ID,
		&i.DateTime,
		&i.PvzID,
		&i.Status,
	)
	return i, err
}

const deleteProduct = `-- name: DeleteProduct :exec
DELETE FROM products
WHERE id = $1
`

func (q *Queries) DeleteProduct(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.ExecContext(ctx, deleteProduct, id)
	return err
}

const finishReception = `-- name: FinishReception :one
UPDATE receptions
SET status='close'
WHERE id=(SELECT id FROM receptions R WHERE R.pvz_id=$1 AND R.status='in_progress' LIMIT 1)
    RETURNING id, date_time, pvz_id, status
`

func (q *Queries) FinishReception(ctx context.Context, pvzID uuid.UUID) (Reception, error) {
	row := q.db.QueryRowContext(ctx, finishReception, pvzID)
	var i Reception
	err := row.Scan(
		&i.ID,
		&i.DateTime,
		&i.PvzID,
		&i.Status,
	)
	return i, err
}

const getLastProductInReception = `-- name: GetLastProductInReception :one
SELECT id, date_time, type, reception_id FROM products
WHERE reception_id = $1
ORDER BY date_time DESC
    LIMIT 1
`

func (q *Queries) GetLastProductInReception(ctx context.Context, receptionID uuid.UUID) (Product, error) {
	row := q.db.QueryRowContext(ctx, getLastProductInReception, receptionID)
	var i Product
	err := row.Scan(
		&i.ID,
		&i.DateTime,
		&i.Type,
		&i.ReceptionID,
	)
	return i, err
}

const getOpenReceptionByPvzID = `-- name: GetOpenReceptionByPvzID :one
SELECT id, date_time, pvz_id, status FROM receptions
WHERE pvz_id = $1 AND status = 'in_progress'
    LIMIT 1
`

func (q *Queries) GetOpenReceptionByPvzID(ctx context.Context, pvzID uuid.UUID) (Reception, error) {
	row := q.db.QueryRowContext(ctx, getOpenReceptionByPvzID, pvzID)
	var i Reception
	err := row.Scan(
		&i.ID,
		&i.DateTime,
		&i.PvzID,
		&i.Status,
	)
	return i, err
}

const getProductsFromReception = `-- name: GetProductsFromReception :many
SELECT id, date_time, type, reception_id FROM products
WHERE reception_id IN ($1)
`

func (q *Queries) GetProductsFromReception(ctx context.Context, receptionID uuid.UUID) ([]Product, error) {
	rows, err := q.db.QueryContext(ctx, getProductsFromReception, receptionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Product{}
	for rows.Next() {
		var i Product
		if err := rows.Scan(
			&i.ID,
			&i.DateTime,
			&i.Type,
			&i.ReceptionID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const searchReceptionsByPvzsAndTime = `-- name: SearchReceptionsByPvzsAndTime :many
SELECT id, date_time, pvz_id, status FROM receptions
WHERE pvz_id = ANY($1::uuid[]) AND date_time BETWEEN $2 AND $3
`

type SearchReceptionsByPvzsAndTimeParams struct {
	PvzIds    []uuid.UUID
	StartDate time.Time
	EndDate   time.Time
}

func (q *Queries) SearchReceptionsByPvzsAndTime(ctx context.Context, arg SearchReceptionsByPvzsAndTimeParams) ([]Reception, error) {
	rows, err := q.db.QueryContext(ctx, searchReceptionsByPvzsAndTime, pq.Array(arg.PvzIds), arg.StartDate, arg.EndDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Reception{}
	for rows.Next() {
		var i Reception
		if err := rows.Scan(
			&i.ID,
			&i.DateTime,
			&i.PvzID,
			&i.Status,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const searchReceptionsByTime = `-- name: SearchReceptionsByTime :many
SELECT id, date_time, pvz_id, status FROM receptions
WHERE date_time BETWEEN $1 AND $2
`

type SearchReceptionsByTimeParams struct {
	DateTime   time.Time
	DateTime_2 time.Time
}

func (q *Queries) SearchReceptionsByTime(ctx context.Context, arg SearchReceptionsByTimeParams) ([]Reception, error) {
	rows, err := q.db.QueryContext(ctx, searchReceptionsByTime, arg.DateTime, arg.DateTime_2)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Reception{}
	for rows.Next() {
		var i Reception
		if err := rows.Scan(
			&i.ID,
			&i.DateTime,
			&i.PvzID,
			&i.Status,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
