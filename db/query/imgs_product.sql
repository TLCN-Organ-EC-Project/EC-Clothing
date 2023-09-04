-- name: CreateImgProduct :one
INSERT INTO imgs_product (
  product_id,
  image
) VALUES (
  $1, $2
) RETURNING *;

-- name: GetImgProduct :one
SELECT * FROM imgs_product
WHERE id = $1 LIMIT 1;

-- name: ListImgProducts :many
SELECT * FROM imgs_product
WHERE product_id = $1
ORDER BY id
LIMIT $2
OFFSET $3;

-- name: UpdateImgProduct :one
UPDATE imgs_product
SET product_id = $2, image = $3
WHERE id = $1
RETURNING *;

-- name: DeleteImgProduct :exec
DELETE FROM imgs_product WHERE id = $1;