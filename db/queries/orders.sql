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

-- name: GetOrder :one
select o.*,
  TO_CHAR(o.created_at, 'DD.MM.YYYY HH:MM') as created_formatted ,
  p.name as payment_name,
  d.name as delivery_name
  from orders o
  left join payment_methods p
  on o.payment_method_id = p.id
  left join delivery_methods d
  on o.delivery_method_id = d.id
where o.id = $1
limit 1;

-- name: GetOrderProducts :many
select 
  p.name,
  p.images,
  p_o.count,
  p.price_int,
  p.price_int * p_o.count as product_total
from products p join product_orders p_o
  on p.id = p_o.product_id
where p_o.order_id = $1
order by p.id;
  
-- name: GetOrders :many
SELECT o.*,
  p.name as payment_method,
  d.name as delivery_method
FROM orders o 
  LEFT JOIN payment_methods p
  ON o.payment_method_id = p.id
  LEFT JOIN delivery_methods d
  ON o.delivery_method_id = d.id;

-- name: GetOrderCount :one
SELECT count(*) FROM orders;

-- name: GetTotalRevenue :one
SELECT SUM(total_int) FROM orders;
