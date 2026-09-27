# DANQIVA

DANQIVA adalah project simulasi **Digital Wallet / Fintech** yang dibuat sebagai portfolio untuk menerapkan konsep backend engineering pada sistem transaksi.

Project ini berfokus pada bagaimana sebuah sistem wallet menangani authentication, transaksi, konsistensi saldo, concurrency, ledger, asynchronous processing, dan keamanan webhook.

> **Catatan:** DANQIVA merupakan project simulasi untuk kebutuhan portfolio dan pembelajaran. Project ini tidak memproses uang sungguhan dan tidak terhubung dengan payment gateway nyata.

---

## Tujuan Project

Project ini dibuat untuk menunjukkan penerapan beberapa konsep backend yang umum digunakan pada sistem transaksi, khususnya:

* Authentication & authorization
* Database transaction
* Concurrency control
* Row-level locking
* Idempotency
* Transaction ledger
* Background processing
* Outbox Pattern
* Webhook security
* Automated testing

Fokus utama project bukan sekadar membuat fitur CRUD, tetapi bagaimana menjaga **konsistensi dan keamanan data ketika transaksi terjadi**.

---

## Teknologi

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

## Fitur

### Authentication

* Registrasi
* Login
* JWT Access Token
* Refresh Token
* Logout
* Password hashing menggunakan bcrypt
* Protected API

### Wallet

* Melihat saldo
* Pengelolaan wallet
* Validasi kepemilikan wallet
* Update saldo dengan database transaction

### Top Up

User dapat melakukan simulasi penambahan saldo ke wallet.

Proses transaksi menggunakan database transaction untuk menjaga konsistensi data.

### Transfer

User dapat melakukan transfer saldo antar pengguna.

Proses transfer mencakup:

* Validasi sender dan receiver
* Validasi saldo
* Database transaction
* Row-level locking
* Concurrency handling
* Transaction record
* Ledger entry
* Rollback ketika transaksi gagal

---

## Concurrency Control

Salah satu fokus utama DANQIVA adalah menangani transaksi yang terjadi secara bersamaan.

Contoh:

```text
Saldo: Rp100.000

Request A → Transfer Rp100.000
Request B → Transfer Rp100.000

        ↓

Database Row Lock

        ↓

Satu transaksi berhasil
Satu transaksi gagal

        ↓

Saldo akhir: Rp0
```

Mekanisme ini digunakan untuk mencegah dua request menggunakan saldo yang sama secara bersamaan.

---

## Transaction Ledger

DANQIVA menggunakan konsep **transaction dan ledger entry** untuk mencatat perubahan saldo.

Contoh:

```text
Transfer Rp50.000

Transaction
├── Sender
├── Receiver
├── Amount
├── Status
└── Reference

Ledger
├── Sender   -Rp50.000
└── Receiver +Rp50.000
```

Ledger digunakan sebagai bagian dari pencatatan aktivitas transaksi dan perubahan saldo.

---

## Idempotency

Project menerapkan konsep **idempotency** untuk membantu mencegah request transaksi yang sama diproses lebih dari satu kali.

Hal ini relevan pada sistem transaksi karena request dapat mengalami retry akibat network timeout atau gangguan komunikasi.

---

## Asynchronous Processing

DANQIVA menggunakan Redis dan Asynq untuk mendukung proses asynchronous dan background job.

Konsep sederhananya:

```text
Application
     │
     ▼
Create Event
     │
     ▼
Queue
     │
     ▼
Worker
     │
     ▼
Background Processing
```

Pendekatan ini digunakan untuk proses yang tidak harus dilakukan secara synchronous pada request utama.

---

## Outbox Pattern

DANQIVA juga menerapkan **Outbox Pattern** untuk membantu menjaga konsistensi antara database transaction dan event asynchronous.

```text
Database Transaction
        │
        ├── Business Data
        │
        └── Outbox Event
                 │
                 ▼
               Worker
```

---

## Webhook Security

Webhook menggunakan **HMAC SHA-256** untuk melakukan validasi signature.

```text
Request
   │
   ▼
Verify Signature
   │
   ├── Invalid → Reject
   │
   └── Valid   → Process
```

Tujuannya adalah memastikan request webhook dapat diverifikasi sebelum diproses.

---

## Testing

Testing difokuskan pada bagian yang memiliki risiko terhadap konsistensi transaksi.

Beberapa skenario yang diuji:

* Successful transfer
* Insufficient balance
* Database rollback
* Concurrent transfer
* Ledger entry

Contoh skenario concurrency:

```text
Concurrent Request : 2
Berhasil            : 1
Gagal               : 1
```

---

## Arsitektur

Backend menggunakan pendekatan layered architecture:

```text
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

Pendekatan ini digunakan untuk memisahkan tanggung jawab antara HTTP layer, business logic, dan data access.

---

## Struktur Project

```text
DANQIVA/
│
├── backend/
│
├── frontend/
│
├── docker-compose.yml
│
├── LICENSE
│
└── README.md
```

Detail implementasi tersedia langsung pada source code repository.

---

## Catatan

Repository ini dipublikasikan sebagai **portfolio source code** untuk menunjukkan implementasi dan pendekatan engineering yang digunakan dalam pengembangan DANQIVA.

Source code tidak dimaksudkan sebagai template atau starter project untuk digunakan kembali pada project lain.
