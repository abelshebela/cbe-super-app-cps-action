# container image builder docker file

FROM golang:1.24 AS builder

WORKDIR /app

ARG GITLAB_USER
ARG GITLAB_PAT

RUN git config --global url."https://${GITLAB_USER}:${GITLAB_PAT}@gitlab.com/".insteadOf "https://gitlab.com/"

ENV GOPRIVATE=gitlab.com/bersufekadgetachew/*

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o supperapp-cps-bin-auth cmd/main.go

FROM scratch

COPY --from=builder /app/supperapp-cps-bin-auth .

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

COPy . .

EXPOSE 8080

CMD ["./supperapp-cps-bin-auth"]
