package com.cloud.auth.service;

import jakarta.annotation.PostConstruct;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.stereotype.Service;

import java.io.File;
import java.sql.*;
import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;
import java.util.Map;

@Service
public class AuthService {

    @Value("${db.path}")
    private String dbPath;

    private Connection getConnection() throws SQLException {
        return DriverManager.getConnection("jdbc:sqlite:" + dbPath);
    }

    @PostConstruct
    public void initDatabase() {
        try {
            File file = new File(dbPath);
            if (file.getParentFile() != null) {
                file.getParentFile().mkdirs();
            }

            try (Connection conn = getConnection(); Statement stmt = conn.createStatement()) {
                // Tabla de usuarios registrados para autenticacion
                stmt.execute("""
                    CREATE TABLE IF NOT EXISTS users (
                        id INTEGER PRIMARY KEY AUTOINCREMENT,
                        username TEXT NOT NULL UNIQUE,
                        password TEXT NOT NULL
                    )
                """);

                // Tabla de registro de sesiones o logs de autenticacion
                stmt.execute("""
                    CREATE TABLE IF NOT EXISTS auth_logs (
                        id INTEGER PRIMARY KEY AUTOINCREMENT,
                        username TEXT NOT NULL,
                        success INTEGER NOT NULL,
                        message TEXT,
                        timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP
                    )
                """);

                // Usuario inicial por defecto si la tabla esta vacia
                ResultSet rs = stmt.executeQuery("SELECT COUNT(*) AS total FROM users");
                if (rs.next() && rs.getInt("total") == 0) {
                    stmt.execute("INSERT INTO users (username, password) VALUES ('admin', 'admin123')");
                    stmt.execute("INSERT INTO users (username, password) VALUES ('user1', 'password123')");
                }
            }
        } catch (SQLException e) {
            throw new RuntimeException("Error initializing SQLite database: " + e.getMessage(), e);
        }
    }

    public boolean register(String username, String password) {
        String sql = "INSERT INTO users (username, password) VALUES (?, ?)";
        try (Connection conn = getConnection(); PreparedStatement pstmt = conn.prepareStatement(sql)) {
            pstmt.setString(1, username);
            pstmt.setString(2, password);
            pstmt.executeUpdate();
            return true;
        } catch (SQLException e) {
            return false;
        }
    }

    public boolean validateLogin(String username, String password) {
        String sql = "SELECT id FROM users WHERE username = ? AND password = ?";
        boolean isValid = false;

        try (Connection conn = getConnection(); PreparedStatement pstmt = conn.prepareStatement(sql)) {
            pstmt.setString(1, username);
            pstmt.setString(2, password);
            try (ResultSet rs = pstmt.executeQuery()) {
                isValid = rs.next();
            }

            // Registrar log de intento de autenticacion
            String logSql = "INSERT INTO auth_logs (username, success, message) VALUES (?, ?, ?)";
            try (PreparedStatement logPstmt = conn.prepareStatement(logSql)) {
                logPstmt.setString(1, username != null ? username : "unknown");
                logPstmt.setInt(2, isValid ? 1 : 0);
                logPstmt.setString(3, isValid ? "Credenciales validas" : "Credenciales invalidas");
                logPstmt.executeUpdate();
            }
        } catch (SQLException e) {
            e.printStackTrace();
        }

        return isValid;
    }

    public List<Map<String, Object>> getAuthLogs() {
        List<Map<String, Object>> logs = new ArrayList<>();
        String sql = "SELECT id, username, success, message, timestamp FROM auth_logs ORDER BY id DESC LIMIT 20";
        try (Connection conn = getConnection();
             Statement stmt = conn.createStatement();
             ResultSet rs = stmt.executeQuery(sql)) {

            while (rs.next()) {
                Map<String, Object> log = new HashMap<>();
                log.put("id", rs.getInt("id"));
                log.put("username", rs.getString("username"));
                log.put("success", rs.getInt("success") == 1);
                log.put("message", rs.getString("message"));
                log.put("timestamp", rs.getString("timestamp"));
                logs.add(log);
            }
        } catch (SQLException e) {
            e.printStackTrace();
        }
        return logs;
    }
}
