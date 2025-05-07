FROM docker.io/golang:1.23.3-alpine AS build

WORKDIR /app

RUN apk add --no-cache \
    git \
    build-base \
    gpgme-dev \
    libassuan-dev \
    libgpg-error-dev \
    btrfs-progs-dev \
    linux-headers

COPY go.mod go.sum ./
RUN go mod download

ENV CGO_ENABLED=1

COPY . .
RUN go build -o panel ./cmd/api

FROM docker.io/library/alpine:3.22.0

RUN apk add --no-cache \
    gpgme \
    libassuan \
    libgpg-error \
    ca-certificates

ARG USERNAME=panel
ARG UID=42873
ARG GID=42873

RUN addgroup -g $GID $USERNAME && adduser -D -u $UID -G $USERNAME $USERNAME

WORKDIR /app
RUN chown $USERNAME:$USERNAME /app
USER $USERNAME

COPY --chown=$USERNAME:$USERNAME --from=build /app/panel .

CMD ["./panel"]

EXPOSE 8000
