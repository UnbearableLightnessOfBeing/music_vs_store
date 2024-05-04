-- name: GetProduct :one
SELECT * FROM products
WHERE id = $1 LIMIT 1;

-- name: GetProductsByCategory :many
SELECT p.* FROM products p, product_categories p_c
WHERE p.id = p_c.product_id
  AND $1 = p_c.category_id 
  AND (CASE WHEN @min_price::integer > 0 THEN p.price_int >= @min_price ELSE p.price_int > 0 END)
  AND (CASE WHEN @max_price::integer > 0 THEN p.price_int <= @max_price ELSE p.price_int < 999999 END)
  AND (CASE WHEN @label_id::integer > 0 THEN p.label_id = @label_id ELSE TRUE END)
ORDER BY 
CASE WHEN @price_sorting::varchar(10) = 'ASC' THEN p.price_int END asc,
CASE WHEN @price_sorting::varchar(10) = 'DESC' THEN p.price_int END desc;

-- name: GetProdutsInCart :many
SELECT p.*, c_i.quantity FROM products p, cart_item c_i
WHERE p.id = c_i.product_id
  AND c_i.session_id = $1
ORDER BY p.id asc;

-- name: GetCartProductsCount :one
SELECT COUNT(*) FROM products p, cart_item c_i
WHERE p.id = c_i.product_id
  AND c_i.session_id = $1;

-- name: SearchProducts :many
select * from products
where LOWER(name) like $1;
