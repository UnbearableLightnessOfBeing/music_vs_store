-- name: CreateCartItem :one
insert into cart_item (
  session_id,
  product_id,
  quantity
)  values (
  $1, $2, $3
) returning *;
