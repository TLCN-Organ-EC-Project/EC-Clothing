-- name: CreateCart :one
INSERT INTO carts (
  username,
  product_id,
  quantity,
  price,
  size
) VALUES (
  $1, $2, $3, $4, $5
) RETURNING *;

-- name: GetCart :one
SELECT * FROM carts
WHERE id = $1 LIMIT 1;

-- name: GetCartDetails :one
SELECT * FROM carts
WHERE username = $1 AND product_id = $2 AND size = $3
LIMIT 1;

-- name: ListCartOfUser :many
SELECT * FROM carts
WHERE username = $1
LIMIT $2
OFFSET $3;

-- name: UpdateCart :one
UPDATE carts
SET quantity = $2, size = $3, price = $4
WHERE id = $1
RETURNING *;

-- name: DeleteCart :exec
DELETE FROM carts WHERE id = $1;