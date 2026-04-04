-- Disable history for account_blocks (one row per disable event; cleared on re-enable from app).
CREATE TABLE account_block_disable_reasons (
    id VARCHAR2(36) NOT NULL,
    account_block_id VARCHAR2(36) NOT NULL,
    reason_text VARCHAR2(512),
    created_by VARCHAR2(128),
    created_at TIMESTAMP DEFAULT SYSTIMESTAMP NOT NULL,
    CONSTRAINT pk_account_block_disable_reasons PRIMARY KEY (id)
);

CREATE INDEX idx_abdr_account_block_id ON account_block_disable_reasons (account_block_id);
