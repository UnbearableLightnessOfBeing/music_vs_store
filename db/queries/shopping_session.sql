-- name: CreateShoppingSession :one
insert into shopping_session 
(user_id)  
values ($1) 
returning *;

-- name: GetShoppingSessionByUserId :one
SELECT * FROM shopping_session
WHERE user_id = $1 LIMIT 1;

-- name: DeleteSessionByUserId :one
DELETE FROM shopping_session
WHERE user_id = $1
RETURNING *;

-- name: UpdateSessionTotal :one
UPDATE shopping_session
SET total_int = $2
WHERE user_id = $1
RETURNING *;
