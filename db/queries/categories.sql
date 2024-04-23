-- name: CreateCategory :one
insert into categories 
(name)  
values ($1) 
returning *;

-- name: GetCategory :one
SELECT * FROM categories
WHERE id = $1 LIMIT 1;

-- name: ListCategories :many
select * from categories
order by id
limit $1
offset $2;

