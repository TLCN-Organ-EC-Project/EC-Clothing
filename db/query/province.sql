-- name: GetProvince :one
SELECT * FROM provinces
WHERE name = $1 LIMIT 1;