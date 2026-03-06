FROM golang:1.26-alpine AS builder

RUN apk add --no-cache git libc-dev linux-headers gcc musl-dev libx11-dev libxcursor-dev libxrandr-dev libxi-dev mesa-dev libxinerama-dev mingw-w64-gcc

WORKDIR /app

COPY go.mod go.sum ./
COPY main.go ./
COPY sprites/ ./sprites/
COPY player/ ./player/
COPY enemies/ ./enemies/
COPY level/ ./level/
COPY camera/ ./camera/
COPY physics/ ./physics/

RUN go mod download

# Build for Linux
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o game-linux .

# Build for Windows  
RUN CGO_ENABLED=1 GOOS=windows GOARCH=amd64 CC=x86_64-w64-mingw32-gcc go build -o game-windows.exe .

# Build for macOS (amd64) - requires osxcross with macOS SDK
# This step is skipped in Docker unless osxcross is properly set up with SDK
# For full macOS support, use the separate macOS job in the workflow

FROM alpine:latest

RUN apk add --no-cache libstdc++ libc-dev mesa libx11 libxcursor libxrandr libxi

WORKDIR /app

COPY --from=builder /app/game-linux ./game-linux
COPY --from=builder /app/game-windows.exe ./game-windows.exe

CMD ["./game-linux"]
