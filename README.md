# ⚽ Real-Time Football Live Score Backend

Backend API untuk aplikasi skor sepak bola real-time, dibangun dengan **Go**, **GraphQL (gqlgen)**, dan **PostgreSQL 17**. Dilengkapi dengan In-Memory PubSub Broker, simulasi pertandingan otomatis, dan koneksi real-time via WebSocket.

Proyek ini siap digunakan oleh **Frontend Developer** yang ingin latihan integrasi GraphQL & WebSocket tanpa perlu mengoding backend sendiri.

---

## 🚀 Quick Start untuk Frontend Developer (Tanpa Perlu Go!)

```bash
docker compose -f compose.fe.yaml up -d
```

Setelah container berjalan:
- 🎮 **GraphQL Playground**: [http://localhost:8080/](http://localhost:8080/)
- 🌐 **HTTP Endpoint (Query & Mutation)**: `http://localhost:8080/query`
- ⚡ **WebSocket Endpoint (Subscription)**: `ws://localhost:8080/query`
- 🩺 **Health Check**: `http://localhost:8080/healthy`

```bash
docker compose -f compose.fe.yaml down
```

---

## 🛠️ Quick Start untuk Backend Developer

```bash
# Jalankan database dan API dalam mode development (live reload)
docker compose up -d

# Menjalankan seluruh test (Unit & Integration)
docker compose run --rm api go test ./... -v

# Menjalankan linter
docker compose run --rm api go vet ./...
```

---

## 📡 Panduan GraphQL API

CORS sudah diaktifkan secara global, sehingga kamu bisa memanggil API ini dari `localhost:3000`, `localhost:5173`, atau port manapun.

### 1. Queries (Mengambil Data)

#### Semua Pertandingan
```graphql
query GetMatches {
  matches {
    id
    minute
    status
    startTime
    score {
      home
      away
    }
    homeTeam {
      id
      name
      shortName
      logoUrl
    }
    awayTeam {
      id
      name
      shortName
      logoUrl
    }
  }
}
```

> Urutan hasil: `LIVE` → `NOT_STARTED` → `FINISHED`

#### Pertandingan yang Sedang LIVE
```graphql
query GetLiveMatches {
  liveMatches {
    id
    minute
    score {
      home
      away
    }
    homeTeam { name }
    awayTeam { name }
  }
}
```

#### Detail 1 Pertandingan Beserta Timeline Events
```graphql
query GetMatchDetail {
  match(id: "match-1") {
    id
    minute
    status
    score {
      home
      away
    }
    events {
      id
      minute
      type
      player
      detail
    }
  }
}
```

#### Semua Tim
```graphql
query GetTeams {
  teams {
    id
    name
    shortName
    logoUrl
  }
}
```

---

### 2. Mutations (Simulasi Manual / Admin)

#### Membuat Pertandingan Baru
```graphql
mutation CreateMatch {
  createMatch(input: {
    homeTeamId: "team-ars"
    awayTeamId: "team-che"
    startTime: "2026-12-01T20:00:00Z"
  }) {
    id
    status
    startTime
    homeTeam { name }
    awayTeam { name }
  }
}
```

> `startTime` menggunakan format **ISO-8601 / RFC3339** (`2006-01-02T15:04:05Z`).

#### Update Skor Pertandingan
```graphql
mutation UpdateScore {
  updateScore(input: {
    matchId: "match-1"
    homeScore: 3
    awayScore: 2
  }) {
    id
    score {
      home
      away
    }
  }
}
```

#### Tambah Event (Gol / Kartu)
```graphql
mutation AddGoalEvent {
  addMatchEvent(
    matchId: "match-1"
    input: {
      minute: 70
      type: GOAL
      teamId: "team-1"
      player: "Bukayo Saka"
      detail: "Tembakan melengkung kaki kiri"
    }
  ) {
    id
    minute
    player
  }
}
```

> `type` yang valid: `GOAL`, `YELLOW_CARD`, `RED_CARD`, `SUBSTITUTION`

#### Ubah Status Pertandingan
```graphql
mutation StartMatch {
  updateMatchStatus(matchId: "match-3", status: LIVE, minute: 1) {
    id
    status
    minute
  }
}
```

> `status` yang valid: `NOT_STARTED`, `LIVE`, `HALFTIME`, `FINISHED`, `POSTPONED`, `CANCELLED`

---

### 3. Subscriptions (Real-Time via WebSocket)

#### Update Skor & Menit Pertandingan Tertentu
```graphql
subscription OnMatchUpdated {
  matchUpdated(matchId: "match-1") {
    id
    minute
    status
    score {
      home
      away
    }
  }
}
```

> Kosongkan `matchId` untuk mendengarkan update dari **semua** pertandingan sekaligus.

#### Event Baru (Gol, Kartu)
```graphql
subscription OnNewEvent {
  newEventAdded(matchId: "match-1") {
    id
    minute
    type
    player
    detail
  }
}
```

---

## 💻 Contoh Integrasi di Frontend (React + Apollo Client)

Backend mendukung protokol modern **`graphql-ws`**:

```bash
npm install @apollo/client graphql graphql-ws
```

```javascript
import { ApolloClient, InMemoryCache, split, HttpLink } from '@apollo/client';
import { GraphQLWsLink } from '@apollo/client/link/subscriptions';
import { createClient } from 'graphql-ws';
import { getMainDefinition } from '@apollo/client/utilities';

const httpLink = new HttpLink({
  uri: 'http://localhost:8080/query',
});

const wsLink = new GraphQLWsLink(
  createClient({
    url: 'ws://localhost:8080/query',
  })
);

const splitLink = split(
  ({ query }) => {
    const definition = getMainDefinition(query);
    return (
      definition.kind === 'OperationDefinition' &&
      definition.operation === 'subscription'
    );
  },
  wsLink,
  httpLink
);

export const client = new ApolloClient({
  link: splitLink,
  cache: new InMemoryCache(),
});
```

---

## 🤖 Mock Match Engine (Simulator Otomatis)

Background worker berjalan otomatis setiap **45 detik** (dapat dikonfigurasi via env var `SIMULATOR_INTERVAL`):
- Mencari pertandingan berstatus **`LIVE`**
- Menambah menit pertandingan secara bertahap
- Memiliki peluang acak mencetak gol (15% per tick)
- Menyiarkan perubahan secara instan via WebSocket ke semua subscriber
- Pertandingan mencapai menit 90 → otomatis `FINISHED`

Untuk mengubah kecepatan simulasi, ubah nilai di `.env`:
```env
SIMULATOR_INTERVAL=30s   # 30 detik
SIMULATOR_INTERVAL=1m    # 1 menit
```

---

## 🏛️ Arsitektur & Tech Stack

```
backend/
├── cmd/api/            # Entry point HTTP & WebSocket server
├── internal/
│   ├── config/         # Environment & konfigurasi
│   ├── database/       # Migrations SQL & sqlc generator
│   ├── domain/         # Pure business domain entities
│   ├── engine/         # Mock Match Simulator (goroutine background)
│   ├── graph/          # GraphQL schema, resolvers, generated code (gqlgen)
│   ├── middleware/     # CORS & preflight handler
│   ├── pubsub/         # In-memory broker (Go channels & RWMutex)
│   ├── repository/     # Data access layer (PostgreSQL pgx/v5)
│   └── service/        # Business logic & validasi aturan
├── tests/
│   ├── integration/    # GraphQL HTTP & subscription integration tests
│   └── testhelper/     # In-memory mock repositories
└── .github/workflows/  # CI lint/test & CD push to GHCR
```

---

## 📜 Lisensi
MIT License. Bebas digunakan untuk bahan latihan, portofolio, dan kolaborasi!
