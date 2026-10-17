# Sistema de Pedidos - Microservicios en Go

> Arquitectura moderna de microservicios con Go + PostgreSQL. Sistema escalable de gestión de pedidos con 4 servicios independientes y una tienda en React.

## 🎯 Descripción

Sistema que implementa **4 microservicios desacoplados** (Users, Products, Orders, Payments), cada uno con su propia base de datos PostgreSQL, y una tienda web que los consume. Demuestra patrones de arquitectura limpia, DDD y API REST.

## 🚀 Características

**Backend:**
- ✅ Arquitectura limpia con separación clara (Handler → Service → Repository)
- ✅ Microservicios independientes con bases de datos desacopladas
- ✅ Autenticación JWT
- ✅ GORM ORM + PostgreSQL
- ✅ Logging estructurado con slog
- ✅ CORS configurado
- ✅ Graceful shutdown
- ✅ Pedidos con transiciones de estado validadas e historial
- ✅ Pagos con idempotencia (la misma clave no cobra dos veces)

**Frontend (tienda):**
- ✅ React + TypeScript + Vite
- ✅ Catálogo con búsqueda y carrito persistente
- ✅ Checkout: crea el pedido y lo paga
- ✅ Mis pedidos (con reintento de pago) y mi cuenta
- ✅ Menú de usuario con avatar

**Infraestructura:**
- ✅ Docker Compose: levanta base de datos, backend, frontend y datos de demo con un comando

## 🛠️ Stack

| Componente | Tech |
|-----------|------|
| **Backend** | Go 1.26, Gin 1.12, GORM 1.31 |
| **Base de Datos** | PostgreSQL 16 (4 bases lógicas) |
| **Auth** | JWT v5 |
| **Frontend** | React, TypeScript, Vite, Zustand, nginx |
| **Contenedores** | Docker, Docker Compose |

## 🏗️ Arquitectura

```
┌─────────────────┐
│   Frontend      │
│  (React + TS)   │
└────────┬────────┘
         │ HTTP/REST
┌────────▼────────────────────────────────┐
│         API Gateway (Gin)                │
├────┬──────────┬──────────┬──────────┐
│    │          │          │          │
▼    ▼          ▼          ▼          ▼
Users  Products  Orders  Payments
```

**4 Microservicios independientes:**
- **Users**: Autenticación y gestión de usuarios
- **Products**: Catálogo y control de stock
- **Orders**: Gestión de pedidos y historial
- **Payments**: Procesamiento de pagos con idempotencia

## ⚡ Puesta en marcha

### Forma rápida: solo Docker (recomendada)

No hace falta instalar PostgreSQL, Go ni Node. Solo necesitás **Docker Desktop** (o Docker Engine con el plugin de Compose).

```bash
git clone <URL_DE_TU_REPO>
cd SistemaDePedidos
docker compose --profile db up --build
```

Eso levanta:
- **PostgreSQL** con las 4 bases creadas.
- **Backend**, con las tablas creadas automáticamente.
- **Seed**: carga 6 productos de demo y termina solo.
- **Frontend**: la tienda.

Después abrí:

- 🛍️ **Tienda:** http://localhost:3000
- ⚙️ **API:** http://localhost:8080/health

Para probar: registrate con cualquier email, agregá productos al carrito y finalizá la compra. El pago es simulado (90% aprobado, 10% rechazado); si se rechaza, el pedido queda pendiente y podés reintentar desde "Mis pedidos".

Para apagar:
```bash
docker compose --profile db down        # apaga los contenedores
docker compose --profile db down -v     # además borra la base de datos
```

> Los primeros `up --build` tardan varios minutos porque descargan las imágenes y compilan el backend.

### Forma con tu PostgreSQL local

Útil si ya tenés PostgreSQL instalado y querés usar tus propias bases.

**1. Creá las 4 bases** (una sola vez). Ejecutá `docker/postgres/create-databases.sql` conectado a la base `postgres`:
```bash
psql -U postgres -h localhost -f docker/postgres/create-databases.sql
```
Crea `users_db`, `products_db`, `orders_db` y `payments_db`.

**2. Configurá el entorno:**
```bash
cp .env.example .env
```
Editá `.env` y poné tu `DB_PASSWORD`.

**3. Levantá el backend en Docker contra tu Postgres:**

La primera vez, para crear las tablas:
```bash
RUN_MIGRATIONS=true docker compose --profile local up --build
```
Cuando veas `migraciones aplicadas` y `servidor arrancando`, cortá con Ctrl+C. Las siguientes veces:
```bash
docker compose --profile local up --build
```

En Linux, tu Postgres tiene que aceptar conexiones desde Docker: en `postgresql.conf` poné `listen_addresses = '*'` y en `pg_hba.conf` agregá `host all all 172.16.0.0/12 scram-sha-256`, y reiniciá Postgres. En Windows y macOS normalmente funciona sin cambios.

**4. Frontend** (con Docker, ya incluido en el perfil `local`) o en modo desarrollo:
```bash
cd frontend
npm ci
npm run dev        # http://localhost:5173
```

### Forma sin Docker (desarrollo)

Necesitás Go 1.26, Node 22 y PostgreSQL.

```bash
# Backend (primera vez, con RUN_MIGRATIONS=true)
cd backend
RUN_MIGRATIONS=true go run ./cmd/api

# Frontend, en otra terminal
cd frontend
npm ci
npm run dev
```

### Puertos

| Servicio | Puerto en tu PC | Para qué |
|---|---|---|
| Tienda (frontend en Docker) | **3000** | http://localhost:3000 |
| Backend (API) | **8080** | http://localhost:8080 |
| Frontend en modo desarrollo | 5173 | http://localhost:5173 |
| PostgreSQL en contenedor (perfil `db`) | 5433 | Conexión desde pgAdmin, opcional |

> Si ya tenés el backend corriendo en el puerto 8080, cortalo antes de levantar Docker.

## 📚 Modelos de Datos

### Users Service
![Diagrama ER - Users](images/1.2.JPG)
- Autenticación JWT
- Roles (`customer`, `admin`)
- Contraseñas hasheadas

### Products Service
![Diagrama ER - Products](images/1.1.JPG)
- Catálogo de productos y categorías
- Control de stock
- Reservas de inventario

### Orders Service
![Diagrama ER - Orders](images/1.3.JPG)
- Estados: pending → paid → processing → shipped → delivered (con fallos y cancelación)
- Historial completo de cambios
- Precios congelados al momento de la compra

### Payments Service
![Diagrama ER - Payments](images/1.4.JPG)
- Idempotencia por clave
- Estados: pending, succeeded, failed
- Pagos simulados para la demo

## 🔌 API Endpoints

Todas las rutas van bajo `/api/v1`. Las marcadas con 🔒 requieren el header `Authorization: Bearer <token>`.

### Users
```bash
POST   /auth/register          # {email, password, full_name}
POST   /auth/login             # {email, password} -> {token, user}
GET    /auth/me                # 🔒
```

### Products
```bash
GET    /products?page=1&limit=20&category_id=<uuid>&search=remera&min_price=0&max_price=5000&sort=price_asc|price_desc|newest
GET    /products/:id
```

### Orders
```bash
POST   /orders                 # 🔒 {items: [{product_id, quantity}], idempotency_key}
GET    /orders                 # 🔒 pedidos del usuario
GET    /orders/:id             # 🔒
PATCH  /orders/:id/status      # 🔒 {status}
```

### Payments
```bash
POST   /payments               # 🔒 {order_id, idempotency_key}
GET    /payments/:id           # 🔒
GET    /payments/order/:orderID  # 🔒
```

## ⚙️ Configuración

Variables principales (ver `.env.example`):

```env
APP_PORT=8080
APP_ENV=development

DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_SSLMODE=disable

USERS_DB_NAME=users_db
PRODUCTS_DB_NAME=products_db
ORDERS_DB_NAME=orders_db
PAYMENTS_DB_NAME=payments_db

JWT_SECRET=your_secret_key
JWT_EXPIRATION_HOURS=24

# Solo para Docker: dónde está tu Postgres desde el contenedor
DOCKER_DB_HOST=host.docker.internal
```

`RUN_MIGRATIONS=true` crea o completa las tablas con los modelos del código. Está apagado por defecto, así no toca una base que ya tiene datos.

## 🧪 Desarrollo

```bash
cd backend

# Tests
( -- estar dentro de la carpeta de backend --)
go test ./tests/ -v

# Lint
go fmt ./...
go vet ./...
```

## 🔐 Seguridad

✅ Implementado:
- JWT para autenticación
- Hash de contraseñas
- CORS configurado
- Prevención de SQL injection (GORM)
- Logging de requests
- Pagos idempotentes

⚠️ Pendiente / producción:
- Restringir `PATCH /orders/:id/status` a administradores
- Secretos fuera del código (`JWT_SECRET` y contraseñas)
- HTTPS/TLS
- Rate limiting
- Protección contra DDoS
- Gestión de secretos (Vault)

## 📈 Escalabilidad

**Horizontal:** cada servicio en una instancia separada
**Vertical:** replicación de PostgreSQL y caché con Redis
**Futuro:** event-driven con Kafka, service mesh y monitoreo

## 🎯 Aprendizajes Clave

Este proyecto demuestra:
- ✅ Arquitectura limpia en Go
- ✅ Microservicios con DDD
- ✅ API REST
- ✅ Manejo de errores
- ✅ Logging estructurado
- ✅ Composición de dependencias
- ✅ Integración frontend-backend
- ✅ Contenedores con Docker Compose

## 📝 Licencia

MIT

## 👨‍💻 Autor

**Rodrigo Fagundez**
GitHub: [@RodrigoFagundez2004](https://github.com/RodrigoFagundez2004)

---

**Estado:** ✅ Funcional con Docker