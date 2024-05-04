-- name: AddProductToOrder :many
INSERT INTO product_orders
(order_id, product_id, count)
VALUES ($1, $2, $3)
RETURNING *;
