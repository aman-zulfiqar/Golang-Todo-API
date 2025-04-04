FROM golang:1.24.2-alpine3.21 AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o main main.go

FROM php:8.2-apache

WORKDIR /var/www/html
RUN docker-php-ext-install pdo pdo_mysql
RUN a2enmod rewrite
COPY templete/ /var/www/html/
COPY --from=builder /app/main /usr/local/bin/main

EXPOSE 80 9010

CMD service apache2 start && main
