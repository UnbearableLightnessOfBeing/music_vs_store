// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: products.sql

package db

import (
	"context"

	"github.com/lib/pq"
)

const getCartProductsBySessionId = `-- name: GetCartProductsBySessionId :many
SELECT p.name, p.price_int, c_i.quantity FROM products p, cart_item c_i
WHERE p.id = c_i.product_id
  AND c_i.session_id = $1
`

type GetCartProductsBySessionIdRow struct {
	Name     string `json:"name"`
	PriceInt int32  `json:"price_int"`
	Quantity int32  `json:"quantity"`
}

func (q *Queries) GetCartProductsBySessionId(ctx context.Context, sessionID int32) ([]GetCartProductsBySessionIdRow, error) {
	rows, err := q.db.QueryContext(ctx, getCartProductsBySessionId, sessionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetCartProductsBySessionIdRow
	for rows.Next() {
		var i GetCartProductsBySessionIdRow
		if err := rows.Scan(&i.Name, &i.PriceInt, &i.Quantity); err != nil {
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

const getProduct = `-- name: GetProduct :one
SELECT id, name, price_int, price_dec, label_id, images, description, in_stock FROM products
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetProduct(ctx context.Context, id int32) (Product, error) {
	row := q.db.QueryRowContext(ctx, getProduct, id)
	var i Product
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.PriceInt,
		&i.PriceDec,
		&i.LabelID,
		pq.Array(&i.Images),
		&i.Description,
		&i.InStock,
	)
	return i, err
}

const getProductsByCategory = `-- name: GetProductsByCategory :many
SELECT p.id, p.name, p.price_int, p.price_dec, p.label_id, p.images, p.description, p.in_stock FROM products p, product_categories p_c
WHERE p.id = p_c.product_id
  AND $1 = p_c.category_id 
  AND (CASE WHEN $2::integer > 0 THEN p.price_int >= $2 ELSE p.price_int > 0 END)
  AND (CASE WHEN $3::integer > 0 THEN p.price_int <= $3 ELSE p.price_int < 999999 END)
  AND (CASE WHEN $4::integer > 0 THEN p.label_id = $4 ELSE TRUE END)
ORDER BY 
CASE WHEN $5::varchar(10) = 'ASC' THEN p.price_int END asc,
CASE WHEN $5::varchar(10) = 'DESC' THEN p.price_int END desc
`

type GetProductsByCategoryParams struct {
	CategoryID   int32  `json:"category_id"`
	MinPrice     int32  `json:"min_price"`
	MaxPrice     int32  `json:"max_price"`
	LabelID      int32  `json:"label_id"`
	PriceSorting string `json:"price_sorting"`
}

func (q *Queries) GetProductsByCategory(ctx context.Context, arg GetProductsByCategoryParams) ([]Product, error) {
	rows, err := q.db.QueryContext(ctx, getProductsByCategory,
		arg.CategoryID,
		arg.MinPrice,
		arg.MaxPrice,
		arg.LabelID,
		arg.PriceSorting,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Product
	for rows.Next() {
		var i Product
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.PriceInt,
			&i.PriceDec,
			&i.LabelID,
			pq.Array(&i.Images),
			&i.Description,
			&i.InStock,
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
