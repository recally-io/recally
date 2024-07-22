-- name: GetCacheByKey :one
SELECT * FROM cache WHERE domain= $1 AND key = $2 AND expires_at > now();

-- name: CreateCache :exec
INSERT INTO cache (domain, key, value, expires_at)
VALUES (@domain, @key, @value, @expires_at);

-- name: UpdateCache :exec
UPDATE cache SET value = @value, expires_at = @expires_at
WHERE key = $1 AND domain = $2;

-- name: DeleteCacheByKey :exec
DELETE FROM cache WHERE key = $1 AND domain = $2;

-- name: DeleteExpiredCache :exec
DELETE FROM cache WHERE expires_at < @expires_at;

-- name: IsCacheExists :one
SELECT EXISTS(SELECT 1 FROM cache WHERE domain = $1 AND key = $2 AND expires_at > now()) as exists;
