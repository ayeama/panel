FROM docker.io/library/node:24.13.0-trixie AS builder-frontend

WORKDIR /app

COPY web/package.json web/package-lock.json .

RUN npm install

COPY web/ .

RUN npm run build

FROM docker.io/library/golang:1.25.5-trixie AS builder-backend

WORKDIR /app

RUN apt-get update -y && \
    apt-get install -y \
        sqlite3 && \
    apt-get clean

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=1 go build -o panel cmd/panel/main.go

FROM docker.io/library/debian:trixie-slim

RUN apt-get update -y && \
    apt-get install -y \
        sqlite3 && \
    apt-get clean

# RUN adduser -D panel

# USER panel

WORKDIR /app

COPY --from=builder-backend /app/panel .
COPY --from=builder-frontend /app/dist ./static

EXPOSE 8000

CMD ["./panel"]
