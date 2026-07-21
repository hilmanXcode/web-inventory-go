# 🔍 Analisis & Rekomendasi Codebase: web-inventory-go

## Gambaran Umum

Ini adalah aplikasi web inventory management berbasis Go (tanpa framework pihak ketiga untuk HTTP) yang menggunakan MySQL, server-side rendering dengan `html/template`, dan session management custom berbasis in-memory map.

---

## 1. 🚨 Masalah Kritis: Keamanan

### 1.1. Credentials Hardcoded di Source Code

**File:** [database.go](file:///home/manz/Documents/Code/web-inventory-go-cloned/database/database.go#L15)

```go
var dsn = "manz:supersecretpassword@tcp(127.0.0.1:3306)/db_inventory"
```

**Masalah:** Username dan password database ditulis langsung di dalam source code. Jika project ini di-push ke GitHub (public), siapa saja bisa melihat credentials ini.

**Rekomendasi:** Gunakan **environment variables** atau file `.env` (dengan library seperti `godotenv`).

```go
// Contoh yang benar:
dsn := os.Getenv("DATABASE_DSN")
// atau
dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s",
    os.Getenv("DB_USER"),
    os.Getenv("DB_PASS"),
    os.Getenv("DB_HOST"),
    os.Getenv("DB_NAME"),
)
```

> [!CAUTION]
> Ini adalah **kesalahan keamanan paling kritis** di codebase ini. Jangan pernah menyimpan credentials di source code.

**Referensi:**
- [The Twelve-Factor App - Config](https://12factor.net/config) — Prinsip industri tentang mengapa konfigurasi harus dari environment
- [OWASP - Hardcoded Credentials](https://owasp.org/www-community/vulnerabilities/Use_of_hard-coded_password)

---

### 1.2. Cookie Tidak Aman

**File:** [sessions.go](file:///home/manz/Documents/Code/web-inventory-go-cloned/sessions/sessions.go#L44-L49)

```go
http.SetCookie(w, &http.Cookie{
    Name:    "session_token",
    Value:   sessionToken,
    Expires: expiresAt,
    Path:    "/",
})
```

**Masalah:** Cookie tidak menggunakan `HttpOnly`, `Secure`, atau `SameSite` flags. Ini membuat aplikasi rentan terhadap **XSS** (pencurian cookie via JavaScript) dan **CSRF** (cross-site request forgery).

**Rekomendasi:**
```go
http.SetCookie(w, &http.Cookie{
    Name:     "session_token",
    Value:    sessionToken,
    Expires:  expiresAt,
    Path:     "/",
    HttpOnly: true,              // Tidak bisa diakses via JavaScript
    Secure:   true,              // Hanya dikirim via HTTPS
    SameSite: http.SameSiteStrictMode, // Cegah CSRF
})
```

**Referensi:**
- [OWASP - Session Management Cheat Sheet](https://cheatsheetseries.owasp.org/cheatsheets/Session_Management_Cheat_Sheet.html)
- [MDN - Set-Cookie](https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Set-Cookie)

---

### 1.3. XSS via Template

**File:** [login.html](file:///home/manz/Documents/Code/web-inventory-go-cloned/views/auth/login.html#L84)

```html
<script>
    toast.error('Gagal', '{{ $elem }}')
</script>
```

**Masalah:** Data user langsung dimasukkan ke dalam tag `<script>`. Meskipun Go `html/template` melakukan escaping untuk HTML context, perilakunya di dalam tag `<script>` **tidak menjamin keamanan penuh**. Attacker bisa menyisipkan string yang merusak JavaScript.

**Rekomendasi:** Render data sebagai data attribute di HTML, lalu baca dari JavaScript:

```html
<!-- Aman: di-escape oleh html/template -->
<div id="error-data" data-errors='{{ .errorJSON }}'></div>
<script>
    const data = document.getElementById('error-data');
    const errors = JSON.parse(data.dataset.errors);
    errors.forEach(e => toast.error('Gagal', e));
</script>
```

**Referensi:**
- [OWASP - XSS Prevention Cheat Sheet](https://cheatsheetseries.owasp.org/cheatsheets/Cross-Site_Scripting_Prevention_Cheat_Sheet.html)

---

## 2. 🏗️ Struktur Folder

### Struktur Saat Ini

```
web-inventory-go-cloned/
├── const/          ← hanya 1 file (4 baris)
├── database/       ← koneksi + raw SQL query helpers
├── entities/       ← model / struct
├── handlers/       ← HTTP handler
├── routes/         ← routing + inline handler logic
├── sessions/       ← session management
├── static/         ← aset statis
├── utils/          ← validator + view helper (campur)
├── views/          ← template cache + HTML templates
├── main.go
├── go.mod
└── go.sum
```

### Masalah

| Masalah | Penjelasan |
|---------|-----------|
| `const/` berisi hanya 1 konstanta | Overkill membuat folder sendiri untuk ini |
| `utils/` campur aduk | `validator.go` dan `viewsWdata.go` punya tanggung jawab yang sangat berbeda tapi digabung satu package |
| Routes punya inline handler | [routes.go](file:///home/manz/Documents/Code/web-inventory-go-cloned/routes/routes.go#L30-L81) berisi handler logic langsung (setCookie, myCookie, clearCookie) alih-alih di `handlers/` |
| Tidak ada layer service/repository | Handler langsung memanggil database — tidak ada pemisahan business logic |
| Tidak ada file `config` | Konfigurasi (port, DSN) hardcoded di berbagai tempat |

### Struktur yang Direkomendasikan

```
web-inventory-go/
├── cmd/
│   └── server/
│       └── main.go              ← entry point
├── internal/
│   ├── config/
│   │   └── config.go            ← load env vars, struct Config
│   ├── handler/
│   │   ├── auth.go              ← auth handlers
│   │   └── session.go           ← session-related handlers
│   ├── middleware/
│   │   └── auth.go              ← auth middleware
│   ├── model/
│   │   └── user.go              ← domain models
│   ├── repository/
│   │   └── user_repo.go         ← database queries
│   ├── service/
│   │   └── auth_service.go      ← business logic
│   └── session/
│       └── session.go           ← session management
├── web/
│   ├── static/                  ← CSS, JS, images
│   └── template/                ← HTML templates
├── go.mod
├── go.sum
├── .env.example                 ← contoh env vars
└── README.md
```

**Kenapa struktur ini lebih baik?**
- **`cmd/`**: Standar Go — entry point dipisahkan sehingga bisa punya multiple binaries (server, CLI migration, dll)
- **`internal/`**: Package di bawah `internal/` **tidak bisa di-import dari luar module** — ini fitur bawaan Go untuk encapsulation
- **Separation of Concerns**: Handler → Service → Repository. Setiap layer punya satu tanggung jawab
- **`web/`**: Template dan asset statis dikelompokkan bersama

**Referensi:**
- [Standard Go Project Layout](https://github.com/golang-standards/project-layout) — Panduan komunitas Go untuk struktur project
- [Go Blog - Organizing a Go Module](https://go.dev/doc/modules/layout)
- [Ben Johnson - Standard Package Layout](https://www.gobeyond.dev/standard-package-layout/)

---

## 3. ⚠️ Masalah Kode

### 3.1. `panic()` di Production Code

**File:** [database.go](file:///home/manz/Documents/Code/web-inventory-go-cloned/database/database.go#L50)

```go
func InsertQuery(query string, params ...any) {
    _, err := DB.Exec(query, params...)
    if err != nil {
        panic(err.Error()) // ❌ Server langsung crash
    }
}
```

**Masalah:** `panic()` akan langsung **menghentikan seluruh aplikasi** jika terjadi error insert. Di production, ini berarti server mati hanya karena satu query gagal.

**Rekomendasi:** Return error dan handle di caller:

```go
func InsertQuery(query string, params ...any) (sql.Result, error) {
    result, err := DB.Exec(query, params...)
    if err != nil {
        return nil, fmt.Errorf("insert query failed: %w", err)
    }
    return result, nil
}
```

**Referensi:**
- [Effective Go - Panic](https://go.dev/doc/effective_go#panic) — Kapan boleh dan tidak boleh pakai panic
- [Go Blog - Error Handling](https://go.dev/blog/error-handling-and-go)

---

### 3.2. `fmt.Println` untuk Logging

**File:** [database.go](file:///home/manz/Documents/Code/web-inventory-go-cloned/database/database.go#L53), [auth_handler.go](file:///home/manz/Documents/Code/web-inventory-go-cloned/handlers/auth_handler.go#L56)

```go
fmt.Println("Sukses melakukan insert")
fmt.Println(sessions.SessionData[myKey].OldInput)
```

**Masalah:** `fmt.Println` tidak ada timestamp, level, atau context. Sangat sulit untuk debugging di production.

**Rekomendasi:** Gunakan `log/slog` (bawaan Go 1.21+) atau library seperti `zerolog`:

```go
slog.Info("insert query executed",
    "table", "users",
    "email", reqs.Email,
)
```

**Referensi:**
- [Go 1.21 - slog package](https://pkg.go.dev/log/slog) — Structured logging bawaan Go

---

### 3.3. Module Path Tidak Konsisten

**File:** [go.mod](file:///home/manz/Documents/Code/web-inventory-go-cloned/go.mod#L1) vs [main.go](file:///home/manz/Documents/Code/web-inventory-go-cloned/main.go#L7-L8)

```
go.mod:  module github.com/hilmanxcode/web-inventory-go-cloned
main.go: import "github.com/hilmanxcode/web-inventory-go/database"   ← BEDA!
```

`go.mod` mendeklarasikan module sebagai `web-inventory-go-cloned` tapi `main.go` meng-import `web-inventory-go` (tanpa `-cloned`). Ini harusnya menyebabkan **build error**.

**Rekomendasi:** Konsistenkan semua import path sesuai dengan apa yang ditulis di `go.mod`.

---

### 3.4. Session In-Memory Tidak Scalable

**File:** [sessions.go](file:///home/manz/Documents/Code/web-inventory-go-cloned/sessions/sessions.go#L13)

```go
var SessionData = map[string]Session{}
```

**Masalah:**
- Data session hilang setiap kali server restart
- Tidak ada goroutine untuk membersihkan session yang sudah expired — **memory leak**
- Jika deploy multiple instances (load balancer), session tidak bisa di-share

**Rekomendasi (bertahap):**
1. **Jangka pendek:** Tambahkan goroutine cleanup untuk session expired
2. **Jangka panjang:** Gunakan Redis atau database untuk session storage

```go
// Contoh cleanup goroutine
func StartSessionCleanup(interval time.Duration) {
    ticker := time.NewTicker(interval)
    go func() {
        for range ticker.C {
            mu.Lock()
            for key, session := range SessionData {
                if session.IsExpired() {
                    delete(SessionData, key)
                }
            }
            mu.Unlock()
        }
    }()
}
```

**Referensi:**
- [Alex Edwards - Sessions in Go](https://www.alexedwards.net/blog/scs-session-manager) — SCS session manager untuk Go

---

### 3.5. Global Mutable State (`var DB *sql.DB`)

**File:** [database.go](file:///home/manz/Documents/Code/web-inventory-go-cloned/database/database.go#L11)

**Masalah:** Global variable membuat kode sulit di-test (tidak bisa mock database) dan rawan race condition.

**Rekomendasi:** Gunakan **dependency injection** — pass `*sql.DB` sebagai parameter atau field struct:

```go
type UserRepository struct {
    db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
    return &UserRepository{db: db}
}

func (r *UserRepository) Create(user User) error {
    _, err := r.db.Exec("INSERT INTO users ...", ...)
    return err
}
```

**Referensi:**
- [Go Wiki - Style Decisions - Global State](https://google.github.io/styleguide/go/decisions#global-state)
- [Uber Go Style Guide](https://github.com/uber-go/guide/blob/master/style.md)

---

### 3.6. Exported Variable `SessionData`

**File:** [sessions.go](file:///home/manz/Documents/Code/web-inventory-go-cloned/sessions/sessions.go#L13)

```go
var SessionData = map[string]Session{}
```

`SessionData` di-export (huruf besar) sehingga package lain bisa langsung memanipulasi map ini tanpa melalui fungsi yang proper. Ini melanggar prinsip **encapsulation**.

**Rekomendasi:** Jadikan `unexported` dan expose hanya via functions:
```go
var sessionData = map[string]Session{} // huruf kecil = private

func GetSession(key string) (Session, bool) { ... }
```

---

## 4. 📝 Kualitas Kode

### 4.1. Commented-out Code

**File:** [main.go](file:///home/manz/Documents/Code/web-inventory-go-cloned/main.go#L15-L27), [validator.go](file:///home/manz/Documents/Code/web-inventory-go-cloned/utils/validator.go#L39-L47)

```go
// reqs := entities.User{...}
// err, message := utils.Validate(reqs)

// tidak di hapus karna ini buat catatan hehehe :D
```

**Masalah:** Commented-out code membuat codebase berantakan dan membingungkan. Git sudah menyimpan history — jadi tidak perlu "menyimpan" kode lama di comment.

**Rekomendasi:** Hapus semua commented-out code. Gunakan Git untuk melihat history jika perlu.

**Referensi:**
- [Clean Code by Robert C. Martin](https://www.oreilly.com/library/view/clean-code-a/9780136083238/) — Bab tentang Comments

---

### 4.2. Penamaan File Tidak Konsisten

| File | Masalah |
|------|---------|
| `viewsWdata.go` | CamelCase — harusnya `views_with_data.go` |
| `viewsconst.go` | Tanpa separator — harusnya `views_const.go` |
| `auth_handler.go` | ✅ Ini sudah benar (snake_case) |

**Konvensi Go:** Nama file menggunakan **snake_case** (huruf kecil, dipisah underscore).

**Referensi:**
- [Go Code Review Comments](https://go.dev/wiki/CodeReviewComments)

---

### 4.3. LoginHandler Tidak Benar-Benar Login

**File:** [auth_handler.go](file:///home/manz/Documents/Code/web-inventory-go-cloned/handlers/auth_handler.go#L20-L27)

```go
func LoginHandler(w http.ResponseWriter, r *http.Request) {
    var email = r.FormValue("email")
    var password = r.FormValue("password")
    var result = fmt.Sprintf("Email: %s\nPassword: %s", email, password)
    w.Write([]byte(result))  // ← Hanya menampilkan email & password!
}
```

**Masalah:** Handler ini **tidak melakukan autentikasi sama sekali** — hanya mencetak kembali email dan password. Ini juga membocorkan password di response.

---

## 5. 🧪 Tidak Ada Testing

Saat ini tidak ada satupun file test (`*_test.go`) di seluruh project.

**Rekomendasi:** Mulai dengan unit test untuk logic yang paling penting:

```go
// utils/validator_test.go
func TestValidate_RequiredFieldEmpty(t *testing.T) {
    user := entities.User{Email: "", NamaLengkap: "", Password: ""}
    invalid, messages := utils.Validate(user)

    if !invalid {
        t.Error("expected validation to fail for empty fields")
    }
    if len(messages) == 0 {
        t.Error("expected error messages")
    }
}
```

**Referensi:**
- [Go Blog - Table Driven Tests](https://go.dev/blog/subtests) — Pola testing idiomatis di Go
- [Learn Go with Tests](https://quii.gitbook.io/learn-go-with-tests/) — Tutorial gratis belajar Go via TDD

---

## 6. 📚 Ringkasan Referensi Utama

| Topik | Referensi |
|-------|-----------|
| Struktur Project Go | [golang-standards/project-layout](https://github.com/golang-standards/project-layout) |
| Go Style Guide | [Uber Go Style Guide](https://github.com/uber-go/guide/blob/master/style.md) |
| Google Go Style | [Google Go Style Decisions](https://google.github.io/styleguide/go/) |
| Error Handling | [Go Blog - Error Handling](https://go.dev/blog/error-handling-and-go) |
| Effective Go | [Effective Go (Official)](https://go.dev/doc/effective_go) |
| Keamanan Web | [OWASP Top 10](https://owasp.org/www-project-top-ten/) |
| Session Management | [OWASP Session Cheat Sheet](https://cheatsheetseries.owasp.org/cheatsheets/Session_Management_Cheat_Sheet.html) |
| Clean Code | [Clean Code - Robert C. Martin](https://www.oreilly.com/library/view/clean-code-a/9780136083238/) |
| Belajar Go via TDD | [Learn Go with Tests](https://quii.gitbook.io/learn-go-with-tests/) |
| 12 Factor App | [12factor.net](https://12factor.net/) |

---

## 7. ✅ Prioritas Perbaikan (Urutan Disarankan)

| Prioritas | Item | Effort |
|-----------|------|--------|
| 🔴 P0 | Hapus hardcoded credentials, gunakan env vars | Kecil |
| 🔴 P0 | Fix module path yang tidak konsisten (`main.go`) | Kecil |
| 🟠 P1 | Tambahkan `HttpOnly`, `Secure`, `SameSite` ke cookie | Kecil |
| 🟠 P1 | Ganti `panic()` dengan proper error return | Kecil |
| 🟠 P1 | Implementasikan `LoginHandler` yang sebenarnya | Sedang |
| 🟡 P2 | Pindahkan inline handlers dari `routes.go` ke `handlers/` | Kecil |
| 🟡 P2 | Hapus commented-out code | Kecil |
| 🟡 P2 | Konsistenkan penamaan file (snake_case) | Kecil |
| 🔵 P3 | Tambahkan session cleanup goroutine | Sedang |
| 🔵 P3 | Refactor ke pattern Repository + Service | Besar |
| 🔵 P3 | Tambahkan unit tests | Besar |
| 🔵 P3 | Restructure ke layout `cmd/internal/web` | Besar |

---

> [!TIP]
> Mulai dari **P0** dan **P1** yang effort-nya kecil tapi dampak keamanannya besar. Untuk P3 (refactor besar), lakukan secara bertahap — tidak perlu sekaligus.
