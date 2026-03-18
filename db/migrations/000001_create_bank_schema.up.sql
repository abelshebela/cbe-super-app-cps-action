CREATE TABLE banks (
    id RAW(16) DEFAULT SYS_GUID() PRIMARY KEY,
    bank_name VARCHAR2(32) NOT NULL,
    logo VARCHAR2(255) NOT NULL,
    bic_code VARCHAR2(16) NOT NULL,
    is_enabled NUMBER(1) DEFAULT 1,
    account_length NUMBER NOT NULL,
    has_alpha_numeric NUMBER(1) DEFAULT 0,
    create_at TIMESTAMP DEFAULT SYSTIMESTAMP,
    update_at TIMESTAMP DEFAULT SYSTIMESTAMP
);