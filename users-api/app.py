import os
import sqlite3
from flask import Flask, request, jsonify

app = Flask(__name__)

DB_PATH = os.getenv("DB_PATH", "/app/data/users.db")

def init_db():
    os.makedirs(os.path.dirname(DB_PATH), exist_ok=True)
    with sqlite3.connect(DB_PATH) as conn:
        cursor = conn.cursor()
        cursor.execute("""
            CREATE TABLE IF NOT EXISTS users (
                id INTEGER PRIMARY KEY AUTOINCREMENT,
                name TEXT NOT NULL,
                email TEXT NOT NULL UNIQUE,
                password TEXT NOT NULL,
                created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
            )
        """)
        conn.commit()

def get_db():
    conn = sqlite3.connect(DB_PATH)
    conn.row_factory = sqlite3.Row
    return conn

@app.route("/health", methods=["GET"])
def health():
    return jsonify({"status": "UP", "service": "users-api"}), 200

# 1. Registrar usuario
@app.route("/users", methods=["POST"])
def create_user():
    data = request.get_json(silent=True) or {}
    name = data.get("name")
    email = data.get("email")
    password = data.get("password")

    if not name or not email or not password:
        return jsonify({"error": "name, email and password are required"}), 400

    try:
        with get_db() as conn:
            cursor = conn.cursor()
            cursor.execute(
                "INSERT INTO users (name, email, password) VALUES (?, ?, ?)",
                (name, email, password)
            )
            user_id = cursor.lastrowid
            conn.commit()
            return jsonify({
                "message": "User registered successfully",
                "user": {"id": user_id, "name": name, "email": email}
            }), 201
    except sqlite3.IntegrityError:
        return jsonify({"error": "User with this email already exists"}), 409
    except Exception as e:
        return jsonify({"error": str(e)}), 500

# 2. Consultar lista de usuarios
@app.route("/users", methods=["GET"])
def list_users():
    with get_db() as conn:
        cursor = conn.cursor()
        cursor.execute("SELECT id, name, email, created_at FROM users")
        rows = cursor.fetchall()
        users = [dict(row) for row in rows]
        return jsonify(users), 200

# 3. Consultar usuario por ID
@app.route("/users/<int:user_id>", methods=["GET"])
def get_user(user_id):
    with get_db() as conn:
        cursor = conn.cursor()
        cursor.execute("SELECT id, name, email, created_at FROM users WHERE id = ?", (user_id,))
        row = cursor.fetchone()
        if not row:
            return jsonify({"error": "User not found"}), 404
        return jsonify(dict(row)), 200

# 4. Actualizar usuario
@app.route("/users/<int:user_id>", methods=["PUT"])
def update_user(user_id):
    data = request.get_json(silent=True) or {}
    name = data.get("name")
    email = data.get("email")
    password = data.get("password")

    if not name or not email or not password:
        return jsonify({"error": "name, email and password are required"}), 400

    try:
        with get_db() as conn:
            cursor = conn.cursor()
            cursor.execute(
                "UPDATE users SET name = ?, email = ?, password = ? WHERE id = ?",
                (name, email, password, user_id)
            )
            if cursor.rowcount == 0:
                return jsonify({"error": "User not found"}), 404
            conn.commit()
            return jsonify({
                "message": "User updated successfully",
                "user": {"id": user_id, "name": name, "email": email}
            }), 200
    except sqlite3.IntegrityError:
        return jsonify({"error": "Email already in use by another user"}), 409
    except Exception as e:
        return jsonify({"error": str(e)}), 500

# 5. Eliminar usuario
@app.route("/users/<int:user_id>", methods=["DELETE"])
def delete_user(user_id):
    with get_db() as conn:
        cursor = conn.cursor()
        cursor.execute("DELETE FROM users WHERE id = ?", (user_id,))
        if cursor.rowcount == 0:
            return jsonify({"error": "User not found"}), 404
        conn.commit()
        return jsonify({"message": f"User {user_id} deleted successfully"}), 200

if __name__ == "__main__":
    init_db()
    app.run(host="0.0.0.0", port=8000)
