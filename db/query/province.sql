-- name: CreateProvince :one
INSERT INTO provinces (
  name
) VALUES (
  $1
) RETURNING *;

-- name: GetProvince :one
SELECT * FROM provinces
WHERE name = $1 LIMIT 1;

-- name: GetProvinceByID :one
SELECT * FROM provinces
WHERE id = $1 LIMIT 1;

-- name: ListProvinces :many
SELECT name FROM provinces
ORDER BY name;