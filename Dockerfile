# Combine front-end and back-end into a single container
FROM node:22-alpine AS frontend_builder
WORKDIR /frontend

COPY frontend/package*.json ./
RUN npm install
COPY frontend/. .
RUN npm run build

FROM golang:1.24-alpine AS build
RUN apk add --no-cache alpine-sdk

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
COPY --from=frontend_builder /frontend/dist /app/frontend/dist

RUN CGO_ENABLED=1 GOOS=linux go build -o main cmd/api/main.go

FROM alpine:3.22.0 AS prod
WORKDIR /app
COPY --from=build /app/main /app/main
COPY --from=build /app/frontend/dist /app/frontend/dist
EXPOSE ${PORT}
CMD ["./main"]
