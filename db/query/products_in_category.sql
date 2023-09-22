-- name: CreateProductsInCategory :one
INSERT INTO products_in_category (
  product_id,
  category_id
) VALUES (
  $1, $2
) RETURNING *;

-- name: GetProductsInCategoryByID :one
SELECT * FROM products_in_category
WHERE category_id = $1 AND product_id = $2 LIMIT 1;

-- name: ListProductsInCategory :many
SELECT * FROM products_in_category
WHERE category_id = $1
ORDER BY product_id
LIMIT $2
OFFSET $3;

-- name: UpdateProductsInCategory :one
UPDATE products_in_category
SET product_id = $2, category_id = $3
WHERE id = $1
RETURNING *;

-- name: DeleteProductsInCategory :exec
DELETE FROM products_in_category WHERE id = $1;