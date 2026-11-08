-- name: GetMatchByID :one
SELECT id, home_team_id, away_team_id, home_score, away_score, status, minute, start_time
FROM matches
WHERE id = $1 LIMIT 1;

-- name: ListMatches :many
SELECT id, home_team_id, away_team_id, home_score, away_score, status, minute, start_time
FROM matches
ORDER BY 
    CASE 
        WHEN status = 'LIVE' THEN 1 
        WHEN status = 'NOT_STARTED' THEN 2 
        ELSE 3 
    END, 
    start_time DESC;

-- name: CreateMatch :one
INSERT INTO matches (id, home_team_id, away_team_id, home_score, away_score, status, minute, start_time)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING id, home_team_id, away_team_id, home_score, away_score, status, minute, start_time;

-- name: UpdateMatchScore :one
UPDATE matches
SET home_score = $2, away_score = $3
WHERE id = $1
RETURNING id, home_team_id, away_team_id, home_score, away_score, status, minute, start_time;

-- name: UpdateMatchStatus :one
UPDATE matches
SET status = $2, minute = COALESCE(sqlc.narg('minute'), minute)
WHERE id = $1
RETURNING id, home_team_id, away_team_id, home_score, away_score, status, minute, start_time;
