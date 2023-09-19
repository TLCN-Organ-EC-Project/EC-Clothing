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

-- name: DeleteOrder :exec
DELETE FROM orders WHERE booking_id = $1;