CREATE TABLE IF NOT EXISTS match_events (
    id VARCHAR(64) PRIMARY KEY,
    match_id VARCHAR(64) NOT NULL REFERENCES matches(id) ON DELETE CASCADE,
    minute INTEGER NOT NULL,
    type VARCHAR(32) NOT NULL,
    team_id VARCHAR(64) NOT NULL REFERENCES teams(id),
    player VARCHAR(255) NOT NULL,
    detail TEXT
);

CREATE INDEX IF NOT EXISTS idx_match_events_match_id ON match_events(match_id);