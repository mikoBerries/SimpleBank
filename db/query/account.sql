-- name: CreateAccount :one
INSERT INTO accounts (
  owner,
  balance,
  currency
) VALUES (
  $1, $2, $3
)RETURNING *;

-- name: GetAccount :one
SELECT * FROM accounts
WHERE id = $1;

-- name: ListAccount :many
SELECT * FROM accounts
ORDER BY id;

-- name: UpdateAccounts :exec
UPDATE accounts 
SET owner = $2, balance = $3,currency = $4
WHERE id = $1;

-- name: DeleteAccount :exec
DELETE FROM accounts WHERE id = $1;