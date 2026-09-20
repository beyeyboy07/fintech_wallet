# Fintech Wallet Frontend

Frontend aplikasi Fintech Wallet untuk kebutuhan learning dan portfolio.

## Tech Stack

* React
* TypeScript
* Vite
* React Router
* Zustand
* Axios

## Features

* User registration
* Login & logout
* JWT authentication
* Access token refresh
* Refresh token rotation
* Protected routes
* Wallet balance
* Top up
* Transfer
* Transaction history
* Transaction pagination
* Transfer idempotency
* API error handling
* Responsive layout

## Requirements

* Node.js
* npm
* Backend Fintech Wallet berjalan

## Installation

Clone repository kemudian masuk ke folder frontend:

```bash
cd frontend
```

Install dependencies:

```bash
npm install
```

Buat file `.env`:

```env
VITE_API_BASE_URL=http://localhost:8081
```

## Development

Jalankan development server:

```bash
npm run dev
```

Frontend biasanya tersedia di:

```text
http://localhost:5173
```

## Build

Untuk production build:

```bash
npm run build
```

## Lint

Untuk menjalankan ESLint:

```bash
npm run lint
```

## API

Frontend menggunakan backend API dengan base URL yang dikonfigurasi melalui:

```env
VITE_API_BASE_URL
```

Default development:

```text
http://localhost:8081
```

## Authentication

Frontend menyimpan:

* access token
* refresh token

Session akan dipulihkan ketika aplikasi dibuka kembali. Jika access token expired dan refresh token masih valid, frontend akan meminta token baru melalui endpoint refresh.

> Project ini dibuat untuk kebutuhan learning dan portfolio. Bukan aplikasi financial/payment production dan tidak digunakan untuk transaksi uang nyata.
