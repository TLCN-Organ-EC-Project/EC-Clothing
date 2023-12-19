-- name: CreateOrder :one
INSERT INTO orders (
  booking_id,
  user_booking,
  promotion_id,
  address,
  province,
  amount,
  tax
) VALUES (
  $1, $2, $3, $4, $5, $6, $7
) RETURNING *;

-- name: GetOrder :one
SELECT * FROM orders
WHERE booking_id = $1 LIMIT 1;

-- name: ListOrderByUser :many
SELECT * FROM orders
WHERE user_booking = $1
ORDER BY booking_id
LIMIT $2
OFFSET $3;

-- name: ListOrder :many
SELECT * FROM orders
ORDER BY booking_date
LIMIT $1
OFFSET $2;

-- name: UpdateOrder :one
UPDATE orders
SET address = $2, province = $3
WHERE booking_id = $1
RETURNING *;

-- name: UpdateStatusOrder :one
UPDATE orders
SET status = $2
WHERE booking_id = $1
RETURNING *;

-- name: UpdateAmountOfOrder :one
UPDATE orders
SET amount = $2
WHERE booking_id = $1
RETURNING *;

-- name: TotalIncome :one
SELECT CAST(SUM(amount) AS FLOAT) AS TotalIncome
FROM orders  
WHERE (
booking_date BETWEEN $1 AND $2
AND status = $3);

-- name: GetOrderByDate :one
SELECT * FROM orders
WHERE (
booking_date BETWEEN $1 AND $2
AND status = $3) LIMIT 1;

-- name: DeleteOrder :exec
DELETE FROM orders WHERE booking_id = $1;