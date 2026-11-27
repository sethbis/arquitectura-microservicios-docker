FROM node:20-alpine

WORKDIR /app

# Instalar dependencias
COPY package*.json ./
RUN npm install --production

# Copiar código fuente
COPY . .

RUN mkdir -p /app/data

EXPOSE 8003

ENV DB_PATH=/app/data/products.db
ENV PORT=8003

CMD ["node", "index.js"]
