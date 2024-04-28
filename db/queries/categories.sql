-- name: CreateCategory :one
insert into categories 
(name)  
values ($1) 
returning *;

-- name: GetCategory :one
SELECT * FROM categories
WHERE id = $1 LIMIT 1;

-- name: GetCategoryBySlug :one
SELECT * FROM categories
WHERE slug = $1 LIMIT 1;

-- name: UpdateCategoryName :one
UPDATE categories
SET name = $2
WHERE id = $1
RETURNING *;

-- name: UpdateCategoryImageUrl :one
UPDATE categories
SET img_url = $2
WHERE id = $1
RETURNING *;

-- name: ListCategories :many
select * from categories
order by id
limit $1
offset $2;

