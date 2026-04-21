-- Same as 000013_create_account_blocks_schema.up.sql (use numbered migration in CI/CD).
CREATE TABLE account_blocks (
    id RAW(24) PRIMARY KEY,
    name VARCHAR2(255),
    code VARCHAR2(100) UNIQUE,
    address VARCHAR2(500),
    parent_id RAW(24),
    slug VARCHAR2(255),
    type VARCHAR2(1),
    is_enabled NUMBER(1) DEFAULT 1,
    city_id RAW(24),
    region_id RAW(24),
    district_id RAW(24),
    is_deleted NUMBER(1) DEFAULT 0,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    CONSTRAINT fk_account_block_city FOREIGN KEY (city_id) REFERENCES account_blocks (id) ON DELETE CASCADE,
    CONSTRAINT fk_account_block_region FOREIGN KEY (region_id) REFERENCES account_blocks (id) ON DELETE CASCADE,
    CONSTRAINT fk_account_block_district FOREIGN KEY (district_id) REFERENCES account_blocks (id) ON DELETE CASCADE
);

ALTER TABLE account_blocks
    ADD CONSTRAINT chk_ab_type CHECK (type IN ('R', 'D', 'C', 'B'));
