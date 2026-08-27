# Go Microservice — Simple CRUD Blog API

Contoh proyek **microservice sederhana** menggunakan **Golang**, dengan **API Gateway**, **JWT Authentication**, dan **Docker Compose**.

## Arsitektur

```
Client (Postman / Frontend)
           │
           ▼  :8080
   ┌───────────────┐
   │  API Gateway  │  ← Satu-satunya yang exposed ke luar
   │  (validasi JWT│
   │   & routing)  │
   └───────┬───────┘
           │ Internal Docker Network
     ┌─────┴──────┐
     │            │
     ▼            ▼
┌──────────┐  ┌──────────┐
│   User   │  │   Post   │
│ Service  │  │ Service  │
│  :8081   │  │  :8082   │
└────┬─────┘  └────┬─────┘
     │              │
     ▼              ▼
┌─────────┐   ┌─────────┐
│ users_db│   │ posts_db│
│(Postgres│   │(Postgres│
└─────────┘   └─────────┘
```

---

## Struktur Proyek

```text
go-microservice/
├── api-gateway/
│   ├── middleware/
│   │   └── auth.go         # JWT validation logic
│   ├── main.go             # Routing & reverse proxy
│   ├── go.mod
│   └── Dockerfile
│
├── user-service/
│   ├── config/
│   │   └── database.go     # Koneksi PostgreSQL (users_db)
│   ├── handlers/
│   │   └── user_handler.go # Register, Login, Profile
│   ├── models/
│   │   └── user.go         # GORM model + DTOs
│   ├── repositories/
│   │   └── user_repository.go
│   ├── utils/
│   │   └── jwt.go          # Generate & validate token
│   ├── main.go
│   ├── go.mod
│   └── Dockerfile
│
├── post-service/
│   ├── config/
│   │   └── database.go     # Koneksi PostgreSQL (posts_db)
│   ├── handlers/
│   │   └── post_handler.go # CRUD Posts
│   ├── models/
│   │   └── post.go         # GORM model + DTOs
│   ├── repositories/
│   │   └── post_repository.go
│   ├── main.go
│   ├── go.mod
│   └── Dockerfile
│
└── docker-compose.yml      # Orkestrasi semua service
```

---

## Teknologi

| Komponen | Teknologi |
| :--- | :--- |
| Language | Go 1.21 |
| HTTP Router | `gorilla/mux` (user & post service) |
| ORM | GORM + PostgreSQL driver |
| Auth | JWT (`golang-jwt/jwt/v5`) + bcrypt |
| Env | `godotenv` |
| Container | Docker + Docker Compose |
| Database | PostgreSQL 15 (dua database terpisah) |

---

## Prasyarat

- [Docker](https://docs.docker.com/get-docker/) & [Docker Compose](https://docs.docker.com/compose/)

> Anda **tidak perlu** menginstal Go, PostgreSQL, atau apapun secara lokal. Semuanya berjalan di dalam Docker.

---

## Cara Menjalankan

```bash
# Masuk ke folder proyek
cd go-microservice

# Build & jalankan semua service
docker compose up --build

# Jalankan di background
docker compose up --build -d
```

Semua service siap di `http://localhost:8080`.

### Menghentikan Semua Service

```bash
docker compose down

# Hapus juga data database (volume)
docker compose down -v
```

---

## API Endpoints

Semua request dikirim ke **API Gateway** di port `8080`.

### 👤 User Service

| Method | Endpoint | Auth | Deskripsi |
| :--- | :--- | :---: | :--- |
| POST | `/api/users/register` | ❌ | Daftar akun baru |
| POST | `/api/users/login` | ❌ | Login & dapatkan JWT |
| GET | `/api/users/profile` | ✅ | Lihat profil sendiri |

### 📝 Post Service

| Method | Endpoint | Auth | Deskripsi |
| :--- | :--- | :---: | :--- |
| GET | `/api/posts` | ❌ | Ambil semua post |
| GET | `/api/posts/{id}` | ❌ | Ambil detail post |
| POST | `/api/posts` | ✅ | Buat post baru |
| PUT | `/api/posts/{id}` | ✅ | Update post |
| DELETE | `/api/posts/{id}` | ✅ | Hapus post |

> ✅ = Wajib menyertakan header `Authorization: Bearer <token>`

---

## Contoh Penggunaan (cURL)

### 1. Register

```bash
curl -X POST http://localhost:8080/api/users/register \
  -H "Content-Type: application/json" \
  -d '{"name": "Hary", "email": "hary@example.com", "password": "secret123"}'
```

### 2. Login → Simpan token

```bash
curl -X POST http://localhost:8080/api/users/login \
  -H "Content-Type: application/json" \
  -d '{"email": "hary@example.com", "password": "secret123"}'
```

Response:
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR...",
  "user": { "id": 1, "name": "Hary", "email": "hary@example.com" }
}
```

### 3. Buat Post (gunakan token dari step 2)

```bash
curl -X POST http://localhost:8080/api/posts \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR..." \
  -d '{"title": "Post Pertama", "content": "Isi konten post pertama saya."}'
```

### 4. Ambil Semua Post (tanpa token)

```bash
curl http://localhost:8080/api/posts
```

---

## Konsep Microservice yang Dipelajari

- **Service Isolation**: Setiap service punya codebase, database, dan container-nya sendiri.
- **API Gateway Pattern**: Satu pintu masuk untuk semua request dari client.
- **JWT Auth di Gateway**: Token divalidasi di Gateway, tidak perlu validasi ulang di setiap service.
- **Header Forwarding**: Gateway meneruskan `X-User-ID` ke downstream service setelah JWT valid.
- **Docker Networking**: Service berkomunikasi via nama container (bukan IP) dalam jaringan internal Docker.
- **Database per Service**: `users_db` dan `posts_db` terpisah — tidak boleh saling akses langsung.
