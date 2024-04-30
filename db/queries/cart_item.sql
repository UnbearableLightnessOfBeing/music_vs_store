-- name: CreateCartItem :one
INSERT INTO cart_item (
  session_id,
  product_id,
  quantity
)  VALUES (
  $1, $2, $3
) RETURNING *;

-- name: GetCartItem :one
SELECT * FROM cart_item
WHERE session_id = $1 
  AND product_id = $2
LIMIT 1;

-- name: UpdateCartItemQuantity :one
UPDATE cart_item
SET quantity = $2
WHERE id = $1
RETURNING *;

-- name: DeleteCartItem :one
DELETE FROM  cart_item
WHERE session_id = $1
  AND product_id = $2
RETURNING *;
