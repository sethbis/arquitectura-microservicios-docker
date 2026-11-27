# 🚀 Arquitectura de Microservicios con APIs Multilenguaje, SQLite y Docker

Solución desarrollada para la materia de **Arquitectura en la Nube**. Implementa cuatro microservicios completamente independientes, desarrollados en diferentes lenguajes de programación (**Python, Java, Go y JavaScript**), cada uno con su propio contenedor Docker, su propia base de datos relacional SQLite y persistencia garantizada mediante volúmenes Docker independientes gestionados por Docker Compose.

---

## 🏛️ Arquitectura del Sistema

```
                    ┌───────────────────────────────────────────┐
                    │          Docker Compose Orquestador       │
                    └─────────────────────┬─────────────────────┘
                                          │
       ┌──────────────────┬───────────────┴───────────────┬──────────────────┐
       │                  │                               │                  │
       ▼                  ▼                               ▼                  ▼
┌──────────────┐   ┌──────────────┐                ┌──────────────┐   ┌──────────────┐
│  Users API   │   │   Auth API   │                │  Orders API  │   │ Products API │
│ Python:8000  │   │  Java:8001   │                │   Go:8002    │   │  Node:8003   │
└──────┬───────┘   └──────┬───────┘                └──────┬───────┘   └──────┬───────┘
       │                  │                               │                  │
       ▼                  ▼                               ▼                  ▼
┌──────────────┐   ┌──────────────┐                ┌──────────────┐   ┌──────────────┐
│ users_data   │   │  auth_data   │                │ orders_data  │   │products_data │
│  (users.db)  │   │  (auth.db)   │                │ (orders.db)  │   │(products.db) │
└──────────────┘   └──────────────┘                └──────────────┘   └──────────────┘
```

### Tabla de Microservicios

| Microservicio | Lenguaje / Framework | Responsabilidad | Puerto | Base de Datos | Volumen Docker | Dependencias |
|---|---|---|---|---|---|---|
| **Users API** | Python 3.11 / Flask | Gestión y CRUD de usuarios | `8000` | `users.db` | `users_sqlite_data` | `requirements.txt` |
| **Auth API** | Java 17 / Spring Boot | Autenticación, login y logs | `8001` | `auth.db` | `auth_sqlite_data` | `pom.xml` |
| **Orders API** | Go 1.22 / net/http | Gestión y CRUD de pedidos | `8002` | `orders.db` | `orders_sqlite_data` | `go.mod` |
| **Products API** | JavaScript / Express | Gestión y CRUD de productos | `8003` | `products.db` | `products_sqlite_data` | `package.json` |

---

## 📂 Estructura del Repositorio

```text
.
├── users-api/              # Microservicio en Python
│   ├── app.py              # Endpoints CRUD y conexión SQLite
│   ├── requirements.txt    # Dependencias (Flask)
│   └── Dockerfile          # Imagen Python 3.11-slim
├── auth-api/               # Microservicio en Java
│   ├── src/                # Código fuente Spring Boot
│   ├── pom.xml             # Dependencias Maven (Spring Boot + SQLite JDBC)
│   └── Dockerfile          # Multi-stage build (Maven + JRE 17 Alpine)
├── orders-api/             # Microservicio en Go
│   ├── main.go             # Endpoints CRUD y SQLite puro (sin CGO)
│   ├── go.mod / go.sum     # Módulos y dependencias de Go
│   └── Dockerfile          # Multi-stage build (Go Alpine -> Alpine mínimo)
├── products-api/           # Microservicio en JavaScript (Node.js)
│   ├── index.js            # Endpoints CRUD Express y SQLite
│   ├── package.json        # Dependencias (Express + SQLite3)
│   └── Dockerfile          # Imagen Node 20 Alpine
├── docker-compose.yml      # Orquestación de los 4 servicios y 4 volúmenes
├── test-endpoints.sh       # Script bash de pruebas automatizadas
├── test-endpoints.ps1      # Script PowerShell de pruebas automatizadas
└── README.md               # Documentación completa y guía de uso
```

---

## ⚡ Ejecución Rápida con Docker Compose

### 1. Clonar el repositorio
```bash
git clone https://github.com/sethbis/arquitectura-microservicios-docker.git
cd arquitectura-microservicios-docker
```

### 2. Construir y levantar todos los microservicios
```bash
docker compose up -d --build
```

### 3. Verificar que los contenedores están en ejecución
```bash
docker compose ps
```
Salida esperada (los 4 contenedores en estado `Up` / `running`):
```text
NAME               IMAGE                                     COMMAND                  SERVICE        CREATED         STATUS         PORTS
auth-service       arquitectura-microservicios-auth-api      "java -jar app.jar"      auth-api       1 minute ago    Up 1 minute    0.0.0.0:8001->8001/tcp
orders-service     arquitectura-microservicios-orders-api    "./orders-api"           orders-api     1 minute ago    Up 1 minute    0.0.0.0:8002->8002/tcp
products-service   arquitectura-microservicios-products-api  "docker-entrypoint.s…"   products-api   1 minute ago    Up 1 minute    0.0.0.0:8003->8003/tcp
users-service      arquitectura-microservicios-users-api     "python app.py"          users-api      1 minute ago    Up 1 minute    0.0.0.0:8000->8000/tcp
```

### 4. Ver los logs en tiempo real
```bash
docker compose logs -f
```

---

## 🧪 Pruebas de Endpoints (Guía Completa con `curl`)

### 1. Users API — Python (Puerto 8000)

* **Health Check:**
  ```bash
  curl -X GET http://localhost:8000/health
  ```
* **Registrar Usuario (POST):**
  ```bash
  curl -X POST http://localhost:8000/users \
    -H "Content-Type: application/json" \
    -d '{"name": "Carlos Gomez", "email": "carlos@example.com", "password": "securePass123"}'
  ```
* **Consultar Lista de Usuarios (GET):**
  ```bash
  curl -X GET http://localhost:8000/users
  ```
* **Consultar Usuario por ID (GET):**
  ```bash
  curl -X GET http://localhost:8000/users/1
  ```
* **Actualizar Información de Usuario (PUT):**
  ```bash
  curl -X PUT http://localhost:8000/users/1 \
    -H "Content-Type: application/json" \
    -d '{"name": "Carlos Gomez Modificado", "email": "carlos_mod@example.com", "password": "newPassword987"}'
  ```
* **Eliminar Usuario (DELETE):**
  ```bash
  curl -X DELETE http://localhost:8000/users/1
  ```

---

### 2. Auth API — Java (Puerto 8001)

* **Health Check:**
  ```bash
  curl -X GET http://localhost:8001/health
  ```
* **Login con Credenciales Válidas (POST):** *(Usuarios pre-poblados por defecto: `admin` / `admin123`, `user1` / `password123`)*
  ```bash
  curl -X POST http://localhost:8001/auth/login \
    -H "Content-Type: application/json" \
    -d '{"username": "admin", "password": "admin123"}'
  ```
  *Respuesta (200 OK):*
  ```json
  {"valid": true, "message": "Credenciales válidas. Inicio de sesión exitoso.", "username": "admin"}
  ```
* **Login con Credenciales Inválidas (POST):**
  ```bash
  curl -X POST http://localhost:8001/auth/login \
    -H "Content-Type: application/json" \
    -d '{"username": "admin", "password": "passwordIncorrecto"}'
  ```
  *Respuesta (401 Unauthorized):*
  ```json
  {"valid": false, "message": "Credenciales inválidas. Usuario o contraseña incorrectos."}
  ```
* **Registrar Nuevo Usuario para Auth (POST):**
  ```bash
  curl -X POST http://localhost:8001/auth/register \
    -H "Content-Type: application/json" \
    -d '{"username": "devuser", "password": "mypassword"}'
  ```
* **Consultar Logs de Intentos de Autenticación (GET):**
  ```bash
  curl -X GET http://localhost:8001/auth/logs
  ```

---

### 3. Orders API — Go (Puerto 8002)

* **Health Check:**
  ```bash
  curl -X GET http://localhost:8002/health
  ```
* **Crear Pedido (POST):**
  ```bash
  curl -X POST http://localhost:8002/orders \
    -H "Content-Type: application/json" \
    -d '{"customer_name": "Empresa Cloud S.A.", "product": "Cluster Kubernetes", "quantity": 3, "total": 450.00, "status": "PENDING"}'
  ```
* **Consultar Todos los Pedidos (GET):**
  ```bash
  curl -X GET http://localhost:8002/orders
  ```
* **Consultar Pedido por ID (GET):**
  ```bash
  curl -X GET http://localhost:8002/orders/1
  ```
* **Actualizar Estado de Pedido (PUT):**
  ```bash
  curl -X PUT http://localhost:8002/orders/1 \
    -H "Content-Type: application/json" \
    -d '{"status": "COMPLETED"}'
  ```
* **Eliminar Pedido (DELETE):**
  ```bash
  curl -X DELETE http://localhost:8002/orders/1
  ```

---

### 4. Products API — JavaScript (Puerto 8003)

* **Health Check:**
  ```bash
  curl -X GET http://localhost:8003/health
  ```
* **Registrar Producto (POST):**
  ```bash
  curl -X POST http://localhost:8003/products \
    -H "Content-Type: application/json" \
    -d '{"name": "Teclado Mecanico Inalambrico", "price": 89.99, "description": "Switches Red Hot-Swap RGB", "stock": 40}'
  ```
* **Consultar Lista de Productos (GET):**
  ```bash
  curl -X GET http://localhost:8003/products
  ```
* **Consultar Producto por ID (GET):**
  ```bash
  curl -X GET http://localhost:8003/products/1
  ```
* **Actualizar Información de Producto (PUT):**
  ```bash
  curl -X PUT http://localhost:8003/products/1 \
    -H "Content-Type: application/json" \
    -d '{"name": "Teclado Mecanico Custom", "price": 99.99, "description": "Switches Gateron Lubed", "stock": 35}'
  ```
* **Eliminar Producto (DELETE):**
  ```bash
  curl -X DELETE http://localhost:8003/products/1
  ```

---

## 💾 Demostración y Validación de Persistencia SQLite

Para demostrar que los datos **NO se pierden** al detener, eliminar o recrear los contenedores:

1. **Crear registros en las APIs** (ejemplo: un usuario y un producto usando los curls anteriores).
2. **Detener y eliminar los contenedores:**
   ```bash
   docker compose down
   ```
   *(Nota: no usar la bandera `-v` para no borrar los volúmenes).*
3. **Volver a levantar los servicios:**
   ```bash
   docker compose up -d
   ```
4. **Consultar nuevamente los datos:**
   ```bash
   curl -X GET http://localhost:8000/users
   curl -X GET http://localhost:8003/products
   ```
   ✅ Los registros creados siguen intactos, comprobando la persistencia y el aislamiento de cada base de datos SQLite en su volumen independiente.

---

## 🛑 Detener y Limpiar la Arquitectura

* **Detener los servicios sin perder los datos:**
  ```bash
  docker compose stop
  ```
* **Detener y remover contenedores preservando volúmenes:**
  ```bash
  docker compose down
  ```
* **Detener y eliminar contenedores junto con los volúmenes (reinicio total de fábrica):**
  ```bash
  docker compose down -v
  ```

---

## 🌿 Estrategia de Ramas Git Utilizada

Siguiendo las mejores prácticas de integración continua y control de versiones, el desarrollo se organizó en ramas descriptivas de funcionalidad (`feat/...`):

* `feat/users-api`: Implementación del microservicio en Python (Flask), SQLite `users.db` y Dockerfile.
* `feat/auth-api`: Implementación del microservicio en Java (Spring Boot), SQLite `auth.db` y Dockerfile.
* `feat/orders-api`: Implementación del microservicio en Go, SQLite `orders.db` y Dockerfile.
* `feat/products-api`: Implementación del microservicio en JavaScript (Node/Express), SQLite `products.db` y Dockerfile.
* `feat/docker-compose`: Configuración del orquestador `docker-compose.yml`, volúmenes persistentes y scripts de verificación.
* `main`: Rama principal de producción consolidada con todas las características validadas.
