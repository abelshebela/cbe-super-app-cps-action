# container image builder docker file

FROM golang:1.24 AS builder

WORKDIR /app

ARG GITLAB_USER
ARG GITLAB_PAT

ENV GITLAB_USER=${GITLAB_USER}
ENV GITLAB_PAT=${GITLAB_PAT}

RUN git config --global url."https://${GITLAB_USER}:${GITLAB_PAT}@gitlab.com/".insteadOf "https://gitlab.com/"

ENV GOPRIVATE=gitlab.com/bersufekadgetachew/*

COPY go.mod .

COPY go.sum .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o supperapp-cps-action cmd/main.go

FROM scratch

COPY --from=builder /app/supperapp-cps-action  .

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

COPY --from=builder /app/.env .

EXPOSE 50051

CMD ["./supperapp-cps-action"]
