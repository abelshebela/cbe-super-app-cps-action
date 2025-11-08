# container image builder docker file

FROM golang:1.25 AS builder

WORKDIR /app

ARG GITLAB_USER
ARG GITLAB_PAT

ENV GITLAB_USER=${GITLAB_USER}
ENV GITLAB_PAT=${GITLAB_PAT}

RUN git config --global url."https://${GITLAB_USER}:${GITLAB_PAT}@gitlab.com/".insteadOf "https://gitlab.com/"

ENV GOPRIVATE=gitlab.com/bersufekadgetachew/*

# COPY go.mod .
# COPY go.sum .

RUN go mod init cbe-super-app-cps-action
RUN go mod tidy
# Download dependencies
# RUN go mod download

# Copy all source files
COPY . .

# Build binary
RUN CGO_ENABLED=0 GOOS=linux go build -o supperapp-cps-action cmd/main.go

# Copy runtime config
# COPY corebanking.json .

EXPOSE 8080
EXPOSE 50051

CMD ["./supperapp-cps-action"]
