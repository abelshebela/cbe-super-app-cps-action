CREATE TABLE app_access_list (
    id               VARCHAR2(36) PRIMARY KEY,
    key              VARCHAR2(200) NOT NULL,
    enabled          NUMBER(1) DEFAULT 1,
    access_list_name VARCHAR2(255) NOT NULL,
    ussd_enabled     NUMBER(1) DEFAULT 0,
    created_at       TIMESTAMP(6) WITH TIME ZONE DEFAULT SYSTIMESTAMP,
    last_modified_at TIMESTAMP(6) WITH TIME ZONE
);

CREATE TABLE app_sub_access_list (
    id                VARCHAR2(36) PRIMARY KEY,
    parent_id         VARCHAR2(36) NOT NULL,
    sub_key           VARCHAR2(200),
    name              VARCHAR2(255),
    enabled           NUMBER(1) DEFAULT 1,
    
    CONSTRAINT fk_app_sub_access_list_parent
        FOREIGN KEY (parent_id)
        REFERENCES app_access_list(id)
        ON DELETE CASCADE
);