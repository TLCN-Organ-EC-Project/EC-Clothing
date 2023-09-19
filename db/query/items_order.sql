-- name: CreateItemsOrder :one
INSERT INTO items_order (
  booking_id,
  product_id,
  quantity,
  price
) VALUES (
  $1, $2, $3, $4
) RETURNING *;

-- name: GetItemsOrder :one
SELECT * FROM items_order
WHERE id = $1 LIMIT 1;

-- name: ListItemsOrderByBookingID :many
SELECT * FROM items_order
WHERE booking_id = $1
ORDER BY id
LIMIT $2
OFFSET $3;

-- name: DeleteItemsOrder :exec
DELETE FROM items_order WHERE id = $1;