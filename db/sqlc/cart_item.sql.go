// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: cart_item.sql

package db

import (
	"context"
)

const createCartItem = `-- name: CreateCartItem :one
INSERT INTO cart_item (
  session_id,
  product_id,
  quantity
)  VALUES (
  $1, $2, $3
) RETURNING id, session_id, product_id, quantity
`

type CreateCartItemParams struct {
	SessionID int32 `json:"session_id"`
	ProductID int32 `json:"product_id"`
	Quantity  int32 `json:"quantity"`
}

func (q *Queries) CreateCartItem(ctx context.Context, arg CreateCartItemParams) (CartItem, error) {
	row := q.db.QueryRowContext(ctx, createCartItem, arg.SessionID, arg.ProductID, arg.Quantity)
	var i CartItem
	err := row.Scan(
		&i.ID,
		&i.SessionID,
		&i.ProductID,
		&i.Quantity,
	)
	return i, err
}

const deleteCartItem = `-- name: DeleteCartItem :one
DELETE FROM  cart_item
WHERE session_id = $1
  AND product_id = $2
RETURNING id, session_id, product_id, quantity
`

type DeleteCartItemParams struct {
	SessionID int32 `json:"session_id"`
	ProductID int32 `json:"product_id"`
}

func (q *Queries) DeleteCartItem(ctx context.Context, arg DeleteCartItemParams) (CartItem, error) {
	row := q.db.QueryRowContext(ctx, deleteCartItem, arg.SessionID, arg.ProductID)
	var i CartItem
	err := row.Scan(
		&i.ID,
		&i.SessionID,
		&i.ProductID,
		&i.Quantity,
	)
	return i, err
}

const getCartItem = `-- name: GetCartItem :one
SELECT id, session_id, product_id, quantity FROM cart_item
WHERE session_id = $1 
  AND product_id = $2
LIMIT 1
`

type GetCartItemParams struct {
	SessionID int32 `json:"session_id"`
	ProductID int32 `json:"product_id"`
}

func (q *Queries) GetCartItem(ctx context.Context, arg GetCartItemParams) (CartItem, error) {
	row := q.db.QueryRowContext(ctx, getCartItem, arg.SessionID, arg.ProductID)
	var i CartItem
	err := row.Scan(
		&i.ID,
		&i.SessionID,
		&i.ProductID,
		&i.Quantity,
	)
	return i, err
}

const updateCartItemQuantity = `-- name: UpdateCartItemQuantity :one
UPDATE cart_item
SET quantity = $2
WHERE id = $1
RETURNING id, session_id, product_id, quantity
`

type UpdateCartItemQuantityParams struct {
	ID       int32 `json:"id"`
	Quantity int32 `json:"quantity"`
}

func (q *Queries) UpdateCartItemQuantity(ctx context.Context, arg UpdateCartItemQuantityParams) (CartItem, error) {
	row := q.db.QueryRowContext(ctx, updateCartItemQuantity, arg.ID, arg.Quantity)
	var i CartItem
	err := row.Scan(
		&i.ID,
		&i.SessionID,
		&i.ProductID,
		&i.Quantity,
	)
	return i, err
}
