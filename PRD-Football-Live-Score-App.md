# Product Requirements Document (PRD)

## Real-Time Football Live Score App

| Atribut | Detail |
|---|---|
| **Nama Proyek** | Real-Time Football Live Score App |
| **Versi Dokumen** | 1.0 (Draft) |
| **Tanggal** | 7 November 2026 |
| **Jenis Proyek** | Latihan / Portofolio |
| **Status** | Siap untuk implementasi |

---

## Daftar Isi

1. [Ringkasan Proyek & Tujuan](#1-ringkasan-proyek--tujuan)
2. [Ruang Lingkup](#2-ruang-lingkup)
3. [Tech Stack & Arsitektur](#3-tech-stack--arsitektur)
4. [Fitur Utama & Spesifikasi Teknis](#4-fitur-utama--spesifikasi-teknis)
5. [Struktur GraphQL Schema](#5-struktur-graphql-schema)
6. [Model Data & Database](#6-model-data--database)
7. [Mock Data Generator (Event Engine)](#7-mock-data-generator-event-engine)
8. [Persyaratan Non-Fungsional](#8-persyaratan-non-fungsional)
9. [Struktur Proyek](#9-struktur-proyek)
10. [Rencana Pengujian](#10-rencana-pengujian)
11. [Milestone & Roadmap](#11-milestone--roadmap)
12. [Risiko & Mitigasi](#12-risiko--mitigasi)
13. [Metrik Keberhasilan](#13-metrik-keberhasilan)
14. [Pertanyaan Terbuka & Pengembangan Lanjutan](#14-pertanyaan-terbuka--pengembangan-lanjutan)

---

## 1. Ringkasan Proyek & Tujuan

### 1.1 Deskripsi

Aplikasi **Real-Time Football Live Score** menampilkan skor pertandingan sepak bola secara real-time, jadwal pertandingan, serta detail pertandingan. Komunikasi data dirancang efisien dengan latency rendah menggunakan GraphQL (Query, Mutation, Subscription) di atas backend Go.

### 1.2 Tujuan Produk

- Menampilkan daftar pertandingan (Live, Finished, Scheduled) dengan akurat.
- Mendorong (*push*) pembaruan skor dan menit pertandingan ke klien secara otomatis tanpa *polling*.
- Menyediakan antarmuka admin untuk mengubah skor dan status pertandingan secara manual.

### 1.3 Tujuan Utama Pengembang

Proyek ini adalah **latihan/portofolio** untuk:

1. Mengimplementasikan **GraphQL API** secara lengkap: Query, Mutation, dan Subscription.
2. Mendalami **Go (Golang)**, terutama konkurensi via goroutine dan channel.
3. Menerapkan pendekatan *schema-first* dan *type-safe* dari API sampai database.

### 1.4 Target Pengguna

| Persona | Kebutuhan |
|---|---|
| **Penggemar sepak bola (Klien)** | Melihat skor live, jadwal, dan detail pertandingan tanpa refresh manual. |
| **Admin** | Memperbarui skor dan status pertandingan secara manual. |
| **Pengembang / Reviewer portofolio** | Melihat kualitas arsitektur, kode, dan dokumentasi. |

---

## 2. Ruang Lingkup

### 2.1 Dalam Lingkup (In Scope)

- Backend GraphQL API dengan Query, Mutation, dan Subscription.
- Pembaruan real-time via WebSocket.
- Persistensi data pertandingan (SQLite, opsional PostgreSQL).
- Mock data generator berbasis goroutine dan channel.
- Dokumentasi API dan panduan menjalankan secara lokal.

### 2.2 Di Luar Lingkup (Out of Scope)

- Integrasi API data sepak bola eksternal berbayar.
- Autentikasi pengguna lengkap (OAuth, manajemen akun). Pengamanan admin pada versi ini dibatasi (lihat 8.3).
- Aplikasi klien produksi (mobile/web). Klien yang dibuat hanya sebatas demo sederhana atau GraphQL Playground.
- Fitur taruhan, statistik lanjutan, dan komentar pertandingan.

---

## 3. Tech Stack & Arsitektur

### 3.1 Tech Stack

| Komponen | Pilihan | Alasan |
|---|---|---|
| **Bahasa Backend** | Go (Golang) | Performa tinggi, efisiensi memori, konkurensi kuat via goroutine/channel. |
| **Protokol API** | GraphQL (MIT, 100% gratis & open source) | Klien meminta tepat data yang dibutuhkan; mendukung subscription. |
| **Framework GraphQL** | `gqlgen` (schema-first) | Kode resolver dan tipe digenerate dari schema, type-safe. |
| **WebSocket** | `gorilla/websocket` | Di-handle oleh *transport layer* `gqlgen`. |
| **Database** | SQLite (lokal) / PostgreSQL (opsional produksi) | SQLite mudah untuk latihan; PostgreSQL untuk skala lebih besar. |
| **Database Tool / ORM** | `sqlc` (type-safe SQL generator) atau `ent` | Query type-safe tanpa *boilerplate*. |
| **Sumber Data / Event Engine** | Mock Data Generator (goroutine + channel) | Simulasi skor/menit otomatis tanpa biaya API eksternal. |

### 3.2 Diagram Arsitektur

```
┌──────────────────┐        HTTP (Query/Mutation)        ┌─────────────────────────┐
│  Klien (Web/App) │ ──────────────────────────────────► │   GraphQL Server (Go)   │
│  / GraphQL       │                                     │   gqlgen + resolvers    │
│  Playground      │ ◄────────────────────────────────── │                         │
└──────────────────┘     WebSocket (Subscription)        │  ┌───────────────────┐  │
                                                         │  │  PubSub / Broker  │  │
                                                         │  └─────────▲─────────┘  │
                                                         └────────────┼────────────┘
                                                                      │ publish
                        ┌───────────────────────┐          ┌──────────┴─────────┐
                        │ Mock Data Generator   │ ───────► │  Match Service     │
                        │ (goroutine + channel) │  events  │  (business logic)  │
                        └───────────────────────┘          └──────────┬─────────┘
                                                                      │
                                                           ┌──────────▼─────────┐
                                                           │ SQLite / PostgreSQL│
                                                           │ (sqlc / ent)       │
                                                           └────────────────────┘
```

### 3.3 Alur Data Utama

1. **Query:** Klien mengirim HTTP POST ke `/query` → resolver memanggil Match Service → data dibaca dari database → hanya field yang diminta dikembalikan.
2. **Mutation:** Admin mengirim `updateScore` → Match Service memvalidasi dan menyimpan ke database → event dipublikasikan ke PubSub.
3. **Subscription:** Klien membuka koneksi WebSocket dan berlangganan `matchUpdated(matchId)` → server mendaftarkan channel per-subscriber → setiap event pada match tersebut didorong ke klien.
4. **Mock Generator:** Goroutine berjalan di latar belakang, mengubah menit/skor pertandingan berstatus `LIVE`, menyimpan perubahan, lalu mempublikasikan event ke PubSub.

### 3.4 Pola Konkurensi

- **PubSub Broker:** menyimpan `map[matchID][]chan *Match` yang dilindungi `sync.RWMutex`.
- **Satu channel per subscriber** (buffered) agar subscriber lambat tidak memblokir publisher.
- **Pembersihan:** saat `context` subscription dibatalkan (klien disconnect), channel ditutup dan dihapus dari broker untuk mencegah *goroutine leak*.
- **Graceful shutdown:** semua goroutine generator dihentikan lewat `context.Context`.

---

## 4. Fitur Utama & Spesifikasi Teknis

### 4.1 Queries (Pengambilan Data)

| ID | Persyaratan | Prioritas |
|---|---|---|
| Q-01 | Mengambil daftar semua pertandingan (Live, Finished, Scheduled) melalui `matches`. | Must |
| Q-02 | Mengambil detail satu pertandingan berdasarkan `matchId` melalui `match(id)`; mengembalikan `null` jika tidak ditemukan. | Must |
| Q-03 | Klien dapat memilih field yang diperlukan sehingga tidak terjadi *over-fetching* maupun *under-fetching*. | Must |
| Q-04 | Daftar pertandingan diurutkan secara konsisten (mis. LIVE lebih dulu, lalu berdasarkan waktu mulai). | Should |

**Contoh Query:**

```graphql
query {
  matches {
    id
    homeTeam
    awayTeam
    homeScore
    awayScore
    status
    minute
  }
}
```

### 4.2 Mutations (Pengubahan Data)

| ID | Persyaratan | Prioritas |
|---|---|---|
| M-01 | Admin dapat memperbarui skor secara manual melalui `updateScore(matchId, homeScore, awayScore)`. | Must |
| M-02 | Admin dapat mengubah status pertandingan: `NOT_STARTED` → `LIVE` → `FINISHED`. | Must |
| M-03 | Validasi: skor tidak boleh negatif; `matchId` harus ada; transisi status harus valid (tidak boleh mundur). | Must |
| M-04 | Setiap mutation yang berhasil memicu event ke subscriber yang relevan. | Must |

**Contoh Mutation:**

```graphql
mutation {
  updateScore(matchId: "1", homeScore: 2, awayScore: 1) {
    id
    homeScore
    awayScore
  }
}
```

> **Catatan:** Schema pada Bagian 5 hanya mendefinisikan `updateScore`. Perubahan status (M-02) dapat diimplementasikan sebagai mutation tambahan `updateStatus(matchId: ID!, status: String!): Match!` pada iterasi berikutnya, atau dijalankan lewat mock generator. Lihat Bagian 14.

### 4.3 Subscriptions (Real-Time via WebSocket)

| ID | Persyaratan | Prioritas |
|---|---|---|
| S-01 | Klien dapat berlangganan `matchUpdated(matchId)` untuk satu pertandingan. | Must |
| S-02 | Server mendorong data otomatis setiap ada **gol** atau **perubahan menit**. | Must |
| S-03 | Koneksi dibersihkan otomatis ketika klien terputus (tanpa goroutine leak). | Must |
| S-04 | Mendukung banyak subscriber bersamaan pada pertandingan yang sama. | Must |
| S-05 | Payload event berisi objek `Match` lengkap sesuai field yang diminta klien. | Must |

**Contoh Subscription:**

```graphql
subscription {
  matchUpdated(matchId: "1") {
    id
    homeScore
    awayScore
    status
    minute
  }
}
```

---

## 5. Struktur GraphQL Schema

Schema Definition Language (SDL) yang disepakati:

```graphql
type Match {
  id: ID!
  homeTeam: String!
  awayTeam: String!
  homeScore: Int!
  awayScore: Int!
  status: String!
  minute: Int!
}

type Query {
  matches: [Match!]!
  match(id: ID!): Match
}

type Mutation {
  updateScore(matchId: ID!, homeScore: Int!, awayScore: Int!): Match!
}

type Subscription {
  matchUpdated(matchId: ID!): Match!
}
```

### 5.1 Nilai yang Diizinkan untuk `status`

Field `status` bertipe `String!`. Nilai yang valid dan divalidasi di sisi server:

| Nilai | Arti |
|---|---|
| `NOT_STARTED` | Pertandingan terjadwal, belum dimulai. |
| `LIVE` | Pertandingan sedang berlangsung. |
| `FINISHED` | Pertandingan selesai. |

### 5.2 Konfigurasi `gqlgen` (contoh `gqlgen.yml`)

```yaml
schema:
  - graph/schema.graphqls
exec:
  filename: graph/generated.go
model:
  filename: graph/model/models_gen.go
resolver:
  layout: follow-schema
  dir: graph
autobind:
  - github.com/<username>/football-live-score/graph/model
```

---

## 6. Model Data & Database

### 6.1 Tabel `matches`

| Kolom | Tipe | Keterangan |
|---|---|---|
| `id` | INTEGER / UUID (PK) | Identitas pertandingan. |
| `home_team` | TEXT NOT NULL | Nama tim tuan rumah. |
| `away_team` | TEXT NOT NULL | Nama tim tamu. |
| `home_score` | INTEGER NOT NULL DEFAULT 0 | Skor tuan rumah. |
| `away_score` | INTEGER NOT NULL DEFAULT 0 | Skor tamu. |
| `status` | TEXT NOT NULL | `NOT_STARTED` / `LIVE` / `FINISHED`. |
| `minute` | INTEGER NOT NULL DEFAULT 0 | Menit berjalan (0–90+). |
| `kickoff_at` | DATETIME | Waktu mulai (untuk jadwal). |
| `updated_at` | DATETIME | Waktu pembaruan terakhir. |

### 6.2 Contoh SQL (`sqlc`)

```sql
-- name: ListMatches :many
SELECT * FROM matches
ORDER BY (status = 'LIVE') DESC, kickoff_at ASC;

-- name: GetMatch :one
SELECT * FROM matches WHERE id = ?;

-- name: UpdateScore :one
UPDATE matches
SET home_score = ?, away_score = ?, updated_at = CURRENT_TIMESTAMP
WHERE id = ?
RETURNING *;

-- name: UpdateStatus :one
UPDATE matches
SET status = ?, updated_at = CURRENT_TIMESTAMP
WHERE id = ?
RETURNING *;

-- name: UpdateMinute :one
UPDATE matches
SET minute = ?, updated_at = CURRENT_TIMESTAMP
WHERE id = ?
RETURNING *;
```

### 6.3 Strategi Database

- **Pengembangan lokal:** SQLite (file tunggal, tanpa instalasi server).
- **Produksi (opsional):** PostgreSQL; perbedaan dialek SQL dikelola melalui konfigurasi `sqlc` per engine.
- **Migrasi:** menggunakan file SQL berversi (mis. `golang-migrate`).
- **Seed data:** skrip untuk mengisi beberapa pertandingan awal dalam tiga status.

---

## 7. Mock Data Generator (Event Engine)

### 7.1 Tujuan

Mensimulasikan pertandingan live tanpa API eksternal sehingga fitur subscription dapat didemonstrasikan secara otomatis.

### 7.2 Perilaku

| Parameter | Nilai Default | Konfigurasi |
|---|---|---|
| Interval tick | 1 detik nyata = 1 menit pertandingan (dipercepat) | `TICK_INTERVAL` |
| Durasi pertandingan | 90 menit + tambahan waktu opsional | `MATCH_DURATION` |
| Probabilitas gol per tick | ~2% | `GOAL_PROBABILITY` |
| Transisi otomatis | `NOT_STARTED` → `LIVE` pada kickoff, `LIVE` → `FINISHED` pada menit akhir | – |

### 7.3 Alur Kerja

1. Satu goroutine per pertandingan `LIVE` (atau satu goroutine dengan `time.Ticker` yang mengiterasi semua match live).
2. Setiap tick: naikkan `minute`, secara acak tentukan apakah terjadi gol dan tim mana yang mencetak.
3. Simpan perubahan ke database.
4. Publikasikan event `*Match` ke PubSub broker.
5. Hentikan saat `context` dibatalkan atau pertandingan `FINISHED`.

### 7.4 Contoh Kerangka Kode

```go
func (g *Generator) runMatch(ctx context.Context, m *model.Match) {
    ticker := time.NewTicker(g.tickInterval)
    defer ticker.Stop()

    for {
        select {
        case <-ctx.Done():
            return
        case <-ticker.C:
            m.Minute++
            if rand.Float64() < g.goalProbability {
                if rand.Intn(2) == 0 {
                    m.HomeScore++
                } else {
                    m.AwayScore++
                }
            }
            if m.Minute >= g.matchDuration {
                m.Status = "FINISHED"
            }
            g.store.Save(ctx, m)
            g.broker.Publish(m.ID, m)
            if m.Status == "FINISHED" {
                return
            }
        }
    }
}
```

---

## 8. Persyaratan Non-Fungsional

### 8.1 Performa

| Metrik | Target |
|---|---|
| Latency Query (p95, lokal) | < 50 ms |
| Latency propagasi event ke subscriber (mutation → klien) | < 200 ms |
| Jumlah subscriber bersamaan | ≥ 1.000 koneksi pada satu instance (target latihan) |
| Penggunaan memori idle | Rendah (target < 50 MB) |

### 8.2 Keandalan

- Subscriber yang lambat tidak boleh memblokir publisher (channel ber-buffer; event lama dapat dibuang/di-*drop* bila penuh).
- Tidak ada *goroutine leak* setelah klien terputus (diverifikasi dengan pengujian).
- *Graceful shutdown* menutup semua koneksi WebSocket dan menghentikan generator.
- Penggunaan `sync.RWMutex` atau pola setara untuk mencegah *race condition* (diverifikasi dengan `go test -race`).

### 8.3 Keamanan

- Mutation bersifat khusus admin. Pada versi latihan, perlindungan dilakukan dengan **API key/token statis** via HTTP header `Authorization` yang diperiksa oleh middleware; implementasi autentikasi penuh berada di luar lingkup.
- Validasi seluruh input mutation (rentang skor, status valid, ID ada).
- Batasi *query depth* dan *complexity* untuk mencegah query berlebihan.
- Atur `CheckOrigin` WebSocket dan konfigurasi CORS sesuai lingkungan.
- Nonaktifkan introspection dan GraphQL Playground pada mode produksi.

### 8.4 Observabilitas

- Logging terstruktur (mis. `log/slog`) untuk event server, mutation, dan koneksi subscription.
- Endpoint `/healthz` untuk pemeriksaan kesehatan.
- (Opsional) Metrik jumlah subscriber aktif dan jumlah event terkirim.

### 8.5 Maintainability

- Kode mengikuti konvensi Go (`gofmt`, `go vet`, `golangci-lint`).
- Pemisahan lapisan: resolver → service → repository.
- Kode hasil generate (`gqlgen`, `sqlc`) tidak diedit manual.

---

## 9. Struktur Proyek

```
football-live-score/
├── cmd/
│   └── server/
│       └── main.go            # entry point, wiring, graceful shutdown
├── graph/
│   ├── schema.graphqls        # SDL
│   ├── schema.resolvers.go    # implementasi resolver
│   ├── resolver.go            # dependency injection resolver
│   ├── generated.go           # (generated) gqlgen
│   └── model/                 # (generated) model
├── internal/
│   ├── broker/                # PubSub broker (channels)
│   ├── generator/             # mock data generator
│   ├── service/               # business logic pertandingan
│   ├── store/                 # akses database (sqlc)
│   │   ├── queries/           # file .sql
│   │   └── db/                # (generated) sqlc
│   └── middleware/            # auth admin, logging, CORS
├── migrations/                # migrasi SQL
├── gqlgen.yml
├── sqlc.yaml
├── go.mod
├── Makefile
└── README.md
```

---

## 10. Rencana Pengujian

| Jenis | Cakupan | Alat |
|---|---|---|
| **Unit test** | Business logic service, validasi status/skor, broker publish/subscribe. | `testing`, `testify` |
| **Integrasi** | Resolver terhadap database SQLite in-memory. | `testing`, `httptest` |
| **Real-time / WebSocket** | Klien subscribe → mutation dijalankan → event diterima; disconnect membersihkan channel. | `gqlgen` client / `gorilla/websocket` |
| **Race detection** | Seluruh paket konkuren. | `go test -race ./...` |
| **Beban (opsional)** | Ratusan sampai ribuan subscriber simultan. | `k6` / skrip Go kustom |

### 10.1 Kriteria Penerimaan Utama

- [ ] `matches` mengembalikan seluruh pertandingan dengan status `LIVE`, `FINISHED`, dan `NOT_STARTED`.
- [ ] `match(id)` mengembalikan detail yang benar, dan `null` untuk ID tidak ada.
- [ ] `updateScore` memperbarui database dan memicu event subscription.
- [ ] Subscriber menerima update gol dan perubahan menit dari generator tanpa aksi manual.
- [ ] Dua klien atau lebih pada match yang sama menerima event bersamaan.
- [ ] Setelah klien terputus, jumlah goroutine kembali ke baseline.
- [ ] `go test -race ./...` lulus tanpa temuan.

---

## 11. Milestone & Roadmap

| Fase | Deliverable | Estimasi |
|---|---|---|
| **Fase 1 – Fondasi** | Inisialisasi modul Go, setup `gqlgen`, schema, server dasar, skema database + `sqlc`, seed data. | Minggu 1 |
| **Fase 2 – Query & Mutation** | Resolver `matches`, `match`, `updateScore`; validasi; middleware admin. | Minggu 2 |
| **Fase 3 – Subscription** | PubSub broker, resolver `matchUpdated`, transport WebSocket, pembersihan koneksi. | Minggu 3 |
| **Fase 4 – Mock Generator** | Generator goroutine, transisi status otomatis, integrasi dengan broker. | Minggu 4 |
| **Fase 5 – Kualitas & Dokumentasi** | Pengujian (termasuk `-race`), logging, README, demo klien sederhana. | Minggu 5 |

---

## 12. Risiko & Mitigasi

| Risiko | Dampak | Mitigasi |
|---|---|---|
| Goroutine leak pada subscription | Memori naik seiring waktu | Selalu pantau `ctx.Done()`, hapus subscriber dari broker, uji dengan `runtime.NumGoroutine()`. |
| Race condition pada broker/generator | Data korup atau panic | Mutex yang tepat, `go test -race` di CI. |
| Subscriber lambat memblokir publisher | Latensi event meningkat | Channel ber-buffer + kebijakan *drop* saat penuh, `select` dengan `default`. |
| SQLite terbatas pada konkurensi tulis | Error `database is locked` | Aktifkan mode WAL, serialisasi penulisan, opsi migrasi ke PostgreSQL. |
| `status` bertipe String rentan nilai tidak valid | Data tidak konsisten | Validasi di server; pertimbangkan `enum` GraphQL pada iterasi berikut. |
| Cakupan membengkak | Proyek tidak selesai | Pertahankan ruang lingkup Bagian 2; fitur tambahan di Bagian 14. |

---

## 13. Metrik Keberhasilan

**Teknis**

- Seluruh kriteria penerimaan (10.1) terpenuhi.
- Latency propagasi event < 200 ms pada lingkungan lokal.
- Tidak ada temuan race dan tidak ada goroutine leak.

**Pembelajaran / Portofolio**

- Ketiga operasi GraphQL (Query, Mutation, Subscription) berfungsi dan terdokumentasi.
- Pola konkurensi Go (goroutine, channel, context, mutex) diterapkan dan dapat dijelaskan.
- Repositori memiliki README jelas, instruksi menjalankan lokal, dan contoh query.

---

## 14. Pertanyaan Terbuka & Pengembangan Lanjutan

### 14.1 Pertanyaan Terbuka

1. Apakah `status` tetap `String!` atau diubah menjadi `enum MatchStatus { NOT_STARTED LIVE FINISHED }` demi keamanan tipe?
2. Apakah perubahan status akan diekspos lewat mutation terpisah (`updateStatus`) atau hanya lewat generator?
3. Apakah diperlukan subscription untuk **seluruh** pertandingan (mis. `matchesUpdated`) selain per-`matchId`?
4. SQLite saja, atau sejak awal disiapkan dukungan PostgreSQL?
5. `sqlc` atau `ent` sebagai alat database final?

### 14.2 Ide Pengembangan Lanjutan (Backlog)

- Entitas tambahan: `Team`, `League`, `Player`, `Event` (gol, kartu, pergantian pemain).
- Filter dan pagination pada `matches` (berdasarkan liga, tanggal, status).
- *Cursor-based pagination* (Relay connection).
- Subscription global untuk semua pertandingan live.
- Autentikasi JWT dan peran (admin vs publik).
- Integrasi API data sepak bola nyata sebagai pengganti mock generator.
- Klien demo (web React/Next.js atau Flutter) dengan Apollo/urql.
- Deployment (Docker, CI/CD) dan *rate limiting*.
- Dukungan *horizontal scaling* dengan Redis Pub/Sub sebagai broker lintas instance.

---

## Lampiran A – Perintah Pengembangan (Contoh)

```bash
# Inisialisasi
go mod init github.com/<username>/football-live-score
go get github.com/99designs/gqlgen
go run github.com/99designs/gqlgen init

# Generate kode dari schema
go run github.com/99designs/gqlgen generate

# Generate kode SQL type-safe
sqlc generate

# Jalankan server
go run ./cmd/server

# Pengujian
go test -race ./...
```

## Lampiran B – Referensi

- gqlgen: https://gqlgen.com
- GraphQL Specification: https://spec.graphql.org
- sqlc: https://sqlc.dev
- gorilla/websocket: https://github.com/gorilla/websocket
- Go Concurrency Patterns: https://go.dev/blog/pipelines
