
CREATE TABLE services (
    id                  RAW(16) DEFAULT SYS_GUID() PRIMARY KEY,
    service_name        VARCHAR2(255)    NOT NULL,
    service_code        VARCHAR2(100)    NOT NULL,
    service_key         VARCHAR2(150)    NOT NULL,
    minimum_fraud_amount VARCHAR2(50),
    product_gl_account  VARCHAR2(100),
    product_gl_account_currency VARCHAR2(10),
    enabled             NUMBER(1)        DEFAULT 1,
    is_deleted          NUMBER(1)        DEFAULT 0,
    created_at          TIMESTAMP(6) WITH TIME ZONE DEFAULT SYSTIMESTAMP,
    last_modified_at    TIMESTAMP(6) WITH TIME ZONE DEFAULT SYSTIMESTAMP NOT NULL,
    deleted_at          TIMESTAMP(6) WITH TIME ZONE DEFAULT SYSTIMESTAMP NOT NULL
);

CREATE TABLE service_cap (
    id                      RAW(16) DEFAULT SYS_GUID() PRIMARY KEY,
    service_id              RAW(16) NOT NULL,
    source                  VARCHAR2(50),      -- POS, USSD, ATM, ECOMMERCE
    currency                VARCHAR2(10) NOT NULL,
    single_cap              VARCHAR2(50),
    minimum_transfer_cap    VARCHAR2(50),

    CONSTRAINT fk_service_cap_service
        FOREIGN KEY (service_id)
        REFERENCES services(id)
        ON DELETE CASCADE
);

CREATE TABLE service_keys (
    id               RAW(16) DEFAULT SYS_GUID() PRIMARY KEY,
    service_name     VARCHAR2(255) NOT NULL,
    service_key      VARCHAR2(150) NOT NULL,
    is_enabled       NUMBER(1) DEFAULT 1,
    created_at       TIMESTAMP(6) WITH TIME ZONE DEFAULT SYSTIMESTAMP,
    last_modified_at TIMESTAMP(6) WITH TIME ZONE DEFAULT SYSTIMESTAMP NOT NULL
);
