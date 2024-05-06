-- name: CreateOrder :one
insert into orders (
	user_id,
	product_count,
	price_int,   
	delivery_price_int,
	total_int,            
	country_id,
	district,
	city,
	postal_code,
	delivery_method_id,
	payment_method_id,
	customer_firstname,
	cusotmer_middlename,
	customer_lastname,
	customer_phone_number,
	customer_email,
	customer_address,
	customer_comment
)  values (
  $1, 
  $2, 
  $3,
  $4,
  $5,
  $6,
  $7,
  $8,
  $9,
  $10,
  $11,
  $12,
  $13,
  $14,
  $15,
  $16,
  $17,
  $18
) returning id;

-- name: GetOrdersByUserId :many
select *, TO_CHAR(created_at, 'DD.MM.YYYY HH:MM:SS') as created_formatted from orders
where user_id = $1
  order by id;
