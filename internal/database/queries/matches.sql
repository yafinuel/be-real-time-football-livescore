-- name: GetMatchByID :one
SELECT id, home_team_id, away_team_id, home_score, away_score, status, minute, start_time
FROM matches
WHERE id = $1 LIMIT 1;

-- name: ListMatches :many
SELECT id, home_team_id, away_team_id, home_score, away_score, status, minute, start_time
FROM matches
ORDER BY start_time DESC;

-- name: UpdateMatchScore :one
UPDATE matches
SET home_score = $2, away_score = $3
WHERE id = $1
RETURNING id, home_team_id, away_team_id, home_score, away_score, status, minute, start_time;