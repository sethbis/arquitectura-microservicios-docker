# Arquitectura de Microservicios con APIs Multilenguaje, SQLite y Docker Compose

Este proyecto implementa una arquitectura basada en cuatro microservicios independientes, desarrollados en diferentes lenguajes de programación, cada uno con su propia base de datos SQLite y orquestados mediante Docker Compose.

---

## 🏛️ Arquitectura del Sistema

| Microservicio | Lenguaje | Responsabilidad | Puerto | Base de Datos |
|---|---|---|---|---|
| **Users API** | Python (Flask) | CRUD de usuarios | `:8000` | `users.db` |
| **Auth API** | Java (Spring Boot) | Autenticación y login | `:8001` | `auth.db` |
| **Orders API** | Go (Golang) | CRUD de pedidos | `:8002` | `orders.db` |
| **Products API** | JavaScript (Node.js/Express) | CRUD de productos | `:8003` | `products.db` |

---

## 🚀 Despliegue Rápido con Docker Compose

```bash
# Clonar repositorio
git clone https://github.com/sethbis/arquitectura-microservicios-docker.git
cd arquitectura-microservicios-docker

# Construir y levantar todos los microservicios en segundo plano
docker compose up -d --build

# Verificar estado de los contenedores
docker compose ps
```
