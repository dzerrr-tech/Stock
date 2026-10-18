# Config, context, pool pgx, constraint

- `internal/`: hanya bisa di-import dari dalam module sendiri (ditegakkan compiler).
- Config dari env + godotenv (dev saja); env var host selalu menang di produksi.
- `context.Context` membawa timeout dan pembatalan; context request HTTP diteruskan sampai query.
- `pgxpool.Pool` aman dipakai banyak goroutine; satu pool per aplikasi, dibuat di `main`.
- Wrap error dengan `fmt.Errorf("...: %w", err)` supaya error asli tetap bisa diperiksa.
- Invarian bisnis (stok >= 0, SKU unik) dijaga CHECK/UNIQUE di database; constraint diberi nama.
- Migrasi lewat koneksi direct, bukan pooler PgBouncer.
- Liveness (`/health`) dipisah dari readiness (`/health/db`) agar health check platform tidak membangunkan DB.

# Service, repository, dan tes

- Handler -> service -> repository: handler urus HTTP, service urus aturan bisnis, repository urus SQL.
- `errors.As`: target harus sama dengan tipe konkret error. `ValidationError` (nilai) -> `var ve ValidationError`; `*pgconn.PgError` (pointer) -> `var pgErr *pgconn.PgError`.
- Sentinel error (`ErrDuplicateSKU`) menerjemahkan error DB antar layer.
- Table-driven test + `t.Run`; `t.Fatalf` (stop) vs `t.Errorf` (lanjut).
- Tes integrasi: `t.Skip` bila env kosong, `t.Cleanup` (LIFO) untuk membersihkan data.
- Mutation check: ubah kode sebentar untuk memastikan tes benar-benar bisa gagal.