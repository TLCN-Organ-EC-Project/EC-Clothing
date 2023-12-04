-- name: CreateItemsOrder :one
INSERT INTO items_order (
  booking_id,
  product_id,
  quantity,
  price,
  size
) VALUES (
  $1, $2, $3, $4, $5
) RETURNING *;

-- name: GetItemsOrder :one
SELECT * FROM items_order
WHERE id = $1 LIMIT 1;

-- name: ListItemsOrderByBookingID :many
SELECT * FROM items_order
WHERE booking_id = $1
ORDER BY id;

-- name: StatisticsProduct :many
SELECT product_id, COUNT(*) AS quantity
FROM items_order
GROUP BY product_id
LIMIT 10;

-- name: DeleteItemsOrder :exec
DELETE FROM items_order WHERE id = $1;

-- name: DeleteItemsOrderByBookingID :exec
DELETE FROM items_order WHERE booking_id = $1;