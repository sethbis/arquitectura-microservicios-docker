# Script de pruebas automatizadas para validar los 4 microservicios
Write-Host "======================================================" -ForegroundColor Cyan
Write-Host " PRUEBAS DE ENDPOINTS - ARQUITECTURA DE MICROSERVICIOS " -ForegroundColor Cyan
Write-Host "======================================================" -ForegroundColor Cyan

# 1. Users API (:8000)
Write-Host "`n>>> [1/4] Probando Users API (Python :8000)..." -ForegroundColor Yellow
try {
    Invoke-RestMethod -Uri "http://localhost:8000/health" -Method GET | ConvertTo-Json
    Write-Host "[POST] Crear Usuario:" -ForegroundColor Green
    $userBody = @{ name = "Nicolas Alvear"; email = "nicolas@example.com"; password = "secretPassword123" } | ConvertTo-Json
    $userRes = Invoke-RestMethod -Uri "http://localhost:8000/users" -Method POST -Body $userBody -ContentType "application/json"
    $userRes | ConvertTo-Json
    $userId = $userRes.user.id

    Write-Host "[GET] Listar Usuarios:" -ForegroundColor Green
    Invoke-RestMethod -Uri "http://localhost:8000/users" -Method GET | ConvertTo-Json

    Write-Host "[GET] Consultar Usuario por ID ($userId):" -ForegroundColor Green
    Invoke-RestMethod -Uri "http://localhost:8000/users/$userId" -Method GET | ConvertTo-Json

    Write-Host "[PUT] Actualizar Usuario ($userId):" -ForegroundColor Green
    $userUpdate = @{ name = "Nicolas A. Actualizado"; email = "nicolas_updated@example.com"; password = "newSecretPassword" } | ConvertTo-Json
    Invoke-RestMethod -Uri "http://localhost:8000/users/$userId" -Method PUT -Body $userUpdate -ContentType "application/json" | ConvertTo-Json
} catch {
    Write-Host "Error en Users API: $_" -ForegroundColor Red
}

# 2. Auth API (:8001)
Write-Host "`n>>> [2/4] Probando Auth API (Java :8001)..." -ForegroundColor Yellow
try {
    Invoke-RestMethod -Uri "http://localhost:8001/health" -Method GET | ConvertTo-Json

    Write-Host "[POST] Login Credenciales Válidas (admin / admin123):" -ForegroundColor Green
    $authOk = @{ username = "admin"; password = "admin123" } | ConvertTo-Json
    Invoke-RestMethod -Uri "http://localhost:8001/auth/login" -Method POST -Body $authOk -ContentType "application/json" | ConvertTo-Json

    Write-Host "[POST] Login Credenciales Inválidas (admin / wrongpass):" -ForegroundColor Green
    try {
        $authFail = @{ username = "admin"; password = "wrongpass" } | ConvertTo-Json
        Invoke-RestMethod -Uri "http://localhost:8001/auth/login" -Method POST -Body $authFail -ContentType "application/json"
    } catch {
        Write-Host "Respuesta 401 esperada recibida correctamente:" $_.Exception.Response.StatusCode -ForegroundColor Green
    }

    Write-Host "[GET] Logs de autenticación registrados:" -ForegroundColor Green
    Invoke-RestMethod -Uri "http://localhost:8001/auth/logs" -Method GET | ConvertTo-Json
} catch {
    Write-Host "Error en Auth API: $_" -ForegroundColor Red
}

# 3. Orders API (:8002)
Write-Host "`n>>> [3/4] Probando Orders API (Go :8002)..." -ForegroundColor Yellow
try {
    Invoke-RestMethod -Uri "http://localhost:8002/health" -Method GET | ConvertTo-Json

    Write-Host "[POST] Crear Pedido:" -ForegroundColor Green
    $orderBody = @{ customer_name = "Cliente Cloud"; product = "Servidor VPS 4GB"; quantity = 2; total = 79.99; status = "PENDING" } | ConvertTo-Json
    $orderRes = Invoke-RestMethod -Uri "http://localhost:8002/orders" -Method POST -Body $orderBody -ContentType "application/json"
    $orderRes | ConvertTo-Json
    $orderId = $orderRes.order.id

    Write-Host "[GET] Listar Todos los Pedidos:" -ForegroundColor Green
    Invoke-RestMethod -Uri "http://localhost:8002/orders" -Method GET | ConvertTo-Json

    Write-Host "[GET] Consultar Pedido por ID ($orderId):" -ForegroundColor Green
    Invoke-RestMethod -Uri "http://localhost:8002/orders/$orderId" -Method GET | ConvertTo-Json

    Write-Host "[PUT] Actualizar Estado de Pedido ($orderId) a COMPLETED:" -ForegroundColor Green
    $orderUpdate = @{ status = "COMPLETED" } | ConvertTo-Json
    Invoke-RestMethod -Uri "http://localhost:8002/orders/$orderId" -Method PUT -Body $orderUpdate -ContentType "application/json" | ConvertTo-Json
} catch {
    Write-Host "Error en Orders API: $_" -ForegroundColor Red
}

# 4. Products API (:8003)
Write-Host "`n>>> [4/4] Probando Products API (JavaScript :8003)..." -ForegroundColor Yellow
try {
    Invoke-RestMethod -Uri "http://localhost:8003/health" -Method GET | ConvertTo-Json

    Write-Host "[POST] Registrar Producto:" -ForegroundColor Green
    $prodBody = @{ name = "Laptop Gamer Pro"; price = 1299.99; description = "Intel Core i7 16GB RAM 1TB SSD"; stock = 15 } | ConvertTo-Json
    $prodRes = Invoke-RestMethod -Uri "http://localhost:8003/products" -Method POST -Body $prodBody -ContentType "application/json"
    $prodRes | ConvertTo-Json
    $prodId = $prodRes.product.id

    Write-Host "[GET] Listar Productos:" -ForegroundColor Green
    Invoke-RestMethod -Uri "http://localhost:8003/products" -Method GET | ConvertTo-Json

    Write-Host "[GET] Consultar Producto por ID ($prodId):" -ForegroundColor Green
    Invoke-RestMethod -Uri "http://localhost:8003/products/$prodId" -Method GET | ConvertTo-Json

    Write-Host "[PUT] Actualizar Información de Producto ($prodId):" -ForegroundColor Green
    $prodUpdate = @{ name = "Laptop Gamer Pro Max"; price = 1399.99; description = "Intel Core i9 32GB RAM 1TB SSD"; stock = 12 } | ConvertTo-Json
    Invoke-RestMethod -Uri "http://localhost:8003/products/$prodId" -Method PUT -Body $prodUpdate -ContentType "application/json" | ConvertTo-Json
} catch {
    Write-Host "Error en Products API: $_" -ForegroundColor Red
}

Write-Host "`n======================================================" -ForegroundColor Cyan
Write-Host " PRUEBAS COMPLETADAS EXITOSAMENTE " -ForegroundColor Cyan
Write-Host "======================================================" -ForegroundColor Cyan
