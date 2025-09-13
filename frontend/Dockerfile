FROM node:20-alpine AS builder

WORKDIR /app/frontend

COPY package*.json ./
RUN if [ -f package-lock.json ]; then npm ci --ignore-scripts; else npm install --ignore-scripts; fi

COPY . .
RUN npm run build

FROM nginx:1.27-alpine AS runtime

WORKDIR /app/frontend

COPY --from=builder /app/frontend/dist /usr/share/nginx/html

RUN printf '%s\n' \
  'server {' \
  '  listen 80;' \
  '  server_name _;' \
  '  root /usr/share/nginx/html;' \
  '  index index.html;' \
  '  # Serve static files with long cache (except index.html)' \
  '  location ~* \.(?:js|css|png|jpg|jpeg|gif|ico|svg|woff2?)$ {' \
  '    expires 7d;' \
  '    add_header Cache-Control "public, max-age=604800, immutable";' \
  '  }' \
  '  # SPA fallback: send all non-file routes to index.html' \
  '  location / {' \
  '    try_files $uri /index.html;' \
  '  }' \
  '}' \
  > /etc/nginx/conf.d/default.conf

EXPOSE 80