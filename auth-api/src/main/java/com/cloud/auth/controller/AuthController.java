package com.cloud.auth.controller;

import com.cloud.auth.service.AuthService;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

import java.util.HashMap;
import java.util.Map;

@RestController
@RequestMapping
public class AuthController {

    private final AuthService authService;

    public AuthController(AuthService authService) {
        this.authService = authService;
    }

    // Endpoint para verificar el estado de la API
    @GetMapping({"/health", "/auth/health"})
    public ResponseEntity<Map<String, Object>> health() {
        Map<String, Object> response = new HashMap<>();
        response.put("status", "UP");
        response.put("service", "auth-api");
        return ResponseEntity.ok(response);
    }

    // Endpoint para registrar credenciales
    @PostMapping("/auth/register")
    public ResponseEntity<Map<String, Object>> register(@RequestBody Map<String, String> payload) {
        String username = payload.get("username");
        String password = payload.get("password");

        Map<String, Object> response = new HashMap<>();
        if (username == null || username.isBlank() || password == null || password.isBlank()) {
            response.put("valid", false);
            response.put("error", "Username and password are required");
            return ResponseEntity.badRequest().body(response);
        }

        boolean created = authService.register(username, password);
        if (!created) {
            response.put("valid", false);
            response.put("error", "Username already exists or registration failed");
            return ResponseEntity.status(HttpStatus.CONFLICT).body(response);
        }

        response.put("valid", true);
        response.put("message", "User registered successfully");
        response.put("username", username);
        return ResponseEntity.status(HttpStatus.CREATED).body(response);
    }

    // Endpoint para login y validacion de credenciales
    @PostMapping({"/login", "/auth/login"})
    public ResponseEntity<Map<String, Object>> login(@RequestBody Map<String, String> payload) {
        String username = payload.get("username");
        String password = payload.get("password");

        Map<String, Object> response = new HashMap<>();
        if (username == null || password == null) {
            response.put("valid", false);
            response.put("message", "Username and password are required");
            return ResponseEntity.badRequest().body(response);
        }

        boolean isValid = authService.validateLogin(username, password);
        response.put("valid", isValid);

        if (isValid) {
            response.put("message", "Credenciales válidas. Inicio de sesión exitoso.");
            response.put("username", username);
            return ResponseEntity.ok(response);
        } else {
            response.put("message", "Credenciales inválidas. Usuario o contraseña incorrectos.");
            return ResponseEntity.status(HttpStatus.UNAUTHORIZED).body(response);
        }
    }

    // Consulta de registros de sesiones / logs
    @GetMapping("/auth/logs")
    public ResponseEntity<?> getLogs() {
        return ResponseEntity.ok(authService.getAuthLogs());
    }
}
