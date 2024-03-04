FROM golang:1.22.0-alpine AS builder

WORKDIR /app

COPY . .

RUN go mod download
RUN go build -o moodfy main.go

FROM scratch
WORKDIR /app
COPY --from=builder /app/moodfy .

CMD ["/app/moodfy"]
