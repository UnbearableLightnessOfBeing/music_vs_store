-- name: ListPaymentMethods :many
select * from payment_methods
limit $1
offset $2;
