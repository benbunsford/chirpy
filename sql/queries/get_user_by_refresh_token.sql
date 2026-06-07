-- name: GetUserByRefreshToken :one
SELECT users.* FROM users
JOIN refresh_tokens ON users.id = refresh_tokens.user_id
WHERE token = $1
AND expires_at > NOW()
AND revoked_at IS NULL;
