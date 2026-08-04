# --------------------- Stage 1: Go Builder ---------------------
FROM golang:1.26 AS builder

WORKDIR /app

ARG GITLAB_USER
ARG GITLAB_PAT

RUN git config --global url."https://${GITLAB_USER}:${GITLAB_PAT}@gitlab.com/".insteadOf "https://gitlab.com/"

ENV GOPRIVATE=gitlab.com/bersufekadgetachew/*

COPY go.mod go.sum ./

RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

COPY . .

RUN --mount=type=cache,target=/go/pkg/mod \
    go mod tidy

RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg/mod \
    CGO_ENABLED=1 go build -buildvcs=false -o supperapp-cps-action cmd/main.go

# --------------------- Stage 2: Minimal Ubuntu Runtime ---------------------
FROM ubuntu:24.04

# Set environments
ENV LD_LIBRARY_PATH=/opt/oracle/instantclient
ENV CGO_ENABLED=1
ENV OCI_LIB_DIR=/opt/oracle/instantclient
ENV OCI_INC_DIR=/opt/oracle/instantclient/sdk/include

WORKDIR /opt/oracle

# Install dependencies and setup Oracle Client in a single cached layer
RUN apt-get update && apt-get install -y wget unzip libaio1t64 ca-certificates && \
    wget -q https://download.oracle.com/otn_software/linux/instantclient/2390000/instantclient-basic-linux.x64-23.9.0.25.07.zip && \
    wget -q https://download.oracle.com/otn_software/linux/instantclient/2390000/instantclient-sdk-linux.x64-23.9.0.25.07.zip && \
    unzip -q instantclient-basic-linux.x64-23.9.0.25.07.zip && \
    unzip -q -o instantclient-sdk-linux.x64-23.9.0.25.07.zip && \
    rm *.zip && \
    ln -s instantclient_23_9 instantclient && \
    ln -s /usr/lib/x86_64-linux-gnu/libaio.so.1t64 /usr/lib/x86_64-linux-gnu/libaio.so.1 && \
    rm -rf /var/lib/apt/lists/*

WORKDIR /app

# Copy compiled Go binary and certificates from builder
COPY --from=builder /app/supperapp-cps-action .
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/.env . 

EXPOSE 8080 50051

CMD ["./supperapp-cps-action"]