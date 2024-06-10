-- name: ListLabels :many
select * from labels
order by id
limit $1
offset $2;

-- name: GetLabel :one
select * from labels
where id = $1
limit 1;

-- name: CreateLabel :one
INSERT INTO labels
(name) values ($1)
RETURNING *;

-- name: UpdateLabel :one
UPDATE labels
  SET name = $2
WHERE id = $1
RETURNING *;

-- name: DeleteLabel :one
DELETE FROM labels
WHERE id = $1
RETURNING *;

-- name: RemoveLabelProductRelations :exec
UPDATE products
  SET label_id = NULL
WHERE label_id = $1;
