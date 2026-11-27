#!/usr/bin/env bash
set -e

echo "======================================================"
echo " PRUEBAS DE ENDPOINTS - ARQUITECTURA DE MICROSERVICIOS "
echo "======================================================"

# 1. Users API (Python :8000)
echo -e "\n>>> [1/4] Probando Users API (Python :8000)..."
curl -s http://localhost:8000/health
echo ""
echo "[POST] Crear Usuario:"
USER_RES=$(curl -s -X POST http://localhost:8000/users \
  -H "Content-Type: application/json" \
  -d '{"name":"Nicolas Alvear","email":"nicolas@example.com","password":"secretPassword123"}')
echo "$USER_RES"

echo "[GET] Listar Usuarios:"
curl -s http://localhost:8000/users
echo ""

# 2. Auth API (Java :8001)
echo -e "\n>>> [2/4] Probando Auth API (Java :8001)..."
curl -s http://localhost:8001/health
echo ""
echo "[POST] Login Credenciales Válidas (admin / admin123):"
curl -s -X POST http://localhost:8001/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}'
echo ""
echo "[POST] Login Credenciales Inválidas:"
curl -s -X POST http://localhost:8001/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"wrongpassword"}'
echo ""

# 3. Orders API (Go :8002)
echo -e "\n>>> [3/4] Probando Orders API (Go :8002)..."
curl -s http://localhost:8002/health
echo ""
echo "[POST] Crear Pedido:"
ORDER_RES=$(curl -s -X POST http://localhost:8002/orders \
  -H "Content-Type: application/json" \
  -d '{"customer_name":"Cliente Cloud","product":"VPS 4GB RAM","quantity":2,"total":50.00,"status":"PENDING"}')
echo "$ORDER_RES"
echo "[GET] Listar Pedidos:"
curl -s http://localhost:8002/orders
echo ""

# 4. Products API (JavaScript :8003)
echo -e "\n>>> [4/4] Probando Products API (JavaScript :8003)..."
curl -s http://localhost:8003/health
echo ""
echo "[POST] Registrar Producto:"
PROD_RES=$(curl -s -X POST http://localhost:8003/products \
  -H "Content-Type: application/json" \
  -d '{"name":"Monitor 27 Pulgadas","price":249.99,"description":"4K UHD IPS 144Hz","stock":25}')
echo "$PROD_RES"
echo "[GET] Listar Productos:"
curl -s http://localhost:8003/products
echo ""

echo "======================================================"
echo " PRUEBAS COMPLETADAS "
echo "======================================================"
