// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.16.0
// source: store.sql

package db

import (
	"context"
)

const createStore = `-- name: CreateStore :one
INSERT INTO store (
  product_id,
  size,
  quantity
) VALUES (
  $1, $2, $3
) RETURNING id, product_id, size, quantity
`

type CreateStoreParams struct {
	ProductID int64  `json:"product_id"`
	Size      string `json:"size"`
	Quantity  int32  `json:"quantity"`
}

func (q *Queries) CreateStore(ctx context.Context, arg CreateStoreParams) (Store, error) {
	row := q.db.QueryRowContext(ctx, createStore, arg.ProductID, arg.Size, arg.Quantity)
	var i Store
	err := row.Scan(
		&i.ID,
		&i.ProductID,
		&i.Size,
		&i.Quantity,
	)
	return i, err
}

const deleteStore = `-- name: DeleteStore :exec
DELETE FROM store WHERE product_id = $1 AND size = $2
`

type DeleteStoreParams struct {
	ProductID int64  `json:"product_id"`
	Size      string `json:"size"`
}

func (q *Queries) DeleteStore(ctx context.Context, arg DeleteStoreParams) error {
	_, err := q.db.ExecContext(ctx, deleteStore, arg.ProductID, arg.Size)
	return err
}

const getStore = `-- name: GetStore :one
SELECT id, product_id, size, quantity FROM store
WHERE product_id = $1 AND size = $2
LIMIT 1
`

type GetStoreParams struct {
	ProductID int64  `json:"product_id"`
	Size      string `json:"size"`
}

func (q *Queries) GetStore(ctx context.Context, arg GetStoreParams) (Store, error) {
	row := q.db.QueryRowContext(ctx, getStore, arg.ProductID, arg.Size)
	var i Store
	err := row.Scan(
		&i.ID,
		&i.ProductID,
		&i.Size,
		&i.Quantity,
	)
	return i, err
}

const listStore = `-- name: ListStore :many
SELECT id, product_id, size, quantity FROM store
WHERE product_id = $1
ORDER BY id
`

func (q *Queries) ListStore(ctx context.Context, productID int64) ([]Store, error) {
	rows, err := q.db.QueryContext(ctx, listStore, productID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Store{}
	for rows.Next() {
		var i Store
		if err := rows.Scan(
			&i.ID,
			&i.ProductID,
			&i.Size,
			&i.Quantity,
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

const updateStore = `-- name: UpdateStore :one
UPDATE store
SET quantity = $3
WHERE product_id = $1 AND size = $2
RETURNING id, product_id, size, quantity
`

type UpdateStoreParams struct {
	ProductID int64  `json:"product_id"`
	Size      string `json:"size"`
	Quantity  int32  `json:"quantity"`
}

func (q *Queries) UpdateStore(ctx context.Context, arg UpdateStoreParams) (Store, error) {
	row := q.db.QueryRowContext(ctx, updateStore, arg.ProductID, arg.Size, arg.Quantity)
	var i Store
	err := row.Scan(
		&i.ID,
		&i.ProductID,
		&i.Size,
		&i.Quantity,
	)
	return i, err
}
