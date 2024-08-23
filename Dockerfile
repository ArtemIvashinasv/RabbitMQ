FROM golang:1.20-alpine AS builder

WORKDIR /app

# Dependences
COPY ["app/go.mod", "app/go.sum", "./"]
RUN go mod download

# Build
COPY app ./

RUN go build -o ./bin/app cmd/main.go

FROM alpine AS runner

WORKDIR /app

COPY --from=builder /app/bin/app /app/app
COPY --from=builder /app/.env /app/.env

CMD [ "/app/app" ]