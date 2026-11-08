-- ─────────────────────────────────────────────
-- 1. Seed Teams
-- ─────────────────────────────────────────────
INSERT INTO teams (id, name, short_name, logo_url) VALUES 
('team-1', 'Arsenal', 'ARS', 'https://resources.premierleague.com/premierleague/badges/50/t3.png'),
('team-2', 'Chelsea', 'CHE', 'https://resources.premierleague.com/premierleague/badges/50/t8.png'),
('team-3', 'Liverpool', 'LIV', 'https://resources.premierleague.com/premierleague/badges/50/t14.png'),
('team-4', 'Manchester City', 'MCI', 'https://resources.premierleague.com/premierleague/badges/50/t43.png'),
('team-5', 'Real Madrid', 'RMA', 'https://upload.wikimedia.org/wikipedia/en/5/56/Real_Madrid_CF.svg'),
('team-6', 'Barcelona', 'FCB', 'https://upload.wikimedia.org/wikipedia/en/4/47/FC_Barcelona_%28crest%29.svg'),
('team-7', 'Bayern Munich', 'BAY', 'https://upload.wikimedia.org/wikipedia/commons/1/1b/FC_Bayern_M%C3%BCnchen_logo_%282017%29.svg'),
('team-8', 'Borussia Dortmund', 'BVB', 'https://upload.wikimedia.org/wikipedia/commons/6/67/Borussia_Dortmund_logo.svg')
ON CONFLICT (id) DO NOTHING;

-- ─────────────────────────────────────────────
-- 2. Seed Matches (Variasi Status: LIVE, NOT_STARTED, FINISHED)
-- ─────────────────────────────────────────────
INSERT INTO matches (id, home_team_id, away_team_id, home_score, away_score, status, minute, start_time) VALUES 
-- Match 1: Premier League Derby (LIVE)
('match-1', 'team-1', 'team-2', 2, 1, 'LIVE', 65, NOW() - INTERVAL '65 minutes'),

-- Match 2: Big Match (LIVE)
('match-2', 'team-3', 'team-4', 1, 1, 'LIVE', 32, NOW() - INTERVAL '32 minutes'),

-- Match 3: El Clasico (NOT_STARTED / Jadwal Nanti)
('match-3', 'team-5', 'team-6', 0, 0, 'NOT_STARTED', 0, NOW() + INTERVAL '2 hours'),

-- Match 4: Der Klassiker (FINISHED / Pertandingan Selesai)
('match-4', 'team-7', 'team-8', 3, 2, 'FINISHED', 90, NOW() - INTERVAL '3 hours')
ON CONFLICT (id) DO NOTHING;

-- ─────────────────────────────────────────────
-- 3. Seed Match Events
-- ─────────────────────────────────────────────
INSERT INTO match_events (id, match_id, minute, type, team_id, player, detail) VALUES 
-- Events untuk match-1 (Arsenal vs Chelsea)
('event-1', 'match-1', 14, 'GOAL', 'team-1', 'Bukayo Saka', 'Assist dari Martin Odegaard'),
('event-2', 'match-1', 28, 'YELLOW_CARD', 'team-2', 'Moises Caicedo', 'Pelanggaran keras di lini tengah'),
('event-3', 'match-1', 41, 'GOAL', 'team-2', 'Cole Palmer', 'Tendangan penalti akurat'),
('event-4', 'match-1', 58, 'GOAL', 'team-1', 'Kai Havertz', 'Sundulan menyambut sepak pojok'),

-- Events untuk match-4 (Bayern vs Dortmund - Finished)
('event-5', 'match-4', 12, 'GOAL', 'team-7', 'Harry Kane', 'Tembakan mendatar ke sudut kiri gawang'),
('event-6', 'match-4', 35, 'GOAL', 'team-8', 'Julian Brandt', 'Gol dari luar kotak penalti'),
('event-7', 'match-4', 62, 'GOAL', 'team-7', 'Jamal Musiala', 'Aksi solo melewati dua bek lawan'),
('event-8', 'match-4', 77, 'GOAL', 'team-8', 'Serhou Guirassy', 'Rebound dari tendangan sudut'),
('event-9', 'match-4', 88, 'GOAL', 'team-7', 'Leroy Sane', 'Gol penentu kemenangan')
ON CONFLICT (id) DO NOTHING;
