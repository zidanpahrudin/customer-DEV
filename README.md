# Supplier Management API

API untuk manajemen supplier menggunakan Golang, Gin, dan PostgreSQL.

## Persyaratan

- Go 1.21 atau lebih tinggi
- PostgreSQL
- Git

## Instalasi

1. Clone repository:
```bash
git clone <repository-url>
cd supplier-api
```

2. Install dependencies:
```bash
go mod tidy
```

3. Setup database PostgreSQL dan sesuaikan konfigurasi di file `.env`:
```env
DB_HOST=localhost
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=supplier_db
DB_PORT=5432
JWT_SECRET=your_jwt_secret_key
```

4. Jalankan aplikasi:
```bash
go run main.go
```

## Endpoint API

### Autentikasi

#### Register User
- **POST** `/register`
```json
{
    "username": "user123",
    "password": "password123"
}
```

#### Login
- **POST** `/login`
```json
{
    "username": "user123",
    "password": "password123"
}
```

### Supplier Management

Semua endpoint supplier memerlukan token JWT di header Authorization: `Bearer <token>`

#### Create Supplier
- **POST** `/api/suppliers`
```json
{
    "code": "SUPP001",
    "name": "PT Setroom Indonesia",
    "address": "Jakarta, Indonesia",
    "contact": "Albert Einstein",
    "status": "Active",
    "average_cost": 320000000
}
```

#### Get All Suppliers
- **GET** `/api/suppliers`
- **GET** `/api/suppliers?status=Active` (filter by status)

#### Get Supplier by ID
- **GET** `/api/suppliers/:id`

#### Update Supplier
- **PUT** `/api/suppliers/:id`
```json
{
    "name": "PT Setroom Indonesia Updated",
    "address": "Jakarta Selatan, Indonesia",
    "contact": "Albert Einstein",
    "status": "Active",
    "average_cost": 350000000
}
```

#### Delete Supplier
- **DELETE** `/api/suppliers/:id`

## Response Format

### Success Response
```json
{
    "suppliers": [...],
    "stats": {
        "total_suppliers": 1869,
        "new_suppliers": 27,
        "avg_cost": 320300000,
        "blocked_suppliers": 31
    }
}
```

### Error Response
```json
{
    "error": "Error message"
}
```