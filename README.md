# ⚽ Real-Time Football Live Score Backend

Backend API performa tinggi untuk aplikasi skor sepak bola real-time yang dibangun dengan **Go (Golang)**, **GraphQL (gqlgen)**, dan **PostgreSQL 17**. Dilengkapi dengan **In-Memory PubSub Broker** berbasis Go channels, simulasi pertandingan otomatis (*Mock Engine*), dan koneksi real-time via **WebSocket**.

Proyek ini siap digunakan oleh **Frontend Developer** yang ingin latihan integrasi GraphQL & WebSocket tanpa perlu mengoding backend sendiri.

---

## 🚀 Quick Start untuk Frontend Developer (Tanpa Perlu Go!)

Kamu tidak perlu menginstal Go, compiler, atau mengerti kodingan backend. Cukup jalankan Docker:

```bash
# 1. Jalankan backend dari image resmi GHCR
docker compose -f compose.fe.yaml up -d
```

Setelah container berjalan:
- 🎮 **GraphQL Playground**: [http://localhost:8080/](http://localhost:8080/)
- 🌐 **HTTP Endpoint (Query & Mutation)**: `http://localhost:8080/query`
- ⚡ **WebSocket Endpoint (Subscription)**: `ws://localhost:8080/query`
- 🩺 **Health Check**: `http://localhost:8080/healthy`

Untuk mematikan container:
```bash
docker compose -f compose.fe.yaml down
```

---

## 🛠️ Quick Start untuk Backend Developer

```bash
# 1. Jalankan database dan API dalam mode development (live reload)
docker compose up -d

# 2. Menjalankan seluruh test (Unit & Integration)
docker compose run --rm api go test ./... -v

# 3. Menjalankan linter
docker compose run --rm api go vet ./...
```

---

## 📡 Panduan GraphQL API

CORS sudah diaktifkan secara global (`*`), sehingga kamu bisa menembak API ini dari `localhost:3000`, `localhost:5173`, atau port manapun!

### 1. Queries (Mengambil Data)

#### Mengambil Semua Pertandingan
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

#### Mengambil Pertandingan yang Sedang LIVE Saja
```graphql
query GetLiveMatches {
  liveMatches {
    id
    minute
    score {
      home
      away
    }
    homeTeam {
      name
    }
    awayTeam {
      name
    }
  }
}
```

#### Mengambil Detail 1 Pertandingan Beserta Timeline Events
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

---

### 2. Subscriptions (Mendengarkan Update Real-Time)

#### Mendengarkan Update Skor & Menit Pertandingan Tertentu
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
*(Kosongkan parameter `matchId` untuk mendengarkan update dari **semua** pertandingan sekaligus).*

#### Mendengarkan Event Baru (Gol, Kartu Kuning/Merah)
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

### 3. Mutations (Admin / Simulasi Manual)

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

#### Menambah Event Gol / Kartu
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

#### Mengubah Status Pertandingan
Status yang valid: `NOT_STARTED`, `LIVE`, `HALFTIME`, `FINISHED`, `POSTPONED`, `CANCELLED`.
```graphql
mutation StartMatch {
  updateMatchStatus(matchId: "match-3", status: LIVE, minute: 1) {
    id
    status
    minute
  }
}
```

---

## 💻 Contoh Integrasi di Frontend (React + Apollo Client)

Backend ini mendukung protokol modern **`graphql-ws`**:

```bash
npm install @apollo/client graphql graphql-ws
```

```javascript
import { ApolloClient, InMemoryCache, split, HttpLink } from '@apollo/client';
import { GraphQLWsLink } from '@apollo/client/link/subscriptions';
import { createClient } from 'graphql-ws';
import { getMainDefinition } from '@apollo/client/utilities';

// 1. HTTP Link untuk Query & Mutation
const httpLink = new HttpLink({
  uri: 'http://localhost:8080/query',
});

// 2. WebSocket Link untuk Subscription
const wsLink = new GraphQLWsLink(
  createClient({
    url: 'ws://localhost:8080/query',
  })
);

// 3. Split link berdasarkan tipe operasi
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

Backend memiliki background worker yang otomatis berjalan setiap **5 detik**:
- Mencari pertandingan yang berstatus **`LIVE`**.
- Menambah menit pertandingan secara bertahap.
- Memiliki peluang acak mencetak gol dan membuat event gol baru.
- Menyiarkan perubahan secara instan melalui WebSocket ke semua client yang sedang subscribe.
- Jika pertandingan mencapai menit 90, statusnya otomatis berubah menjadi `FINISHED`.

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
