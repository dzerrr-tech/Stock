# PROGRESS

## Selesai
- [x] Langkah 1: kerangka backend Go + Gin, `GET /api/v1/health` (diuji manual dengan curl, 200 OK)
- [x] Langkah 2: PostgreSQL di Neon (branch `dev`), migrasi `products` (golang-migrate), config dari env, pool pgx, `GET /api/v1/health/db`

## Perintah yang berjalan
Done:
- Jalankan API: `cd apps\api` lalu `go run ./cmd/api` (port 8080)
- Cek: `curl.exe http://localhost:8080/api/v1/health`
Done:
- Migrasi: `cd apps\api`, muat `DATABASE_URL` dari `.env`, lalu `migrate -path migrations -database $env:DATABASE_URL up`
- Jalankan API: `go run ./cmd/api` dari `apps\api` (harus dari folder itu, supaya `.env` terbaca)
- Cek: `/api/v1/health` (liveness, tanpa DB) dan `/api/v1/health/db` (ping DB)

## Catatan teknis / known issues
Done:
- Port 8080 masih hardcode, nanti dibaca dari env (`internal/config`)
- Gin warning "trusted all proxies": atur `SetTrustedProxies` saat setup config/deploy
- Handler health masih inline di `main.go`, pindahkan ke `internal/handler/` saat endpoint kedua muncul
Done:
- Jalur 503 di `/health/db` belum teruji (celah tes)
- Neon Free suspend setelah 5 menit idle; ukur cold start dan putuskan pooler vs direct saat deploy
- Belum ada branch Neon untuk demo publik (baru `dev`)

## Tugas berikutnya
- Langkah 2: PostgreSQL lokal + koneksi pgx + migrasi tabel `products`

## Selesai
- Membuat Config, Database dan Handler, Diuji dari go run ./cmd/api 

## Tugas berikutnya
-