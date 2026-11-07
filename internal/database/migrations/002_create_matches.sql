CREATE TABLE IF NOT EXISTS matches (
    id VARCHAR(64) PRIMARY KEY,
    home_team_id VARCHAR(64) NOT NULL REFERENCES teams(id),
    away_team_id VARCHAR(64) NOT NULL REFERENCES teams(id),
    home_score INTEGER NOT NULL DEFAULT 0,
    away_score INTEGER NOT NULL DEFAULT 0,
    status VARCHAR(32) NOT NULL DEFAULT 'NOT_STARTED',
    minute INTEGER NOT NULL DEFAULT 0,
    start_time TIMESTAMPTZ NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_matches_status ON matches(status);