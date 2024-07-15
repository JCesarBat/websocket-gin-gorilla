

-- name: CreateUser :one
INSERT INTO "users" (username, email, password)
VALUES($1,$2,$3) RETURNING *;

-- name: GetUserEmail :one
SELECT * FROM "users" WHERE "email" =$1 ;