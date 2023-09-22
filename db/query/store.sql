-- name: CreateStore :one
INSERT INTO store (
  product_id,
  size,
  quantity
) VALUES (
  $1, $2, $3
) RETURNING *;

-- name: GetStore :one
SELECT * FROM store
WHERE product_id = $1 AND size = $2
LIMIT 1;

-- name: ListStore :many
SELECT * FROM store
WHERE product_id = $1
ORDER BY id;

-- name: UpdateStore :one
UPDATE store
SET quantity = $3
WHERE product_id = $1 AND size = $2
RETURNING *;

-- name: DeleteStore :exec
DELETE FROM store WHERE product_id = $1 AND size = $2;