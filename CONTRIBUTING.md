# Panduan Git Branching Strategy & Standar Kontribusi

Dokumen ini adalah acuan wajib bagi seluruh anggota tim dalam mengelola branch, menulis commit, dan mengajukan Pull Request (PR). Tujuannya: riwayat kode yang bersih, rilis yang stabil, dan kolaborasi yang konsisten.

---

## Daftar Isi

1. [Ringkasan Cepat](#1-ringkasan-cepat)
2. [Aturan Utama Branching & Proteksi](#2-aturan-utama-branching--proteksi)
3. [Struktur & Penamaan Branch](#3-struktur--penamaan-branch)
4. [Standar Conventional Commits](#4-standar-conventional-commits)
5. [Alur Kerja Pull Request](#5-alur-kerja-pull-request-pr-workflow)
6. [Alur Hotfix](#6-alur-hotfix)
7. [Template Pull Request](#7-template-pull-request)
8. [Konfigurasi Branch Protection](#8-konfigurasi-branch-protection-untuk-admin-repo)
9. [Contoh CI GitHub Actions](#9-contoh-workflow-ci-github-actions)
10. [Pedoman Code Review](#10-pedoman-code-review)
11. [Troubleshooting](#11-troubleshooting)
12. [Checklist Sebelum Membuka PR](#12-checklist-sebelum-membuka-pr)

---

## 1. Ringkasan Cepat

| Aspek | Aturan |
|---|---|
| Branch utama | `main` (production-ready), `dev` (integrasi & pengujian fitur) |
| Direct push ke `main` / `dev` | **DILARANG KERAS** |
| Perubahan ke `main` / `dev` | **WAJIB** lewat Pull Request |
| Syarat merge | Semua CI/CD **hijau** + minimal **1 approval** |
| Metode merge | `Squash and Merge` (atau `Rebase and Merge`) |
| Branch kerja | `feature/*`, `fix/*`, `chore/*` dari `dev`; `hotfix/*` dari `main` |
| Format commit | `<type>: <short summary>` (Conventional Commits) |

---

## 2. Aturan Utama Branching & Proteksi

### 2.1 Branch Utama

| Branch | Fungsi | Sumber Perubahan | Stabilitas |
|---|---|---|---|
| `main` | Kode siap produksi. Setiap commit di sini dapat di-deploy. | PR dari `dev` atau `hotfix/*` | Paling stabil |
| `dev` (`development`) | Branch integrasi untuk menggabungkan dan menguji fitur sebelum rilis. | PR dari `feature/*`, `fix/*`, `chore/*` | Stabil untuk pengujian |

### 2.2 Aturan Push

> **Direct push ke `main` dan `dev` DILARANG KERAS.** Branch Protection Rules aktif pada kedua branch ini.

Semua perubahan kode ke `dev` atau `main` **WAJIB** melalui Pull Request.

```bash
# ❌ DILARANG - akan ditolak oleh branch protection
git checkout main
git commit -m "feat: sesuatu"
git push origin main

# ✅ BENAR - selalu lewat branch kerja dan PR
git checkout -b feature/nama-fitur
git push -u origin feature/nama-fitur
# lalu buka Pull Request di GitHub
```

### 2.3 Aturan Merge & Pull Request

PR **hanya boleh di-merge** jika memenuhi **semua** syarat berikut:

- [ ] **Seluruh CI/CD di GitHub Actions LOLOS (hijau)**: build, linter, dan testing.
- [ ] Mendapat **minimal 1 Code Review / Approval** dari anggota tim pengembang.
- [ ] Seluruh percakapan review (*conversations*) sudah di-resolve.
- [ ] Branch sudah *up to date* dengan branch target.
- [ ] Menggunakan metode **`Squash and Merge`** (atau `Rebase and Merge`).

### 2.4 Mengapa Squash and Merge?

| Metode | Hasil | Rekomendasi |
|---|---|---|
| **Squash and Merge** | Seluruh commit PR dipadatkan menjadi **1 commit** di branch target. Riwayat bersih dan linear. | ✅ Default |
| **Rebase and Merge** | Commit PR disusun ulang satu per satu di atas branch target, tanpa merge commit. | ✅ Boleh, jika setiap commit sudah rapi dan bermakna |
| Merge Commit | Membuat merge commit tambahan; riwayat bercabang dan berantakan. | ❌ Tidak digunakan |

> **Catatan:** Dengan Squash and Merge, **judul PR menjadi pesan commit** di `main`/`dev`. Karena itu, judul PR **wajib** mengikuti format Conventional Commits (lihat Bagian 4).

---

## 3. Struktur & Penamaan Branch

### 3.1 Alur Branch

```
main      ●─────────────────●──────────────●────────►   (production)
           \               ▲ ▲             ▲
            \  hotfix/*    │ │ (release)   │
             \──●──●───────┘ │             │
dev        ────●──────●──────┴──────●──────┴────────►   (integrasi)
                \     ▲             ▲
                 \    │ PR          │ PR
  feature/*       ●──●─┘            │
  fix/*                       ●──●──┘
```

### 3.2 Jenis Branch

| Prefix | Tujuan | Dicabang dari | Di-merge ke | Contoh |
|---|---|---|---|---|
| `feature/` | Pengembangan fitur baru | `dev` | `dev` | `feature/graphql-subscriptions` |
| `fix/` | Perbaikan bug non-mendesak | `dev` | `dev` | `fix/websocket-reconnect` |
| `chore/` | Tugas teknis: dependency, konfigurasi, CI | `dev` | `dev` | `chore/update-dependencies` |
| `hotfix/` | Patch mendadak untuk masalah di produksi | **`main`** | `main`, lalu **back-merge ke `dev`** | `hotfix/fix-score-crash` |

### 3.3 Aturan Penamaan

- Gunakan **huruf kecil** dan pisahkan kata dengan **tanda hubung** (`kebab-case`).
- Nama singkat, deskriptif, dan menjelaskan tujuan perubahan.
- Disarankan menyertakan nomor issue bila ada, mis. `feature/123-graphql-subscriptions`.
- Jangan gunakan spasi, garis bawah, atau huruf kapital.
- Satu branch = satu tujuan. Hindari mencampur beberapa fitur dalam satu branch.

| ✅ Benar | ❌ Salah | Alasan |
|---|---|---|
| `feature/graphql-subscriptions` | `Feature/GraphQL_Subscriptions` | Huruf kapital dan underscore |
| `fix/websocket-reconnect` | `fix-bug` | Tidak deskriptif, tanpa slash |
| `chore/update-dependencies` | `my-branch` | Tidak memakai prefix |
| `hotfix/fix-score-crash` | `hotfix/fix score crash` | Mengandung spasi |

---

## 4. Standar Conventional Commits

Semua commit **wajib** mengikuti [Conventional Commits](https://www.conventionalcommits.org).

### 4.1 Format

```
<type>: <short summary>

[body opsional]

[footer opsional]
```

| Bagian | Aturan |
|---|---|
| `type` | Wajib. Salah satu dari daftar di 4.2. |
| `short summary` | Wajib. Kalimat perintah (imperatif), huruf kecil di awal, **tanpa titik** di akhir, maksimal ~72 karakter. |
| `body` | Opsional. Jelaskan **apa** dan **mengapa**, bukan *bagaimana*. Pisahkan dengan satu baris kosong. |
| `footer` | Opsional. Referensi issue (`Closes #123`) atau `BREAKING CHANGE: ...`. |

### 4.2 Jenis `type`

| Type | Kegunaan | Contoh |
|---|---|---|
| `feat` | Penambahan fitur baru | `feat(subscription): add matchUpdated subscription` |
| `fix` | Perbaikan bug/error | `fix(websocket): handle reconnect after timeout` |
| `docs` | Perubahan dokumen saja | `docs(readme): add setup instructions` |
| `refactor` | Restrukturisasi kode tanpa mengubah fungsionalitas | `refactor(service): extract score validation` |
| `test` | Menambah atau memperbaiki automated testing | `test(resolver): add tests for updateScore` |
| `chore` | Update build task, paket/dependency, atau konfigurasi CI | `chore(deps): bump gqlgen to latest version` |

### 4.3 Contoh Commit yang BENAR ✅

```bash
git commit -m "feat: add matchUpdated subscription"
git commit -m "fix: handle reconnect after connection timeout"
git commit -m "docs: add branch naming guide"
git commit -m "refactor: extract score validation into helper"
git commit -m "test: add unit tests for updateScore mutation"
git commit -m "chore: add golangci-lint step to workflow"
```

Contoh commit dengan body dan footer:

```bash
git commit -m "fix: prevent goroutine leak on client disconnect" \
  -m "Subscriber channel is now removed from the broker when the request context is cancelled." \
  -m "Closes #42"
```

Contoh *breaking change*:

```bash
git commit -m "feat!: change updateScore to require admin token" \
  -m "BREAKING CHANGE: clients must send the Authorization header when calling updateScore."
```

### 4.4 Contoh Commit yang SALAH ❌

| Pesan Salah | Masalah | Perbaikan |
|---|---|---|
| `update code` | Tidak ada type, tidak jelas | `refactor: simplify match update logic` |
| `fixed bug` | Tidak ada type/scope, tidak spesifik | `fix: return null when match not found` |
| `Feat: Add Login.` | Type kapital, summary kapital, ada titik | `feat: add login endpoint` |
| `feat add subscription` | Tanda titik dua hilang | `feat: add subscription` |
| `feature: add score` | `feature` bukan type valid | `feat: add score field` |
| `fix: fixed the bug where the score was wrong and also changed the config and updated docs` | Terlalu panjang dan mencampur banyak hal | Pecah menjadi beberapa commit terpisah |
| `WIP` | Tidak informatif | Gunakan pesan bermakna, atau rapikan sebelum PR |
| `asdf` | Tidak bermakna | Tulis pesan yang menjelaskan perubahan |

### 4.5 Tips Menulis Commit

- Satu commit = satu perubahan logis.
- Gunakan kalimat perintah: "add", "fix", "update", bukan "added" atau "fixes".
- Jika ragu soal *type*, tanyakan: "Apakah ini mengubah perilaku untuk pengguna?" Jika ya, kemungkinan `feat` atau `fix`.
- Karena PR di-squash, **judul PR** harus juga berformat Conventional Commits.

---

## 5. Alur Kerja Pull Request (PR Workflow)

### Ringkasan Langkah

| # | Langkah | Perintah / Aksi |
|---|---|---|
| 1 | Tarik update terbaru dari `dev` | `git pull origin dev` |
| 2 | Buat branch baru dari `dev` | `git checkout -b feature/<nama>` |
| 3 | Commit sesuai Conventional Commits | `git commit -m "feat: ..."` |
| 4 | Push branch & buka PR ke `dev` | `git push -u origin <branch>` |
| 5 | Isi template PR | Deskripsi, link issue, checklist testing |
| 6 | Pastikan semua Checks GitHub Actions hijau | Perbaiki bila gagal |
| 7 | Minta review & lakukan koreksi | Tangani *request changes* |
| 8 | Merge setelah approve & CI pass | `Squash and Merge` |

### Langkah 1: Tarik Update Terbaru dari `dev`

Selalu mulai dari kondisi `dev` terbaru agar menghindari konflik di kemudian hari.

```bash
git checkout dev
git pull origin dev
```

### Langkah 2: Buat Branch Baru dari `dev`

```bash
git checkout -b feature/graphql-subscriptions
```

### Langkah 3: Lakukan Commit Sesuai Conventional Commits

```bash
git add .
git commit -m "feat: add matchUpdated subscription"

# lanjutkan bekerja, commit secara bertahap
git add internal/broker/broker.go
git commit -m "feat: add pubsub broker for match events"
```

### Langkah 4: Push Branch & Buka PR ke `dev`

```bash
git push -u origin feature/graphql-subscriptions
```

Di GitHub, buka **Pull Request** dengan:

- **Base:** `dev`
- **Compare:** `feature/graphql-subscriptions`
- **Judul:** mengikuti Conventional Commits, mis. `feat: add matchUpdated subscription`

> Jika menggunakan GitHub CLI:
> ```bash
> gh pr create --base dev --title "feat: add matchUpdated subscription"
> ```

### Langkah 5: Isi Template PR

Lengkapi seluruh bagian template (lihat [Bagian 7](#7-template-pull-request)):

- **Deskripsi perubahan:** apa yang berubah dan mengapa.
- **Link issue/task:** mis. `Closes #123`.
- **Checklist testing:** bukti bahwa perubahan sudah diuji.

### Langkah 6: Pastikan Semua Checks GitHub Actions Berhasil

Tunggu hingga seluruh *check* berstatus hijau (build, linter, test). Jika ada yang gagal:

```bash
# jalankan pemeriksaan yang sama secara lokal
go build ./...
golangci-lint run
go test -race ./...

# perbaiki, lalu commit dan push ulang
git add .
git commit -m "fix: resolve lint errors in broker"
git push
```

> PR **tidak dapat di-merge** selama ada check yang merah.

### Langkah 7: Minta Review & Lakukan Koreksi

- Tambahkan minimal **1 reviewer** dari tim.
- Jika reviewer memilih **Request changes**, perbaiki lalu push commit tambahan ke branch yang sama (PR otomatis diperbarui).
- Balas dan **resolve** setiap komentar setelah ditangani.
- Setelah selesai, minta *re-review*.

```bash
git add .
git commit -m "refactor: address review feedback on locking"
git push
```

Jika `dev` sudah berubah dan branch Anda tertinggal:

```bash
git fetch origin
git rebase origin/dev        # atau: git merge origin/dev
# selesaikan konflik jika ada, lalu:
git push --force-with-lease  # hanya setelah rebase, hanya pada branch milik sendiri
```

> ⚠️ **Jangan pernah** melakukan `git push --force` ke `main` atau `dev`. Gunakan `--force-with-lease` hanya pada branch kerja Anda sendiri.

### Langkah 8: Merge PR Setelah Approve & CI Pass

Jika semua syarat terpenuhi (CI hijau, ≥ 1 approval, percakapan ter-resolve):

1. Pilih **`Squash and Merge`** (atau `Rebase and Merge`).
2. Pastikan judul commit hasil squash mengikuti Conventional Commits.
3. Hapus branch kerja setelah merge.

```bash
# bersihkan branch lokal
git checkout dev
git pull origin dev
git branch -d feature/graphql-subscriptions
git fetch --prune
```

---

## 6. Alur Hotfix

Hotfix digunakan **hanya** untuk perbaikan mendesak pada kode produksi (`main`).

```bash
# 1. Cabang dari main
git checkout main
git pull origin main
git checkout -b hotfix/fix-score-crash

# 2. Perbaiki dan commit
git add .
git commit -m "fix: prevent crash when score is negative"

# 3. Push dan buka PR ke main
git push -u origin hotfix/fix-score-crash
```

| Tahap | Aturan |
|---|---|
| PR ke `main` | Tetap wajib CI hijau dan ≥ 1 approval. |
| Setelah merge ke `main` | Buat tag rilis patch (mis. `v1.2.1`). |
| **Back-merge** | Perbaikan **wajib** dibawa kembali ke `dev` agar tidak hilang pada rilis berikutnya. |

```bash
# Back-merge hotfix ke dev (lewat PR)
git checkout dev
git pull origin dev
git checkout -b chore/backmerge-hotfix-fix-score-crash
git merge main
git push -u origin chore/backmerge-hotfix-fix-score-crash
# buka PR ke dev
```

### Rilis dari `dev` ke `main`

Saat fitur di `dev` sudah stabil dan siap rilis, buka PR dari `dev` ke `main` dengan syarat yang sama (CI hijau, ≥ 1 approval).

---

## 7. Template Pull Request

Simpan sebagai `.github/pull_request_template.md` agar otomatis muncul pada setiap PR.

```markdown
## Deskripsi Perubahan
<!-- Jelaskan APA yang berubah dan MENGAPA. -->


## Jenis Perubahan
- [ ] feat: fitur baru
- [ ] fix: perbaikan bug
- [ ] refactor: restrukturisasi kode
- [ ] docs: dokumentasi
- [ ] test: penambahan/perbaikan test
- [ ] chore: dependency/konfigurasi/CI

## Link Issue / Task
<!-- Contoh: Closes #123 -->


## Checklist Testing
- [ ] Unit test baru/yang ada sudah dijalankan dan lulus
- [ ] Sudah diuji secara manual (jelaskan langkahnya di bawah)
- [ ] `go test -race ./...` lulus (jika relevan)
- [ ] Linter lulus tanpa temuan

## Langkah Pengujian Manual
<!-- Jelaskan cara reviewer memverifikasi perubahan ini. -->


## Screenshot / Log (jika relevan)


## Checklist Umum
- [ ] Judul PR mengikuti format Conventional Commits
- [ ] Branch dicabang dari `dev` (atau `main` untuk hotfix)
- [ ] Dokumentasi sudah diperbarui (jika perlu)
- [ ] Tidak ada secret/kredensial yang ikut ter-commit
- [ ] Semua GitHub Actions Checks berwarna hijau
```

---

## 8. Konfigurasi Branch Protection (untuk Admin Repo)

Atur di **Settings → Branches → Add branch protection rule** untuk `main` **dan** `dev`.

| Pengaturan | Nilai |
|---|---|
| Require a pull request before merging | ✅ Aktif |
| Required approvals | **1** (minimal) |
| Dismiss stale approvals when new commits are pushed | ✅ Direkomendasikan |
| Require status checks to pass before merging | ✅ Aktif |
| Required checks | `build`, `lint`, `test` (sesuai nama job CI) |
| Require branches to be up to date before merging | ✅ Direkomendasikan |
| Require conversation resolution before merging | ✅ Direkomendasikan |
| Require linear history | ✅ Direkomendasikan (selaras dengan Squash/Rebase) |
| Do not allow bypassing the above settings | ✅ Aktif (berlaku juga untuk admin) |
| Allow force pushes | ❌ Nonaktif |
| Allow deletions | ❌ Nonaktif |

Pada **Settings → General → Pull Requests**:

| Pengaturan | Nilai |
|---|---|
| Allow squash merging | ✅ (default: *Pull request title and description*) |
| Allow rebase merging | ✅ (opsional) |
| Allow merge commits | ❌ Nonaktif |
| Automatically delete head branches | ✅ Aktif |

---

## 9. Contoh Workflow CI (GitHub Actions)

Contoh `.github/workflows/ci.yml` untuk proyek Go. Sesuaikan dengan stack proyek Anda. Nama job (`build`, `lint`, `test`) harus sama dengan *required checks* di branch protection.

```yaml
name: CI

on:
  pull_request:
    branches: [main, dev]
  push:
    branches: [main, dev]

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version-file: go.mod
      - name: Build
        run: go build ./...

  lint:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version-file: go.mod
      - name: golangci-lint
        uses: golangci/golangci-lint-action@v6

  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version-file: go.mod
      - name: Test (with race detector)
        run: go test -race -cover ./...
```

Opsional: validasi judul PR agar sesuai Conventional Commits.

```yaml
name: PR Title

on:
  pull_request:
    types: [opened, edited, synchronize, reopened]

jobs:
  semantic-title:
    runs-on: ubuntu-latest
    steps:
      - uses: amannn/action-semantic-pull-request@v5
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
```

---

## 10. Pedoman Code Review

### Untuk Author

- Jaga PR tetap **kecil dan fokus** (idealnya < 400 baris perubahan).
- Tulis deskripsi yang membantu reviewer memahami konteks.
- Tanggapi komentar dengan sopan dan segera.
- Jangan merge PR sendiri tanpa approval.

### Untuk Reviewer

| Periksa | Pertanyaan Panduan |
|---|---|
| Kebenaran | Apakah kode melakukan apa yang dijelaskan? Edge case tertangani? |
| Desain | Apakah sesuai arsitektur dan konvensi proyek? |
| Pengujian | Apakah ada test yang memadai dan bermakna? |
| Keamanan | Adakah input tidak tervalidasi, secret ter-commit, atau celah lain? |
| Keterbacaan | Apakah penamaan jelas dan kode mudah dipahami? |
| Performa | Adakah potensi masalah (mis. goroutine leak, query boros)? |

- Berikan umpan balik yang **konstruktif dan spesifik**.
- Bedakan komentar **wajib diperbaiki** dari **saran** (awali dengan `nit:` untuk hal kecil).
- Targetkan review dalam **1 hari kerja**.

---

## 11. Troubleshooting

| Masalah | Solusi |
|---|---|
| Push ke `dev`/`main` ditolak | Memang diblokir. Buat branch kerja dan buka PR. |
| Salah commit langsung di `dev` secara lokal | `git branch feature/nama-baru` lalu `git reset --hard origin/dev` di `dev`, dan lanjutkan di branch baru. |
| Konflik saat PR | `git fetch origin && git rebase origin/dev`, selesaikan konflik, lalu `git push --force-with-lease`. |
| Salah pesan commit terakhir (belum di-push) | `git commit --amend -m "pesan-baru"` |
| Perlu merapikan beberapa commit sebelum PR | `git rebase -i origin/dev` lalu `squash`/`reword`. |
| Check CI gagal hanya di GitHub | Jalankan perintah CI yang sama secara lokal; periksa versi tool dan environment. |
| Ingin membatalkan commit yang sudah di-merge | Buat branch dan jalankan `git revert <hash>`, lalu buka PR. |
| Tidak sengaja commit secret | Segera **cabut/ganti secret**, lalu bersihkan riwayat dan hubungi admin repo. Menghapus commit saja tidak cukup. |

Contoh memindahkan commit yang salah tempat:

```bash
# Anda terlanjur commit di dev lokal (belum di-push)
git checkout dev
git branch feature/nama-baru        # simpan commit ke branch baru
git reset --hard origin/dev         # kembalikan dev ke kondisi remote
git checkout feature/nama-baru      # lanjutkan di sini, lalu push dan buat PR
```

---

## 12. Checklist Sebelum Membuka PR

- [ ] Branch dicabang dari `dev` (atau `main` untuk hotfix) dan namanya sesuai aturan.
- [ ] Semua commit mengikuti Conventional Commits.
- [ ] Branch sudah up to date dengan `dev`.
- [ ] Build, linter, dan test lulus di lokal.
- [ ] Tidak ada file sementara, log debug, atau secret yang ikut ter-commit.
- [ ] Judul PR mengikuti Conventional Commits.
- [ ] Template PR terisi lengkap (deskripsi, link issue, checklist testing).
- [ ] Reviewer sudah ditambahkan.

---

## Lampiran: Cheat Sheet Perintah

```bash
# Mulai pekerjaan baru
git checkout dev && git pull origin dev
git checkout -b feature/nama-fitur

# Simpan pekerjaan
git add .
git commit -m "feat(scope): ringkasan singkat"

# Kirim dan buka PR
git push -u origin feature/nama-fitur

# Sinkronisasi dengan dev
git fetch origin
git rebase origin/dev
git push --force-with-lease

# Bersih-bersih setelah merge
git checkout dev && git pull origin dev
git branch -d feature/nama-fitur
git fetch --prune
```

---

*Dokumen ini adalah living document. Usulan perubahan diajukan melalui PR dengan type `docs`.*
