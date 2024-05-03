// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: countires.sql

package db

import (
	"context"
)

const listCountries = `-- name: ListCountries :many
select id, name from countries
order by id
limit $1
offset $2
`

type ListCountriesParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

func (q *Queries) ListCountries(ctx context.Context, arg ListCountriesParams) ([]Country, error) {
	rows, err := q.db.QueryContext(ctx, listCountries, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Country
	for rows.Next() {
		var i Country
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
