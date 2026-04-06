-- WALLETS (Oracle) — aligned with internal/storage/persistance/wallet/oracle
CREATE TABLE wallets (
    id RAW(16) PRIMARY KEY,
    name VARCHAR2(255),
    wallet_code VARCHAR2(255),
    service_id RAW(16),
    enabled NUMBER(1) DEFAULT 0,
    avatar VARCHAR2(255),
    services_self NUMBER(1) DEFAULT 0,
    services_other NUMBER(1) DEFAULT 0,
    services_agent NUMBER(1) DEFAULT 0,
    is_deleted NUMBER(1) DEFAULT 0,
    created_at TIMESTAMP DEFAULT SYSTIMESTAMP,
    last_modified_at TIMESTAMP DEFAULT SYSTIMESTAMP,
    deleted_at TIMESTAMP,
    CONSTRAINT fk_wallets_service FOREIGN KEY (service_id) REFERENCES services (id) ON DELETE CASCADE
);
