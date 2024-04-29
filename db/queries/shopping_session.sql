-- name: CreateShoppingSession :one
insert into shopping_session 
(user_id)  
values ($1) 
returning *;

-- name: GetShoppingSessionByUserId :one
SELECT * FROM shopping_session
WHERE user_id = $1 LIMIT 1;
