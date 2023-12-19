-- name: CreateProduct :one
INSERT INTO products (
  product_name,
  thumb,
  price
) VALUES (
  $1, $2, $3
) RETURNING *;

-- name: GetProduct :one
SELECT * FROM products
WHERE id = $1 LIMIT 1;

-- name: FindProduct :many
SELECT products.*
FROM products
INNER JOIN descriptions_product
ON products.id = descriptions_product.product_id
WHERE (
  product_name ILIKE '%' || $1 || '%'
  OR descriptions_product.gender ILIKE '%' || $1 || '%'
  OR descriptions_product.material ILIKE '%' || $1 || '%'
)
LIMIT $2
OFFSET $3;

-- name: ListProducts :many
SELECT * FROM products
ORDER BY product_name
LIMIT $1
OFFSET $2;

-- name: ListProductsNoLimit :many
SELECT * FROM products
ORDER BY product_name;

-- name: UpdateProduct :one
UPDATE products
SET product_name = $2, thumb = $3, price = $4
WHERE id = $1
RETURNING *;

-- name: DeleteProduct :exec
DELETE FROM products WHERE id = $1;