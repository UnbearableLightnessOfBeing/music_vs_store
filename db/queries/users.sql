-- name: CreateUser :one
insert into users (
  username,
  email,
  password
)  values (
  $1, $2, $3
) returning *;

-- name: GetUser :one
SELECT * FROM users
WHERE id = $1 LIMIT 1;

-- name: ListUsrs :many
select * from users
order by id
limit $1
offset $2;

-- name: UpdateUserIsAdmin :one
update users
set is_admin = $2
where id = $1
returning *;

-- name: UpdateUserName :one
UPDATE users
SET username = $2
WHERE id = $1
RETURNING *;

-- name: DeleteUser :exec
delete from users
where id = $1;

-- name: DeleteUserByName :exec
delete from users
where username = $1;
