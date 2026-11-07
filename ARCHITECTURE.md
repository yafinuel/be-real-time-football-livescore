# Backend Architecture

## Real-Time Football Live Score App

  -----------------------------------------------------------------------
  Atribut                             Detail
  ----------------------------------- -----------------------------------
  Nama proyek                         Real-Time Football Live Score App

  Komponen                            Backend

  Bahasa                              Go (Golang)

  API                                 GraphQL dengan gqlgen

  Komunikasi real-time                GraphQL Subscription melalui
                                      WebSocket

  Database MVP                        SQLite

  Database access                     sqlc

  Arsitektur                          Modular Monolith dengan pendekatan
                                      Pragmatic Clean Architecture

  Status dokumen                      Rekomendasi arsitektur awal
  -----------------------------------------------------------------------

------------------------------------------------------------------------

## 1. Tujuan Dokumen

Dokumen ini menjelaskan rancangan arsitektur backend untuk aplikasi
**Real-Time Football Live Score App**. Dokumen menjadi acuan
implementasi agar struktur kode, tanggung jawab setiap komponen, alur
data, dan integrasi antarlapisan tetap konsisten.

Dokumen ini berfokus pada keputusan teknis dan batas antar komponen
backend.

## 2. Ruang Lingkup Sistem

Backend bertanggung jawab untuk:

-   Menyediakan data tim dan pertandingan melalui GraphQL.
-   Membuat dan memperbarui pertandingan.
-   Menyimpan skor, status, menit pertandingan, dan event pertandingan.
-   Menyediakan pembaruan pertandingan secara real-time.
-   Menjalankan mock match engine untuk mensimulasikan pertandingan.
-   Mengelola koneksi database dan validasi business rules.

Ruang lingkup dokumen ini terbatas pada arsitektur backend MVP.

## 3. Keputusan Arsitektur

### 3.1 Modular Monolith

Backend dibangun sebagai satu aplikasi Go yang dapat dijalankan dan
di-deploy sebagai satu unit. Kode di dalamnya dibagi menjadi modul
berdasarkan tanggung jawab.

Alasan pemilihan:

-   Tidak memerlukan komunikasi jaringan antarservice.
-   Transaksi database dan business logic lebih sederhana.
-   Batas tanggung jawab tetap jelas sehingga modul dapat dikembangkan
    secara terpisah.

Microservices belum diperlukan untuk kebutuhan deployment MVP.

### 3.2 Pragmatic Clean Architecture

Kode dipisahkan agar aturan bisnis tidak bergantung langsung pada
gqlgen, WebSocket, atau database.

Prinsip dependensi:

-   Domain tidak mengimpor package database atau GraphQL.
-   Service menggunakan interface repository, bukan implementasi
    database secara langsung.
-   GraphQL resolver bertindak sebagai adapter yang meneruskan request
    ke service.
-   Match engine menggunakan service yang sama dengan GraphQL.
-   Detail implementasi dapat diganti tanpa mengubah domain dan business
    rules.

### 3.3 Schema-first GraphQL

Schema GraphQL menjadi kontrak API. gqlgen menghasilkan kode Go
berdasarkan schema tersebut.

Perubahan API dilakukan pada file schema, kemudian jalankan generator
gqlgen. File hasil generate tidak boleh diedit secara manual.

### 3.4 Database access dengan sqlc

SQL ditulis secara eksplisit, lalu sqlc menghasilkan kode Go yang
type-safe. Repository mengadaptasi kode hasil generate ke interface yang
digunakan service.

SQLite digunakan untuk MVP. PostgreSQL dapat dipertimbangkan kemudian,
dengan memperhatikan perbedaan tipe data, migration, dan konfigurasi
driver.

### 3.5 In-memory PubSub

PubSub awal menggunakan Go channels di dalam satu proses backend.
Pendekatan ini cukup untuk satu instance aplikasi dan mock event engine.

Konsekuensinya, event tidak persisten di broker dan tidak otomatis
dibagikan ke beberapa instance backend. Jika aplikasi kelak dijalankan
secara horizontal, pertimbangkan Redis Pub/Sub atau message broker lain.

------------------------------------------------------------------------

## 4. Struktur Direktori

Struktur yang direkomendasikan:

``` text
football-live-score/
├── cmd/
│   └── api/
│       └── main.go
│
├── internal/
│   ├── config/
│   │   └── config.go
│   │
│   ├── domain/
│   │   ├── match.go
│   │   ├── team.go
│   │   ├── match_event.go
│   │   └── errors.go
│   │
│   ├── service/
│   │   ├── match_service.go
│   │   ├── team_service.go
│   │   └── event_service.go
│   │
│   ├── repository/
│   │   ├── interfaces.go
│   │   └── sqlite/
│   │       ├── match_repository.go
│   │       ├── team_repository.go
│   │       └── event_repository.go
│   │
│   ├── database/
│   │   ├── sqlite.go
│   │   ├── migrations/
│   │   │   ├── 001_create_teams.sql
│   │   │   ├── 002_create_matches.sql
│   │   │   └── 003_create_match_events.sql
│   │   └── queries/
│   │       ├── teams.sql
│   │       ├── matches.sql
│   │       └── match_events.sql
│   │
│   ├── graph/
│   │   ├── schema/
│   │   │   ├── match.graphqls
│   │   │   ├── team.graphqls
│   │   │   └── event.graphqls
│   │   ├── resolver.go
│   │   ├── schema.resolvers.go
│   │   ├── generated.go
│   │   └── model/
│   │
│   ├── pubsub/
│   │   └── broker.go
│   │
│   └── engine/
│       ├── simulator.go
│       └── match_runner.go
│
├── graph/
│   └── gqlgen.yml
│
├── scripts/
│   └── seed.go
│
├── tests/
│   ├── integration/
│   └── fixtures/
│
├── .github/
│   └── workflows/
│       └── ci.yml
│
├── go.mod
├── go.sum
├── Makefile
├── .env.example
├── .gitignore
└── README.md
```

> Catatan: Struktur ini adalah target yang direkomendasikan, bukan
> kewajiban untuk membuat seluruh file sekaligus. Folder atau file dapat
> ditambahkan ketika implementasi membutuhkannya. Konfigurasi gqlgen
> harus menunjuk ke lokasi schema, resolver, model, dan generated code
> yang konsisten dengan struktur aktual.

### 4.1 `cmd/api/`

Entry point aplikasi. `main.go` bertanggung jawab untuk:

1.  Membaca konfigurasi.
2.  Membuka koneksi database.
3.  Menjalankan atau memeriksa migration.
4.  Membuat repository.
5.  Membuat service.
6.  Membuat PubSub broker dan match engine.
7.  Menginisialisasi GraphQL server dan HTTP handler.
8.  Menjalankan HTTP server serta menangani graceful shutdown.

Hindari menempatkan business logic di `main.go`. File ini berfungsi
sebagai composition root untuk merakit dependency.

### 4.2 `internal/config/`

Memuat konfigurasi aplikasi, misalnya:

-   `PORT`
-   `DATABASE_URL` atau path database SQLite
-   `SIMULATION_INTERVAL`
-   konfigurasi lingkungan aplikasi

Validasi konfigurasi dilakukan saat startup agar aplikasi gagal lebih
awal jika konfigurasi wajib tidak tersedia.

### 4.3 `internal/domain/`

Berisi model dan aturan inti domain, seperti `Match`, `Team`,
`MatchEvent`, `MatchStatus`, dan `EventType`.

Domain tidak bergantung pada gqlgen, SQL driver, atau framework HTTP.
Tipe domain dapat memiliki method dan validasi yang memang merupakan
aturan bisnis.

Contoh model konseptual:

``` go
type Match struct {
    ID          string
    HomeTeamID  string
    AwayTeamID  string
    HomeScore   int
    AwayScore   int
    Status      MatchStatus
    Minute      int
    KickoffAt   time.Time
    UpdatedAt   time.Time
}
```

Sebaiknya domain menyimpan referensi tim melalui ID. Informasi lengkap
`HomeTeam` dan `AwayTeam` dapat dirakit oleh service atau resolver dari
hasil query repository.

### 4.4 `internal/service/`

Berisi use case dan business logic.

Tanggung jawab contoh:

-   `MatchService`: membuat pertandingan, mengambil pertandingan,
    memperbarui skor, dan mengubah status.
-   `TeamService`: mengambil dan mengelola data tim.
-   `EventService`: menambahkan event pertandingan dan memvalidasi
    event.

Service menerima dependency melalui constructor, misalnya repository dan
publisher interface. Dengan demikian service tidak perlu membuat koneksi
database atau mengelola WebSocket sendiri.

Validasi penting dilakukan di service/domain, bukan hanya di resolver.
Contoh: skor tidak boleh negatif, tim tuan rumah dan tamu tidak boleh
sama, serta transisi status pertandingan harus valid.

### 4.5 `internal/repository/`

Mendefinisikan kontrak penyimpanan data dan implementasinya.

`interfaces.go` berisi interface yang dibutuhkan service. Implementasi
SQLite ditempatkan di `repository/sqlite/` dan menggunakan kode yang
dihasilkan sqlc.

Contoh interface konseptual:

``` go
type MatchRepository interface {
    GetByID(ctx context.Context, id string) (*domain.Match, error)
    List(ctx context.Context, filter MatchFilter) ([]domain.Match, error)
    Create(ctx context.Context, match *domain.Match) error
    Update(ctx context.Context, match *domain.Match) error
}
```

Nama dan signature method dapat disesuaikan ketika query dan kebutuhan
service sudah ditentukan. Operasi yang mengubah beberapa tabel sebaiknya
mendukung transaksi.

### 4.6 `internal/database/`

Mengelola inisialisasi koneksi database, migration, dan file query SQL.

-   `sqlite.go`: membuka koneksi database dan menyiapkan konfigurasi
    koneksi.
-   `migrations/`: perubahan struktur database secara berurutan dan
    berversi.
-   `queries/`: query SQL yang menjadi input sqlc.

Kode hasil generate sqlc dapat ditempatkan di package khusus, misalnya
`internal/database/sqlc/`, sesuai konfigurasi sqlc. Repository
menggunakan package tersebut dan menerjemahkan hasilnya ke domain model.

### 4.7 `internal/graph/`

Berisi GraphQL schema, resolver, model GraphQL, dan kode hasil generate
gqlgen.

Resolver sebaiknya tipis: mengambil input, memanggil service, lalu
memetakan hasil atau error ke bentuk GraphQL. Resolver tidak boleh
menjadi tempat query SQL atau aturan bisnis utama.

Schema dapat dibagi berdasarkan domain, misalnya `match.graphqls`,
`team.graphqls`, dan `event.graphqls`, lalu digabungkan oleh gqlgen.

### 4.8 `internal/pubsub/`

Menyediakan mekanisme publish/subscribe dalam satu proses.

Broker bertanggung jawab untuk:

-   Mendaftarkan subscriber berdasarkan `matchID` atau topic.
-   Mengirim event kepada subscriber yang relevan.
-   Menghapus subscriber saat subscription berakhir.
-   Menangani channel yang sudah ditutup atau client yang lambat.
-   Melindungi struktur data bersama dari race condition.

Gunakan `context.Context` untuk menghentikan subscription ketika client
disconnect. Tentukan kebijakan backpressure, misalnya buffer terbatas
dan perilaku ketika buffer penuh. Jangan membiarkan publisher macet
tanpa batas karena satu subscriber yang lambat.

### 4.9 `internal/engine/`

Menjalankan simulasi pertandingan di background.

-   `simulator.go`: mengelola lifecycle engine dan pertandingan yang
    aktif.
-   `match_runner.go`: menjalankan simulasi satu pertandingan, misalnya
    dengan `time.Ticker`.

Engine tidak membuat aturan bisnis duplikat. Perubahan menit, skor,
status, dan event dilakukan melalui service agar validasi, penyimpanan,
serta publikasi event tetap konsisten.

Setiap goroutine harus memiliki mekanisme penghentian melalui context.
Pastikan pertandingan yang selesai tidak terus disimulasikan dan ticker
selalu dihentikan.

### 4.10 `tests/`, `scripts/`, dan `.github/`

-   `tests/`: integration test dan fixture jika diperlukan. Unit test
    juga dapat diletakkan berdampingan dengan package Go terkait
    menggunakan file `*_test.go`.
-   `scripts/`: utilitas seed data atau perintah pengembangan.
-   `.github/workflows/`: pipeline CI dan pemeriksaan code generation.

------------------------------------------------------------------------

## 5. Model Data

Untuk memenuhi GraphQL schema, model database yang disarankan terdiri
dari tiga tabel utama.

### 5.1 `teams`

  Kolom          Tipe konseptual   Keterangan
  -------------- ----------------- ------------------
  `id`           TEXT / UUID       Primary key
  `name`         TEXT              Nama lengkap tim
  `short_name`   TEXT              Nama singkat
  `logo_url`     TEXT, nullable    URL logo

### 5.2 `matches`

  Kolom            Tipe konseptual   Keterangan
  ---------------- ----------------- -------------------
  `id`             TEXT / UUID       Primary key
  `home_team_id`   TEXT              FK ke `teams.id`
  `away_team_id`   TEXT              FK ke `teams.id`
  `home_score`     INTEGER           Default 0
  `away_score`     INTEGER           Default 0
  `status`         TEXT              Nilai MatchStatus
  `minute`         INTEGER           Default 0
  `kickoff_at`     DATETIME          Waktu kick-off
  `updated_at`     DATETIME          Waktu pembaruan

### 5.3 `match_events`

  Kolom          Tipe konseptual   Keterangan
  -------------- ----------------- ---------------------
  `id`           TEXT / UUID       Primary key
  `match_id`     TEXT              FK ke `matches.id`
  `minute`       INTEGER           Menit kejadian
  `type`         TEXT              Nilai EventType
  `team_id`      TEXT              FK ke `teams.id`
  `player`       TEXT              Nama pemain
  `detail`       TEXT, nullable    Keterangan tambahan
  `created_at`   DATETIME          Waktu event dicatat

Relasi:

``` text
teams 1 ───────< matches (home_team_id)
teams 1 ───────< matches (away_team_id)

teams 1 ───────< match_events
matches 1 ─────< match_events
```

Tambahkan foreign key, index untuk pencarian berdasarkan status dan
kick-off, serta index pada `match_events.match_id`. Aktifkan foreign key
enforcement pada SQLite.

Catatan desain:

-   Nama tim tidak disimpan langsung di tabel `matches`; gunakan foreign
    key agar data tim konsisten.
-   `match_events` diperlukan untuk menyimpan riwayat kejadian
    pertandingan.
-   `minute` dapat menyimpan menit pertandingan yang ditampilkan,
    termasuk nilai tambahan waktu sesuai konvensi aplikasi.
-   Gunakan satu standar timezone, idealnya UTC untuk penyimpanan, lalu
    konversi saat ditampilkan.
-   Pilih satu format ID secara konsisten. UUID dalam bentuk TEXT
    merupakan pilihan praktis untuk SQLite.

------------------------------------------------------------------------

## 6. GraphQL API

GraphQL schema menjadi kontrak antara frontend dan backend.

### 6.1 Query

-   `matches(status, limit)`: mengambil daftar pertandingan dengan
    filter opsional.
-   `match(id)`: mengambil satu pertandingan.
-   `liveMatches`: mengambil pertandingan berstatus LIVE.
-   `teams`: mengambil daftar tim.

### 6.2 Mutation

-   `createMatch(input)`: membuat pertandingan.
-   `updateScore(input)`: memperbarui skor.
-   `addMatchEvent(matchId, input)`: menambahkan event.
-   `updateMatchStatus(matchId, status, minute)`: mengubah status dan
    menit pertandingan.

Setiap mutation harus melalui validasi business logic di service.
Perubahan yang berhasil dan perlu diketahui client harus memicu
publikasi event yang sesuai.

### 6.3 Subscription

-   `matchUpdated(matchId)`: menerima perubahan pertandingan, dengan
    filter opsional berdasarkan ID.
-   `newEventAdded(matchId)`: menerima event baru dari pertandingan
    tertentu.

Subscription menggunakan transport WebSocket yang didukung gqlgen.
Resolver subscription mendaftarkan subscriber ke broker dan
mengembalikan channel atau stream sesuai pola yang digunakan gqlgen.

Pastikan subscription berhenti ketika context dibatalkan atau koneksi
client ditutup.

### 6.4 Catatan tipe waktu

Jika schema menggunakan `startTime: String!`, pertimbangkan custom scalar
`Time` agar format dan parsing waktu terstandar. Jika tetap menggunakan
String, dokumentasikan format ISO 8601 yang wajib digunakan.

------------------------------------------------------------------------

## 7. Alur Data dan Interaksi Komponen

### 7.1 Query pertandingan

``` text
Client
  -> GraphQL Query
  -> Resolver
  -> MatchService
  -> MatchRepository
  -> Database
  -> MatchRepository
  -> MatchService
  -> Resolver
  -> Client
```

### 7.2 Mutation update skor

``` text
Client
  -> GraphQL Mutation
  -> Resolver
  -> MatchService
      -> Validasi skor dan status
      -> MatchRepository: update database
      -> PubSub: publish match-updated
  -> Resolver
  -> Client
```

Untuk konsistensi, tentukan urutan dan penanganan kegagalan antara
penyimpanan database dan publikasi event. Pada MVP satu proses,
publikasi dilakukan setelah operasi database berhasil. Jika kelak
dibutuhkan jaminan event tidak hilang, pertimbangkan transactional
outbox.

### 7.3 Match engine

``` text
Application Startup
  -> Match Engine
  -> Cari pertandingan LIVE
  -> Jalankan runner untuk setiap pertandingan aktif
  -> Ticker menghasilkan tick
  -> MatchService memperbarui menit / skor / status
  -> Repository menyimpan perubahan
  -> PubSub menerbitkan event
  -> GraphQL Subscription mengirim pembaruan ke client
```

Engine harus menghindari dua runner aktif untuk pertandingan yang sama.
Lifecycle runner harus dihentikan ketika pertandingan selesai atau
aplikasi shutdown.

### 7.4 Subscription

``` text
Client
  -> WebSocket
  -> gqlgen Subscription Resolver
  -> PubSub Subscribe(matchID)
  -> Menunggu event
  -> Event diterima
  -> GraphQL mengirim payload ke client
  -> Client disconnect / context cancelled
  -> Unsubscribe dan cleanup
```

------------------------------------------------------------------------

## 8. Concurrency dan Keamanan Goroutine

Concurrency menjadi concern lintas komponen pada backend. Terapkan aturan
berikut:

1.  Gunakan `context.Context` untuk mengendalikan lifecycle engine dan
    subscription.
2.  Gunakan `sync.Mutex` atau `sync.RWMutex` untuk melindungi map
    subscriber yang diakses bersama.
3.  Jangan menutup channel dari banyak goroutine. Tetapkan satu pihak
    yang bertanggung jawab atas penutupan channel.
4.  Hindari pengiriman ke channel yang sudah ditutup.
5.  Tentukan buffer dan kebijakan ketika subscriber lambat.
6.  Hentikan `time.Ticker` dengan `Stop()`.
7.  Pastikan goroutine tidak tertinggal setelah runner atau server
    dihentikan.
8.  Jalankan test dengan race detector: `go test -race ./...`.

Untuk MVP, satu goroutine per pertandingan LIVE dapat digunakan. Jika
jumlah pertandingan meningkat, evaluasi pendekatan worker atau scheduler
agar penggunaan goroutine tetap terkendali.

------------------------------------------------------------------------

## 9. Error Handling dan Validasi

Gunakan error domain/service yang dapat dipetakan oleh resolver menjadi
error GraphQL yang sesuai.

Contoh kondisi yang perlu ditangani:

-   Pertandingan tidak ditemukan.
-   Tim tidak ditemukan.
-   Tim tuan rumah dan tamu sama.
-   Skor negatif.
-   Event menggunakan tim yang bukan peserta pertandingan.
-   Menambahkan event ke pertandingan yang sudah selesai.
-   Transisi status tidak valid.
-   Database tidak tersedia.
-   Subscription dibatalkan.
-   PubSub tidak dapat mengirim karena subscriber lambat.

Jangan mengembalikan detail internal database atau stack trace kepada
client. Catat detail teknis melalui logging di sisi server.

------------------------------------------------------------------------

## 10. Testing Strategy

Dari sisi arsitektur, pengujian harus memastikan:

-   Validasi domain dan business logic dapat diuji tanpa database melalui
    mock atau fake repository.
-   Repository, migration, GraphQL, WebSocket subscription, dan engine
    dapat diuji sebagai integrasi.
-   PubSub aman saat beberapa subscriber publish/unsubscribe secara
    bersamaan.
-   Runner berhenti saat pertandingan selesai atau engine shutdown.
-   Pengujian concurrency dijalankan dengan race detector.

Hindari test yang bergantung pada `time.Sleep` panjang jika dapat
menggunakan synchronization primitive atau clock yang dapat dikendalikan.

------------------------------------------------------------------------

## 11. Konfigurasi Runtime

Konfigurasi runtime mencakup setidaknya:

``` env
APP_ENV=development
PORT=8080
DATABASE_PATH=./data/football.db
SIMULATION_INTERVAL=1s
```

Nilai tersebut adalah contoh pengembangan dan dapat disesuaikan dengan
lingkungan.

Pemeriksaan code generation harus memastikan hasil gqlgen dan sqlc tetap
sesuai dengan schema serta query yang tersimpan di repository.

------------------------------------------------------------------------

## 12. Batasan dan Rencana Pengembangan

Arsitektur MVP memiliki batasan berikut:

  -----------------------------------------------------------------------
  Area                    Batasan MVP             Opsi pengembangan
  ----------------------- ----------------------- -----------------------
  Database                SQLite, satu aplikasi   PostgreSQL

  PubSub                  In-memory, satu         Redis Pub/Sub atau
                          instance                broker

  Match engine            Mock generator          Integrasi data
                                                  pertandingan eksternal

  Deployment              Satu backend instance   Horizontal scaling

  Observability           Logging dasar           Metrics, tracing, dan
                                                  monitoring

  Authentication          Belum menjadi fokus     Authentication dan
                          awal                    authorization
  -----------------------------------------------------------------------

Perubahan teknologi harus didorong oleh kebutuhan nyata, bukan sekadar
menambah komponen.

------------------------------------------------------------------------

## 13. Ringkasan Keputusan

  Aspek                   Keputusan
  ----------------------- ------------------------------------------
  Arsitektur              Modular Monolith
  Design approach         Pragmatic Clean Architecture
  Bahasa                  Go
  GraphQL                 gqlgen, schema-first
  WebSocket               gqlgen WebSocket transport
  Database MVP            SQLite
  Database access         sqlc
  Migration               goose atau golang-migrate
  Real-time broker        In-memory PubSub dengan Go channels
  Background processing   Goroutine dan `time.Ticker`
  Synchronization         `sync.Mutex` / `sync.RWMutex`
  Testing                 Unit, integration, dan race detector

**Dokumen ini merupakan baseline arsitektur backend.** Jika ada
perubahan pada kebutuhan, GraphQL schema, atau keputusan teknologi,
perbarui dokumen ini agar implementasi dan dokumentasi tetap sinkron.
