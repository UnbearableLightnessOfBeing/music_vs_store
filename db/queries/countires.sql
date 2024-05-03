-- name: ListCountries :many
select * from countries
order by id
limit $1
offset $2;
