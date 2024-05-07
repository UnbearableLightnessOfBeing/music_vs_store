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
select o.*, 
  TO_CHAR(o.created_at, 'DD.MM.YYYY HH:MM:SS') as created_formatted ,
  p.name as payment_name,
  d.name as delivery_name
  from orders o
  left join payment_methods p
  on o.payment_method_id = p.id
  left join delivery_methods d
  on o.delivery_method_id = d.id
where o.user_id = $1
order by o.id;
