CREATE TABLE access_list (
    id               RAW(16) DEFAULT SYS_GUID() PRIMARY KEY,
    key              VARCHAR2(200) NOT NULL,
    enabled          NUMBER(1) DEFAULT 1,
    access_list_name VARCHAR2(255) NOT NULL,
    ussd_enabled     NUMBER(1) DEFAULT 0,
    created_at       TIMESTAMP(6) WITH TIME ZONE DEFAULT SYSTIMESTAMP,
    last_modified_at TIMESTAMP(6) WITH TIME ZONE
);

CREATE TABLE sub_access_list (
    id                RAW(16) DEFAULT SYS_GUID() PRIMARY KEY,
    parent_id         RAW(16) NOT NULL,
    sub_key           VARCHAR2(200),
    name              VARCHAR2(255),
    enabled           NUMBER(1) DEFAULT 1,
    
    CONSTRAINT fk_sub_access_list_parent
        FOREIGN KEY (parent_id)
        REFERENCES access_list(id)
        ON DELETE CASCADE
);