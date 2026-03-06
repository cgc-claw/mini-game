FROM golang:1.26-alpine AS builder

RUN apk add --no-cache git libc-dev linux-headers gcc musl-dev libx11-dev libxcursor-dev libxrandr-dev libxi-dev mesa-dev libxinerama-dev

WORKDIR /app

COPY main.go ./
COPY sprites/ ./sprites/
COPY player/ ./player/
COPY enemies/ ./enemies/
COPY level/ ./level/
COPY camera/ ./camera/
COPY physics/ ./physics/

RUN go mod init game && \
    go get github.com/hajimehoshi/ebiten/v2@latest && \
    go get golang.org/x/image/font/basicfont@latest && \
    go mod tidy

RUN CGO_ENABLED=1 go build -o game .

FROM alpine:latest

RUN apk add --no-cache libstdc++ libc-dev mesa libx11 libxcursor libxrandr libxi

WORKDIR /app

COPY --from=builder /app/game .

CMD ["./game"]
