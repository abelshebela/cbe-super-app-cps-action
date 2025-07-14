# container image builder docker file

FROM golang:1.24 AS builder

WORKDIR /app

COPY . .

COPY go.mod .

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -o supperapp-cps-action-bin cmd/main.go

FROM scratch

COPY --from=builder /app/supperapp-cps-action-bin .

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

COPY --from=builder /app/.env .

EXPOSE 8080

CMD ["./supperapp-cps-action-bin"]
