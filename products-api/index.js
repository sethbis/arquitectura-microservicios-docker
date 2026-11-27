const express = require('express');
const sqlite3 = require('sqlite3').verbose();
const path = require('path');
const fs = require('fs');

const app = express();
app.use(express.json());

const PORT = process.env.PORT || 8003;
const DB_PATH = process.env.DB_PATH || '/app/data/products.db';

// Asegurar que el directorio de la base de datos existe
const dbDir = path.dirname(DB_PATH);
if (!fs.existsSync(dbDir)) {
    fs.mkdirSync(dbDir, { recursive: true });
}

// Inicializar conexión SQLite
const db = new sqlite3.Database(DB_PATH, (err) => {
    if (err) {
        console.error('Error opening SQLite database:', err.message);
        process.exit(1);
    }
    console.log(`Connected to SQLite products database at ${DB_PATH}`);
});

// Crear tabla si no existe
db.serialize(() => {
    db.run(`
        CREATE TABLE IF NOT EXISTS products (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            name TEXT NOT NULL,
            price REAL NOT NULL,
            description TEXT,
            stock INTEGER NOT NULL DEFAULT 0,
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
        )
    `, (err) => {
        if (err) {
            console.error('Error creating products table:', err.message);
        } else {
            console.log('Products table verified/created successfully.');
        }
    });
});

// Health check endpoint
app.get(['/health', '/products/health'], (req, res) => {
    res.status(200).json({ status: 'UP', service: 'products-api' });
});

// 1. Registrar productos
app.post('/products', (req, res) => {
    const { name, price, description, stock } = req.body;

    if (!name || price === undefined || stock === undefined) {
        return res.status(400).json({ error: 'name, price and stock are required' });
    }

    const query = `INSERT INTO products (name, price, description, stock) VALUES (?, ?, ?, ?)`;
    db.run(query, [name, price, description || '', stock], function (err) {
        if (err) {
            return res.status(500).json({ error: err.message });
        }
        res.status(201).json({
            message: 'Product registered successfully',
            product: {
                id: this.lastID,
                name,
                price,
                description: description || '',
                stock
            }
        });
    });
});

// 2. Consultar la lista de productos
app.get('/products', (req, res) => {
    const query = `SELECT id, name, price, description, stock, created_at FROM products`;
    db.all(query, [], (err, rows) => {
        if (err) {
            return res.status(500).json({ error: err.message });
        }
        res.status(200).json(rows);
    });
});

// 3. Consultar un producto por su ID
app.get('/products/:id', (req, res) => {
    const id = parseInt(req.params.id, 10);
    if (isNaN(id)) {
        return res.status(400).json({ error: 'Invalid product ID' });
    }

    const query = `SELECT id, name, price, description, stock, created_at FROM products WHERE id = ?`;
    db.get(query, [id], (err, row) => {
        if (err) {
            return res.status(500).json({ error: err.message });
        }
        if (!row) {
            return res.status(404).json({ error: 'Product not found' });
        }
        res.status(200).json(row);
    });
});

// 4. Actualizar la información de un producto
app.put('/products/:id', (req, res) => {
    const id = parseInt(req.params.id, 10);
    if (isNaN(id)) {
        return res.status(400).json({ error: 'Invalid product ID' });
    }

    const { name, price, description, stock } = req.body;
    if (!name || price === undefined || stock === undefined) {
        return res.status(400).json({ error: 'name, price and stock are required' });
    }

    const query = `
        UPDATE products 
        SET name = ?, price = ?, description = ?, stock = ? 
        WHERE id = ?
    `;
    db.run(query, [name, price, description || '', stock, id], function (err) {
        if (err) {
            return res.status(500).json({ error: err.message });
        }
        if (this.changes === 0) {
            return res.status(404).json({ error: 'Product not found' });
        }
        res.status(200).json({
            message: `Product ${id} updated successfully`,
            product: { id, name, price, description: description || '', stock }
        });
    });
});

// 5. Eliminar productos
app.delete('/products/:id', (req, res) => {
    const id = parseInt(req.params.id, 10);
    if (isNaN(id)) {
        return res.status(400).json({ error: 'Invalid product ID' });
    }

    const query = `DELETE FROM products WHERE id = ?`;
    db.run(query, [id], function (err) {
        if (err) {
            return res.status(500).json({ error: err.message });
        }
        if (this.changes === 0) {
            return res.status(404).json({ error: 'Product not found' });
        }
        res.status(200).json({ message: `Product ${id} deleted successfully` });
    });
});

app.listen(PORT, '0.0.0.0', () => {
    console.log(`Products API listening on port ${PORT}`);
});
