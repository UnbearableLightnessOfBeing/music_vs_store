-- name: ListProducts :many
select * from products
limit $1
offset $2;

-- name: GetProductCount :one
SELECT count(*) FROM products;

-- name: GetProductsWithCategory :many
SELECT p.*, p_c.category_id as category_id 
  FROM products p
  LEFT JOIN product_categories p_c
  ON p.id = p_c.product_id;

-- name: CreateProduct :one
INSERT INTO products
  (name, price_int, label_id, description, characteristics, in_stock)
  VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: UpdateProduct :one
UPDATE products
  SET name = $2, 
  price_int = $3, 
  label_id = $4, 
  description = $5, 
  characteristics = $6, 
  in_stock = $7 
WHERE id = $1
RETURNING *;

-- name: GetProductByName :one
SELECT * FROM products
WHERE name = $1
LIMIT 1;

-- name: GetProduct :one
SELECT * FROM products
WHERE id = $1 
LIMIT 1;

-- name: GetProductWithCategory :one
SELECT p.*, p_c.category_id as category_id
  FROM products p
  LEFT JOIN product_categories p_c
  ON p.id = p_c.product_id
WHERE p.id = $1
LIMIT 1;

-- name: DeleteProductCategoryRelations :one
DELETE FROM product_categories p_c
WHERE p_c.product_id = $1
RETURNING *;

-- name: AddProductCategoryRelation :one
INSERT INTO product_categories
  ( product_id, category_id )
  VALUES ( $1, $2 )
RETURNING *;

-- name: AddImageToProduct :exec
UPDATE products
SET images = array_append( images, $2 )
WHERE id = $1;

-- name: RemoveImageFromProduct :exec
UPDATE products
SET images = array_remove( images, $2 )
WHERE id = $1;

-- name: DeleteProduct :one
DELETE FROM products
WHERE id = $1
RETURNING *;

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
where LOWER(name) like $1
 OR LOWER(description) like $1
order by id;
