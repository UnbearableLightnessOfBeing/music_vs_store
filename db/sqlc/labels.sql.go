// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: labels.sql

package db

import (
	"context"
	"database/sql"
)

const createLabel = `-- name: CreateLabel :one
INSERT INTO labels
(name) values ($1)
RETURNING id, name
`

func (q *Queries) CreateLabel(ctx context.Context, name string) (Label, error) {
	row := q.db.QueryRowContext(ctx, createLabel, name)
	var i Label
	err := row.Scan(&i.ID, &i.Name)
	return i, err
}

const deleteLabel = `-- name: DeleteLabel :one
DELETE FROM labels
WHERE id = $1
RETURNING id, name
`

func (q *Queries) DeleteLabel(ctx context.Context, id int32) (Label, error) {
	row := q.db.QueryRowContext(ctx, deleteLabel, id)
	var i Label
	err := row.Scan(&i.ID, &i.Name)
	return i, err
}

const getLabel = `-- name: GetLabel :one
select id, name from labels
where id = $1
limit 1
`

func (q *Queries) GetLabel(ctx context.Context, id int32) (Label, error) {
	row := q.db.QueryRowContext(ctx, getLabel, id)
	var i Label
	err := row.Scan(&i.ID, &i.Name)
	return i, err
}

const listLabels = `-- name: ListLabels :many
select id, name from labels
order by id
limit $1
offset $2
`

type ListLabelsParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

func (q *Queries) ListLabels(ctx context.Context, arg ListLabelsParams) ([]Label, error) {
	rows, err := q.db.QueryContext(ctx, listLabels, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Label
	for rows.Next() {
		var i Label
		if err := rows.Scan(&i.ID, &i.Name); err != nil {
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

const removeLabelProductRelations = `-- name: RemoveLabelProductRelations :exec
UPDATE products
  SET label_id = NULL
WHERE label_id = $1
`

func (q *Queries) RemoveLabelProductRelations(ctx context.Context, labelID sql.NullInt32) error {
	_, err := q.db.ExecContext(ctx, removeLabelProductRelations, labelID)
	return err
}

const updateLabel = `-- name: UpdateLabel :one
UPDATE labels
  SET name = $2
WHERE id = $1
RETURNING id, name
`

type UpdateLabelParams struct {
	ID   int32  `json:"id"`
	Name string `json:"name"`
}

func (q *Queries) UpdateLabel(ctx context.Context, arg UpdateLabelParams) (Label, error) {
	row := q.db.QueryRowContext(ctx, updateLabel, arg.ID, arg.Name)
	var i Label
	err := row.Scan(&i.ID, &i.Name)
	return i, err
}
