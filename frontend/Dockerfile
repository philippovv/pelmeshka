FROM node:16.20.0-alpine3.18 as builder
WORKDIR /app
COPY . .
RUN npm install &&  npm run build

FROM nginx:stable-alpine3.17-slim
WORKDIR /var/www-data/html
COPY default.conf /etc/nginx/conf.d/default.conf
COPY --from=builder /app/dist/ /var/www-data/html/
EXPOSE 80 
