-- name: ListLabels :many
select * from labels
order by id
limit $1
offset $2;
