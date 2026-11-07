-- name: CreateMatchEvent :one
INSERT INTO match_events (id, match_id, minute, type, team_id, player, detail)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING id, match_id, minute, type, team_id, player, detail;

-- name: ListMatchEventsByMatchID :many
SELECT id, match_id, minute, type, team_id, player, detail
FROM match_events
WHERE match_id = $1
ORDER BY minute ASC;