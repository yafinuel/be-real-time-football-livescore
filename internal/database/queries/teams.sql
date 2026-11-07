-- name: GetTeamByID :one
SELECT id, name, short_name, logo_url
FROM teams
WHERE id = $1 LIMIT 1;

-- name: ListTeams :many
SELECT id, name, short_name, logo_url
FROM teams
ORDER BY name ASC;