-- name: CreateRole :one
INSERT INTO roles (
  name
) VALUES (
  $1
) RETURNING *;

-- name: GetRole :one
SELECT * FROM roles
WHERE name = $1 LIMIT 1;