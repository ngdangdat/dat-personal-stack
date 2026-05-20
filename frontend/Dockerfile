# Step 1: Build the Vue application
FROM node:20.18.0-alpine3.20 AS builder

WORKDIR /app

# Copy dependency configs
COPY package*.json ./
RUN npm ci

# Copy source code
COPY index.html vite.config.js ./
COPY src/ ./src/

# Compile build distribution
RUN npm run build

# Step 2: Serve using Nginx
FROM nginx:1.27.2-alpine3.20

# Copy custom configuration
COPY nginx.conf /etc/nginx/conf.d/default.conf

# Copy build artifacts to Nginx docroot
COPY --from=builder /app/dist /usr/share/nginx/html

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
