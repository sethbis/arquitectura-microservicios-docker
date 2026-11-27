package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	_ "modernc.org/sqlite"
)

type Order struct {
	ID           int     `json:"id"`
	CustomerName string  `json:"customer_name"`
	Product      string  `json:"product"`
	Quantity     int     `json:"quantity"`
	Total        float64 `json:"total"`
	Status       string  `json:"status"`
	CreatedAt    string  `json:"created_at"`
}

var db *sql.DB

func initDB() {
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "/app/data/orders.db"
	}

	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		log.Fatalf("Error creating directory for db: %v", err)
	}

	var err error
	db, err = sql.Open("sqlite", dbPath)
	if err != nil {
		log.Fatalf("Error opening SQLite database: %v", err)
	}

	createTableQuery := `
	CREATE TABLE IF NOT EXISTS orders (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		customer_name TEXT NOT NULL,
		product TEXT NOT NULL,
		quantity INTEGER NOT NULL,
		total REAL NOT NULL,
		status TEXT NOT NULL DEFAULT 'PENDING',
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`

	if _, err := db.Exec(createTableQuery); err != nil {
		log.Fatalf("Error creating orders table: %v", err)
	}
	log.Println("Orders database initialized successfully at", dbPath)
}

func jsonResponse(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if data != nil {
		_ = json.NewEncoder(w).Encode(data)
	}
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	jsonResponse(w, http.StatusOK, map[string]string{
		"status":  "UP",
		"service": "orders-api",
	})
}

// Handler for /orders and /orders/{id}
func ordersHandler(w http.ResponseWriter, r *http.Request) {
	path := strings.Trim(r.URL.Path, "/")
	parts := strings.Split(path, "/")

	if len(parts) == 1 && parts[0] == "orders" {
		switch r.Method {
		case http.MethodGet:
			listOrders(w, r)
		case http.MethodPost:
			createOrder(w, r)
		default:
			jsonResponse(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
		}
		return
	}

	if len(parts) == 2 && parts[0] == "orders" {
		id, err := strconv.Atoi(parts[1])
		if err != nil {
			jsonResponse(w, http.StatusBadRequest, map[string]string{"error": "Invalid order ID"})
			return
		}

		switch r.Method {
		case http.MethodGet:
			getOrder(w, r, id)
		case http.MethodPut:
			updateOrder(w, r, id)
		case http.MethodDelete:
			deleteOrder(w, r, id)
		default:
			jsonResponse(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
		}
		return
	}

	jsonResponse(w, http.StatusNotFound, map[string]string{"error": "Endpoint not found"})
}

// 1. Crear pedidos
func createOrder(w http.ResponseWriter, r *http.Request) {
	var input struct {
		CustomerName string  `json:"customer_name"`
		Product      string  `json:"product"`
		Quantity     int     `json:"quantity"`
		Total        float64 `json:"total"`
		Status       string  `json:"status"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		jsonResponse(w, http.StatusBadRequest, map[string]string{"error": "Invalid JSON body"})
		return
	}

	if input.CustomerName == "" || input.Product == "" || input.Quantity <= 0 || input.Total <= 0 {
		jsonResponse(w, http.StatusBadRequest, map[string]string{"error": "customer_name, product, quantity and total are required"})
		return
	}

	if input.Status == "" {
		input.Status = "PENDING"
	}

	res, err := db.Exec(
		"INSERT INTO orders (customer_name, product, quantity, total, status) VALUES (?, ?, ?, ?, ?)",
		input.CustomerName, input.Product, input.Quantity, input.Total, input.Status,
	)
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	id, _ := res.LastInsertId()
	jsonResponse(w, http.StatusCreated, map[string]interface{}{
		"message": "Order created successfully",
		"order": map[string]interface{}{
			"id":            id,
			"customer_name": input.CustomerName,
			"product":       input.Product,
			"quantity":      input.Quantity,
			"total":         input.Total,
			"status":        input.Status,
		},
	})
}

// 2. Consultar todos los pedidos
func listOrders(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT id, customer_name, product, quantity, total, status, created_at FROM orders")
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	defer rows.Close()

	orders := []Order{}
	for rows.Next() {
		var o Order
		if err := rows.Scan(&o.ID, &o.CustomerName, &o.Product, &o.Quantity, &o.Total, &o.Status, &o.CreatedAt); err != nil {
			continue
		}
		orders = append(orders, o)
	}

	jsonResponse(w, http.StatusOK, orders)
}

// 3. Consultar un pedido por su ID
func getOrder(w http.ResponseWriter, r *http.Request, id int) {
	var o Order
	row := db.QueryRow("SELECT id, customer_name, product, quantity, total, status, created_at FROM orders WHERE id = ?", id)
	err := row.Scan(&o.ID, &o.CustomerName, &o.Product, &o.Quantity, &o.Total, &o.Status, &o.CreatedAt)
	if err == sql.ErrNoRows {
		jsonResponse(w, http.StatusNotFound, map[string]string{"error": "Order not found"})
		return
	} else if err != nil {
		jsonResponse(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	jsonResponse(w, http.StatusOK, o)
}

// 4. Actualizar el estado de un pedido (y opcionalmente otros datos)
func updateOrder(w http.ResponseWriter, r *http.Request, id int) {
	var input struct {
		CustomerName string   `json:"customer_name"`
		Product      string   `json:"product"`
		Quantity     int      `json:"quantity"`
		Total        float64  `json:"total"`
		Status       *string  `json:"status"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		jsonResponse(w, http.StatusBadRequest, map[string]string{"error": "Invalid JSON body"})
		return
	}

	// Verificar si existe el pedido
	var current Order
	row := db.QueryRow("SELECT id, customer_name, product, quantity, total, status FROM orders WHERE id = ?", id)
	if err := row.Scan(&current.ID, &current.CustomerName, &current.Product, &current.Quantity, &current.Total, &current.Status); err == sql.ErrNoRows {
		jsonResponse(w, http.StatusNotFound, map[string]string{"error": "Order not found"})
		return
	}

	newStatus := current.Status
	if input.Status != nil && *input.Status != "" {
		newStatus = *input.Status
	}
	newCustomer := current.CustomerName
	if input.CustomerName != "" {
		newCustomer = input.CustomerName
	}
	newProduct := current.Product
	if input.Product != "" {
		newProduct = input.Product
	}
	newQuantity := current.Quantity
	if input.Quantity > 0 {
		newQuantity = input.Quantity
	}
	newTotal := current.Total
	if input.Total > 0 {
		newTotal = input.Total
	}

	_, err := db.Exec(
		"UPDATE orders SET customer_name = ?, product = ?, quantity = ?, total = ?, status = ? WHERE id = ?",
		newCustomer, newProduct, newQuantity, newTotal, newStatus, id,
	)
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	jsonResponse(w, http.StatusOK, map[string]interface{}{
		"message": fmt.Sprintf("Order %d updated successfully", id),
		"order": map[string]interface{}{
			"id":            id,
			"customer_name": newCustomer,
			"product":       newProduct,
			"quantity":      newQuantity,
			"total":         newTotal,
			"status":        newStatus,
		},
	})
}

// 5. Eliminar pedidos
func deleteOrder(w http.ResponseWriter, r *http.Request, id int) {
	res, err := db.Exec("DELETE FROM orders WHERE id = ?", id)
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	count, _ := res.RowsAffected()
	if count == 0 {
		jsonResponse(w, http.StatusNotFound, map[string]string{"error": "Order not found"})
		return
	}

	jsonResponse(w, http.StatusOK, map[string]string{"message": fmt.Sprintf("Order %d deleted successfully", id)})
}

func main() {
	initDB()

	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/orders", ordersHandler)
	http.HandleFunc("/orders/", ordersHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8002"
	}

	log.Printf("Orders API running on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
