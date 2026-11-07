INSERT INTO teams (id, name, short_name, logo_url) VALUES 
('team-1', 'Arsenal', 'ARS', NULL),
('team-2', 'Chelsea', 'CHE', NULL)
ON CONFLICT (id) DO NOTHING;

INSERT INTO matches (id, home_team_id, away_team_id, home_score, away_score, status, minute, start_time) VALUES 
('match-1', 'team-1', 'team-2', 2, 1, 'LIVE', 75, '2026-11-07 20:00:00+00')
ON CONFLICT (id) DO NOTHING;