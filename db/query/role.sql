-- name: GetRole :one
SELECT * FROM roles
WHERE name = $1 LIMIT 1;