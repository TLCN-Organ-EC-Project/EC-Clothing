-- name: CreateDescriptionProduct :one
INSERT INTO descriptions_product (
  product_id,
  gender,
  material,
  size,
  size_of_model
) VALUES (
  $1, $2, $3, $4, $5
) RETURNING *;

-- name: GetDescriptionProductByID :one
SELECT * FROM descriptions_product
WHERE product_id = $1 LIMIT 1;

-- name: ListDescriptionProduct :many
SELECT * FROM descriptions_product
ORDER BY product_id
LIMIT $1
OFFSET $2;

-- name: UpdateDescriptionProduct :one
UPDATE descriptions_product
SET gender = $2, material = $3, size = $4, size_of_model = $5
WHERE product_id = $1
RETURNING *;

-- name: DeleteDescriptionProduct :exec
DELETE FROM descriptions_product WHERE product_id = $1;