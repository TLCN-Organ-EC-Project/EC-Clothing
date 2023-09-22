-- name: GetProvince :one
SELECT * FROM provinces
WHERE name = $1 LIMIT 1;

-- name: ListProvinces :many
SELECT name FROM provinces
ORDER BY name;