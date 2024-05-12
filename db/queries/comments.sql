-- name: GetComments :many
SELECT * FROM comments;

-- name: CreateComment :one
INSERT INTO comments
  (user_id, name, text) values ($1, $2, $3)
RETURNING *;
