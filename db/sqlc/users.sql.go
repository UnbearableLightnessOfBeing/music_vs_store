// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: users.sql

package db

import (
	"context"
	"database/sql"
)

const createUser = `-- name: CreateUser :one
insert into users (
  username,
  email,
  password
)  values (
  $1, $2, $3
) returning id, username, email, is_admin, password, created_at
`

type CreateUserParams struct {
	Username string `form:"username" json:"username"`
	Email    string `form:"email" json:"email"`
	Password string `form:"password" json:"password"`
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.db.QueryRowContext(ctx, createUser, arg.Username, arg.Email, arg.Password)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.Email,
		&i.IsAdmin,
		&i.Password,
		&i.CreatedAt,
	)
	return i, err
}

const deleteUser = `-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1
`

func (q *Queries) DeleteUser(ctx context.Context, id int32) error {
	_, err := q.db.ExecContext(ctx, deleteUser, id)
	return err
}

const deleteUserByName = `-- name: DeleteUserByName :exec
delete from users
where username = $1
`

func (q *Queries) DeleteUserByName(ctx context.Context, username string) error {
	_, err := q.db.ExecContext(ctx, deleteUserByName, username)
	return err
}

const getUser = `-- name: GetUser :one
SELECT id, username, email, is_admin, password, created_at FROM users
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetUser(ctx context.Context, id int32) (User, error) {
	row := q.db.QueryRowContext(ctx, getUser, id)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.Email,
		&i.IsAdmin,
		&i.Password,
		&i.CreatedAt,
	)
	return i, err
}

const getUserByEmail = `-- name: GetUserByEmail :one
SELECT id, username, email, is_admin, password, created_at FROM users
WHERE email = $1 LIMIT 1
`

func (q *Queries) GetUserByEmail(ctx context.Context, email string) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserByEmail, email)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.Email,
		&i.IsAdmin,
		&i.Password,
		&i.CreatedAt,
	)
	return i, err
}

const getUserByName = `-- name: GetUserByName :one
SELECT id, username, email, is_admin, password, created_at FROM users
WHERE username = $1 LIMIT 1
`

func (q *Queries) GetUserByName(ctx context.Context, username string) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserByName, username)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.Email,
		&i.IsAdmin,
		&i.Password,
		&i.CreatedAt,
	)
	return i, err
}

const listUsers = `-- name: ListUsers :many
select id, username, email, is_admin from users
order by id
limit $1
offset $2
`

type ListUsersParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

type ListUsersRow struct {
	ID       int32        `json:"id"`
	Username string       `form:"username" json:"username"`
	Email    string       `form:"email" json:"email"`
	IsAdmin  sql.NullBool `json:"is_admin"`
}

func (q *Queries) ListUsers(ctx context.Context, arg ListUsersParams) ([]ListUsersRow, error) {
	rows, err := q.db.QueryContext(ctx, listUsers, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ListUsersRow
	for rows.Next() {
		var i ListUsersRow
		if err := rows.Scan(
			&i.ID,
			&i.Username,
			&i.Email,
			&i.IsAdmin,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateUserIsAdmin = `-- name: UpdateUserIsAdmin :one
update users
set is_admin = $2
where id = $1
returning id, username, email, is_admin, password, created_at
`

type UpdateUserIsAdminParams struct {
	ID      int32        `json:"id"`
	IsAdmin sql.NullBool `json:"is_admin"`
}

func (q *Queries) UpdateUserIsAdmin(ctx context.Context, arg UpdateUserIsAdminParams) (User, error) {
	row := q.db.QueryRowContext(ctx, updateUserIsAdmin, arg.ID, arg.IsAdmin)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.Email,
		&i.IsAdmin,
		&i.Password,
		&i.CreatedAt,
	)
	return i, err
}

const updateUserName = `-- name: UpdateUserName :one
UPDATE users
SET username = $2
WHERE id = $1
RETURNING id, username, email, is_admin, password, created_at
`

type UpdateUserNameParams struct {
	ID       int32  `json:"id"`
	Username string `form:"username" json:"username"`
}

func (q *Queries) UpdateUserName(ctx context.Context, arg UpdateUserNameParams) (User, error) {
	row := q.db.QueryRowContext(ctx, updateUserName, arg.ID, arg.Username)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.Email,
		&i.IsAdmin,
		&i.Password,
		&i.CreatedAt,
	)
	return i, err
}
