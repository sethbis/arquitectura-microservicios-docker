# Build stage
FROM maven:3.9.6-eclipse-temurin-17-alpine AS builder

WORKDIR /build

COPY pom.xml .
# Descargar dependencias primero para cachear capas
RUN mvn dependency:go-offline -B || true

COPY src ./src
RUN mvn clean package -DskipTests -B

# Runtime stage
FROM eclipse-temurin:17-jre-alpine

WORKDIR /app

COPY --from=builder /build/target/*.jar app.jar

RUN mkdir -p /app/data

EXPOSE 8001

ENV DB_PATH=/app/data/auth.db

ENTRYPOINT ["java", "-jar", "app.jar"]
