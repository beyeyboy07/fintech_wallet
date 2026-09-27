# DANQIVA

DANQIVA adalah project simulasi **Digital Wallet / Fintech** yang dibuat untuk menerapkan konsep backend engineering yang umum digunakan pada sistem transaksi, seperti authentication, wallet transaction, concurrency control, idempotency, asynchronous processing, webhook, dan transaction ledger.

> **Catatan:** DANQIVA merupakan project portfolio dan simulasi untuk pembelajaran. Project ini tidak memproses uang sungguhan dan tidak terhubung dengan payment gateway nyata.

---

## Tentang Project

DANQIVA mensimulasikan sebuah platform digital wallet yang memungkinkan pengguna untuk:

* Registrasi dan login
* Mengelola saldo wallet
* Melakukan top up
* Melakukan transfer antar pengguna
* Melihat riwayat transaksi
* Melakukan refresh session
* Memproses event secara asynchronous
* Menangani webhook dengan validasi signature

Project ini dibuat dengan fokus pada penerapan konsep backend yang lebih dari sekadar CRUD, terutama dalam menjaga konsistensi data dan keamanan proses transaksi.

---

## Teknologi yang Digunakan

### Backend

* Go
* Gin
* GORM
* MySQL
* JWT
* bcrypt
* Redis
* Asynq
* REST API
* Swagger / OpenAPI
* Docker

### Frontend

* React
* TypeScript
* Vite
* React Router
* Zustand
* Axios
* CSS

---

## Arsitektur

Backend DANQIVA menggunakan pendekatan **layered architecture**:

```text id="x9q2kc"
Request
   │
   ▼
Handler
   │
   ▼
Service
   │
   ▼
Repository
   │
   ▼
Database
```

Selain layer utama tersebut, project juga menggunakan Redis dan background worker:

```text id="n2c8xw"
                    ┌─────────────┐
                    │   Client    │
                    └──────┬──────┘
                           │
                           ▼
                    ┌─────────────┐
                    │    Gin      │
                    │   Handler   │
                    └──────┬──────┘
                           │
                           ▼
                    ┌─────────────┐
                    │   Service   │
                    └──────┬──────┘
                           │
              ┌────────────┼────────────┐
              ▼            ▼            ▼
         Repository      Redis        Asynq
              │                         │
              ▼                         ▼
           MySQL                    Worker
```

---

## Fitur Utama

### Authentication

* Registrasi user
* Login
* JWT access token
* Refresh token
* Logout
* Password hashing menggunakan bcrypt
* Protected API

### Wallet

* Melihat saldo wallet
* Validasi kepemilikan wallet
* Update saldo
* Menjaga konsistensi saldo melalui database transaction

### Top Up

User dapat melakukan simulasi penambahan saldo ke wallet.

Proses top up menggunakan database transaction untuk menjaga konsistensi data.

### Transfer

User dapat melakukan transfer saldo ke user lain.

Proses transfer mencakup:

* Validasi sender
* Validasi receiver
* Validasi saldo
* Database transaction
* Row-level locking
* Pencegahan concurrent transfer
* Pembuatan transaction
* Pembuatan ledger entry
* Rollback ketika transaksi gagal

Contoh kasus concurrent transfer:

```text id="t1s5ae"
Saldo awal: Rp100.000

Request A → Transfer Rp100.000
Request B → Transfer Rp100.000

              ↓

        Database Row Lock

              ↓

      Hanya satu transaksi
          yang berhasil

              ↓

Saldo akhir: Rp0
```

Dengan mekanisme tersebut, dua request yang datang secara bersamaan tidak dapat menggunakan saldo yang sama.

---

## Transaction Ledger

DANQIVA menggunakan konsep **transaction dan ledger** untuk mencatat aktivitas transaksi dan perubahan saldo.

Contoh transfer Rp50.000:

```text id="n0v8kg"
Transaction
├── Sender
├── Receiver
├── Amount
├── Status
└── Reference

Ledger Entries
├── Sender   -Rp50.000
└── Receiver +Rp50.000
```

Pendekatan ini memberikan dasar pencatatan transaksi yang lebih terstruktur dibandingkan hanya mengubah nilai saldo.

---

## Idempotency

DANQIVA menerapkan konsep **idempotency** untuk membantu mencegah transaksi yang sama diproses lebih dari satu kali ketika request dikirim ulang.

Hal ini penting pada sistem transaksi karena request dapat mengalami:

* Network timeout
* Retry dari client
* Duplicate request
* Gangguan komunikasi antara service

Dengan idempotency, request yang sama dapat dikenali sehingga tidak menyebabkan transaksi ganda.

---

## Redis & Background Job

Redis digunakan sebagai infrastructure pendukung untuk proses asynchronous.

Asynq digunakan untuk menjalankan background job seperti:

* Pemrosesan event transaksi
* Webhook
* Retryable job
* Proses asynchronous lainnya

Alur sederhananya:

```text id="v5xj5m"
API Request
    │
    ▼
Create Transaction
    │
    ▼
Enqueue Job
    │
    ▼
Redis
    │
    ▼
Asynq Worker
    │
    ▼
Background Processing
```

Dengan pendekatan ini, proses yang tidak harus dilakukan secara synchronous dapat dipindahkan ke background worker.

---

## Outbox Pattern

DANQIVA juga menggunakan **Outbox Pattern** untuk membantu menjaga konsistensi antara database transaction dan asynchronous event.

Alurnya:

```text id="x5c3qd"
Database Transaction
        │
        ├── Business Data
        │
        └── Outbox Event
                 │
                 ▼
              Worker
                 │
                 ▼
        External Processing
```

Dengan pendekatan ini, event tidak langsung dikirim ke proses eksternal sebelum transaksi database berhasil.

---

## Webhook

Webhook menggunakan **HMAC SHA-256** untuk melakukan validasi signature.

Alur sederhananya:

```text id="j0s2ka"
Sender
   │
   │ Payload + Signature
   ▼
DANQIVA
   │
   ▼
Verify HMAC
   │
   ├── Invalid → Reject
   │
   └── Valid   → Process
```

Signature digunakan untuk membantu memastikan bahwa request berasal dari sumber yang dipercaya dan payload tidak mengalami perubahan.

---

## Testing

Backend memiliki test untuk beberapa skenario penting, terutama yang berkaitan dengan transaksi:

* Successful transfer
* Insufficient balance
* Database rollback
* Concurrent transfer
* Ledger entry

Salah satu skenario concurrency:

```text id="1l4jns"
Concurrent Request : 2
Transfer berhasil  : 1
Transfer gagal     : 1
```

Test tersebut memastikan bahwa dua transaksi yang menggunakan saldo yang sama tidak dapat keduanya berhasil ketika saldo hanya mencukupi untuk satu transaksi.

---

## Struktur Project

```text id="7x6p8k"
DANQIVA/
│
├── backend/
│   ├── cmd/
│   ├── config/
│   ├── internal/
│   │   ├── handler/
│   │   ├── middleware/
│   │   ├── model/
│   │   ├── repository/
│   │   └── service/
│   ├── worker/
│   └── ...
│
├── frontend/
│   ├── src/
│   │   ├── assets/
│   │   ├── components/
│   │   ├── config/
│   │   ├── hooks/
│   │   ├── layouts/
│   │   ├── pages/
│   │   ├── routes/
│   │   ├── services/
│   │   ├── stores/
│   │   ├── types/
│   │   └── utils/
│   └── ...
│
├── docker-compose.yml
└── README.md
```

---

## API

Beberapa endpoint utama yang tersedia:

```text id="y4v9pm"
POST   /api/v1/auth/register
POST   /api/v1/auth/login
POST   /api/v1/auth/refresh
POST   /api/v1/auth/logout
GET    /api/v1/auth/me

GET    /api/v1/wallet

POST   /api/v1/transactions/top-up
POST   /api/v1/transactions/transfer
GET    /api/v1/transactions
```

Dokumentasi API tersedia melalui Swagger / OpenAPI ketika backend dijalankan.

---

## Menjalankan Project

### Prasyarat

Pastikan sudah tersedia:

* Go
* Node.js
* MySQL
* Redis
* Docker & Docker Compose

---

### 1. Clone Repository

```bash id="r2n7cb"
git clone <repository-url>
cd DANQIVA
```

---

### 2. Konfigurasi Backend

Buat konfigurasi environment berdasarkan file environment yang tersedia di project.

Contoh konfigurasi:

```env id="j5s3ab"
APP_PORT=8080
APP_ENV=development

DB_HOST=localhost
DB_PORT=3306
DB_NAME=fintech_wallet

JWT_SECRET=your-secret-key
```

---

### 3. Menjalankan Infrastructure

```bash id="f7r4de"
docker compose up -d
```

Perintah tersebut menjalankan infrastructure yang dibutuhkan seperti MySQL dan Redis.

---

### 4. Menjalankan Backend

```bash id="q0m6sy"
cd backend
go run .
```

Backend akan berjalan menggunakan port yang telah dikonfigurasi.

---

### 5. Menjalankan Frontend

Buka terminal baru:

```bash id="c9w2xa"
cd frontend
npm install
npm run dev
```

Frontend akan berjalan menggunakan Vite development server.

---

## Fokus Pembelajaran

DANQIVA dibuat untuk menerapkan beberapa konsep yang umum digunakan pada sistem backend production-oriented:

* Layered architecture
* Authentication & authorization
* JWT & refresh token
* Database transaction
* Concurrency control
* Row-level locking
* Idempotency
* Transaction ledger
* Redis
* Background job
* Asynq
* Outbox Pattern
* Webhook
* HMAC SHA-256
* Automated testing
* Docker

---

## Tujuan Project

Project ini dibuat sebagai portfolio untuk menunjukkan kemampuan dalam membangun aplikasi backend menggunakan Go, khususnya pada sistem yang membutuhkan **data consistency, transaction safety, concurrency handling, dan asynchronous processing**.

---

## License

Project ini dibuat untuk kebutuhan portfolio, pembelajaran, dan demonstrasi.

Copyright (c) 2026 Budy Rachman

All rights reserved.

Perangkat lunak dan source code ini disediakan hanya untuk keperluan
melihat, mempelajari, dan mengevaluasi portfolio.

Dilarang menyalin, memodifikasi, mendistribusikan, mempublikasikan,
memberikan sublisensi, atau menggunakan perangkat lunak ini untuk
keperluan komersial tanpa izin tertulis terlebih dahulu dari pemilik
hak cipta.