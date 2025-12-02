# container image builder docker file

FROM golang:1.25 AS builder

WORKDIR /app

ARG GITLAB_USER
ARG GITLAB_PAT

RUN git config --global url."https://${GITLAB_USER}:${GITLAB_PAT}@gitlab.com/".insteadOf "https://gitlab.com/"

ENV GOPRIVATE=gitlab.com/bersufekadgetachew/*

COPY go.mod .

COPY go.sum .

RUN go mod download

COPY . .

RUN go mod tidy

RUN go build -o supperapp-cps-action cmd/main.go

FROM scratch

COPY --from=builder /app/supperapp-cps-action  .

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

COPY --from=builder /app/.env .


# --------------------- Stage 2: Minimal Ubuntu Runtime ---------------------
FROM ubuntu:24.04

WORKDIR /app

# Install only what's needed to run
RUN apt-get update && apt-get install -y wget unzip libaio1t64 ca-certificates && rm -rf /var/lib/apt/lists/*

WORKDIR /opt/oracle
RUN wget https://download.oracle.com/otn_software/linux/instantclient/2390000/instantclient-basic-linux.x64-23.9.0.25.07.zip && \
    wget https://download.oracle.com/otn_software/linux/instantclient/2390000/instantclient-sdk-linux.x64-23.9.0.25.07.zip
    
RUN unzip instantclient-basic-linux.x64-23.9.0.25.07.zip
    
RUN unzip -o instantclient-sdk-linux.x64-23.9.0.25.07.zip 
    
RUN rm *.zip && \
    ln -s instantclient_23_9 instantclient

RUN  ln -s /usr/lib/x86_64-linux-gnu/libaio.so.1t64 /usr/lib/x86_64-linux-gnu/libaio.so.1


# Set environment for CGO and Oracle
ENV LD_LIBRARY_PATH=/opt/oracle/instantclient
ENV CGO_ENABLED=1
ENV OCI_LIB_DIR=/opt/oracle/instantclient
ENV OCI_INC_DIR=/opt/oracle/instantclient/sdk/include

# Copy compiled Go binary and env/configs
COPY --from=builder /app/supperapp-cps-action .
COPY --from=builder /app .
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Set runtime env for Oracle
ENV LD_LIBRARY_PATH=/opt/oracle/instantclient
# --------------------- Stage 2: Minimal Ubuntu Runtime ---------------------


EXPOSE 8080

EXPOSE 50051

CMD ["./supperapp-cps-action"]