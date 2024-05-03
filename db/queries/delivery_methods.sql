-- name: ListDeliveryMethods :many
select * from delivery_methods
limit $1
offset $2;
