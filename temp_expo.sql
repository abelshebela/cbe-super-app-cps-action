--------------------------------------------------------
--  File created - Friday-April-24-2026
--------------------------------------------------------
--------------------------------------------------------
--  DDL for Table ACCESS_ITEMS_RELATION
--------------------------------------------------------

CREATE TABLE "TEMP_ACCESS_ITEMS_RELATION" (
    "ID" RAW(16),
    "PARENT_KEY" VARCHAR2(128 BYTE),
    "CHILD_KEY" VARCHAR2(128 BYTE)
) SEGMENT CREATION IMMEDIATE PCTFREE 10 PCTUSED 40 INITRANS 1 MAXTRANS 255 NOCOMPRESS LOGGING STORAGE (
    INITIAL 65536 NEXT 1048576 MINEXTENTS 1 MAXEXTENTS 2147483645 PCTINCREASE 0 FREELISTS 1 FREELIST GROUPS 1 BUFFER_POOL DEFAULT FLASH_CACHE DEFAULT CELL_FLASH_CACHE DEFAULT
) TABLESPACE "USERS";

REM INSERTING into TEMP_ACCESS_ITEMS_RELATION SET DEFINE OFF;

Insert into
    TEMP_ACCESS_ITEMS_RELATION (ID, PARENT_KEY, CHILD_KEY)
values (
        '4E7BE33753C38225E0630C6F030A1EEE',
        'cbe_to_cbe',
        'cbe_to_cbe'
    );

Insert into
    TEMP_ACCESS_ITEMS_RELATION (ID, PARENT_KEY, CHILD_KEY)
values (
        '4E7BE33753C48225E0630C6F030A1EEE',
        'mini_statement',
        'mini_statement'
    );

Insert into
    TEMP_ACCESS_ITEMS_RELATION (ID, PARENT_KEY, CHILD_KEY)
values (
        '4E7BE33753C58225E0630C6F030A1EEE',
        'bill_payment',
        'bill_payment'
    );

Insert into
    TEMP_ACCESS_ITEMS_RELATION (ID, PARENT_KEY, CHILD_KEY)
values (
        '4E7BE33753C68225E0630C6F030A1EEE',
        'voucher',
        'voucher'
    );

Insert into
    TEMP_ACCESS_ITEMS_RELATION (ID, PARENT_KEY, CHILD_KEY)
values (
        '4E7BE33753C78225E0630C6F030A1EEE',
        'chat',
        'chat_cbe_to_cbe'
    );

Insert into
    TEMP_ACCESS_ITEMS_RELATION (ID, PARENT_KEY, CHILD_KEY)
values (
        '4E7BE33753C88225E0630C6F030A1EEE',
        'self_transfer',
        'self_transfer'
    );

Insert into
    TEMP_ACCESS_ITEMS_RELATION (ID, PARENT_KEY, CHILD_KEY)
values (
        '4E7BE33753C98225E0630C6F030A1EEE',
        'scheduled_payment',
        'scheduled_payment'
    );

Insert into
    TEMP_ACCESS_ITEMS_RELATION (ID, PARENT_KEY, CHILD_KEY)
values (
        '4E7BE33753CA8225E0630C6F030A1EEE',
        'chat',
        'chat_safaricom_topup'
    );

Insert into
    TEMP_ACCESS_ITEMS_RELATION (ID, PARENT_KEY, CHILD_KEY)
values (
        '4E7BE33753CB8225E0630C6F030A1EEE',
        'merchant',
        'merchant'
    );

Insert into
    TEMP_ACCESS_ITEMS_RELATION (ID, PARENT_KEY, CHILD_KEY)
values (
        '4E7BE33753CC8225E0630C6F030A1EEE',
        'utility',
        'utility'
    );

Insert into
    TEMP_ACCESS_ITEMS_RELATION (ID, PARENT_KEY, CHILD_KEY)
values (
        '4E7BE33753CD8225E0630C6F030A1EEE',
        'sitota',
        'sitota'
    );

Insert into
    TEMP_ACCESS_ITEMS_RELATION (ID, PARENT_KEY, CHILD_KEY)
values (
        '4E7BE33753CE8225E0630C6F030A1EEE',
        'qr_payment',
        'qr_payment'
    );

Insert into
    TEMP_ACCESS_ITEMS_RELATION (ID, PARENT_KEY, CHILD_KEY)
values (
        '4E7BE33753CF8225E0630C6F030A1EEE',
        'money_request',
        'money_request'
    );

Insert into
    TEMP_ACCESS_ITEMS_RELATION (ID, PARENT_KEY, CHILD_KEY)
values (
        '4E7BE33753D08225E0630C6F030A1EEE',
        'split_bill',
        'split_bill'
    );

Insert into
    TEMP_ACCESS_ITEMS_RELATION (ID, PARENT_KEY, CHILD_KEY)
values (
        '4E7BE33753D18225E0630C6F030A1EEE',
        'wallet',
        'telebirr'
    );

Insert into
    TEMP_ACCESS_ITEMS_RELATION (ID, PARENT_KEY, CHILD_KEY)
values (
        '4E7BE33753D28225E0630C6F030A1EEE',
        'ips',
        'ips'
    );

Insert into
    TEMP_ACCESS_ITEMS_RELATION (ID, PARENT_KEY, CHILD_KEY)
values (
        '4E7BE33753D38225E0630C6F030A1EEE',
        'map',
        'map'
    );

Insert into
    TEMP_ACCESS_ITEMS_RELATION (ID, PARENT_KEY, CHILD_KEY)
values (
        '4E7BE33753D48225E0630C6F030A1EEE',
        'news',
        'news'
    );

Insert into
    TEMP_ACCESS_ITEMS_RELATION (ID, PARENT_KEY, CHILD_KEY)
values (
        '4E7BE33753D58225E0630C6F030A1EEE',
        'exchange_rate',
        'exchange_rate'
    );

Insert into
    TEMP_ACCESS_ITEMS_RELATION (ID, PARENT_KEY, CHILD_KEY)
values (
        '4E7BE33753D68225E0630C6F030A1EEE',
        'payment_verification',
        'payment_verification'
    );

Insert into
    TEMP_ACCESS_ITEMS_RELATION (ID, PARENT_KEY, CHILD_KEY)
values (
        '4E7BE33753D78225E0630C6F030A1EEE',
        'wallet',
        'cbe_birr'
    );

Insert into
    TEMP_ACCESS_ITEMS_RELATION (ID, PARENT_KEY, CHILD_KEY)
values (
        '4E7BE33753D88225E0630C6F030A1EEE',
        'donation',
        'donation'
    );

Insert into
    TEMP_ACCESS_ITEMS_RELATION (ID, PARENT_KEY, CHILD_KEY)
values (
        '4E7BE33753D98225E0630C6F030A1EEE',
        'wallet',
        'mpesa'
    );

Insert into
    TEMP_ACCESS_ITEMS_RELATION (ID, PARENT_KEY, CHILD_KEY)
values (
        '4E7BE33753DA8225E0630C6F030A1EEE',
        'transaction_list',
        'transaction_list'
    );

Insert into
    TEMP_ACCESS_ITEMS_RELATION (ID, PARENT_KEY, CHILD_KEY)
values (
        '4E7BE33753DB8225E0630C6F030A1EEE',
        'wallet',
        'wallet'
    );

Insert into
    TEMP_ACCESS_ITEMS_RELATION (ID, PARENT_KEY, CHILD_KEY)
values (
        '4E7BE33753DC8225E0630C6F030A1EEE',
        'beneficiary',
        'beneficiary'
    );

Insert into
    TEMP_ACCESS_ITEMS_RELATION (ID, PARENT_KEY, CHILD_KEY)
values (
        '4E7BE33753DD8225E0630C6F030A1EEE',
        'micro_finance',
        'micro_finance'
    );

Insert into
    TEMP_ACCESS_ITEMS_RELATION (ID, PARENT_KEY, CHILD_KEY)
values (
        '4E7BE33753DE8225E0630C6F030A1EEE',
        'topup',
        'topup'
    );

Insert into
    TEMP_ACCESS_ITEMS_RELATION (ID, PARENT_KEY, CHILD_KEY)
values (
        '4E7BE33753DF8225E0630C6F030A1EEE',
        'topup',
        'ethio_telecom_topup'
    );

Insert into
    TEMP_ACCESS_ITEMS_RELATION (ID, PARENT_KEY, CHILD_KEY)
values (
        '4E7BE33753E08225E0630C6F030A1EEE',
        'topup',
        'safaricom_topup'
    );

Insert into
    TEMP_ACCESS_ITEMS_RELATION (ID, PARENT_KEY, CHILD_KEY)
values (
        '4E7BE33753E18225E0630C6F030A1EEE',
        'chat',
        'chat'
    );

Insert into
    TEMP_ACCESS_ITEMS_RELATION (ID, PARENT_KEY, CHILD_KEY)
values (
        '4E7BE33753E28225E0630C6F030A1EEE',
        'mini_apps',
        'mini_apps'
    );

Insert into
    TEMP_ACCESS_ITEMS_RELATION (ID, PARENT_KEY, CHILD_KEY)
values (
        '4E7BE33753E38225E0630C6F030A1EEE',
        'chat',
        'chat_money_request'
    );

Insert into
    TEMP_ACCESS_ITEMS_RELATION (ID, PARENT_KEY, CHILD_KEY)
values (
        '4E7BE33753E48225E0630C6F030A1EEE',
        'budget_management',
        'budget_management'
    );

Insert into
    TEMP_ACCESS_ITEMS_RELATION (ID, PARENT_KEY, CHILD_KEY)
values (
        '4E7BE33753E58225E0630C6F030A1EEE',
        'micro_lending',
        'micro_lending'
    );

Insert into
    TEMP_ACCESS_ITEMS_RELATION (ID, PARENT_KEY, CHILD_KEY)
values (
        '4E7BE33753E68225E0630C6F030A1EEE',
        'locked_amount',
        'locked_amount'
    );

Insert into
    TEMP_ACCESS_ITEMS_RELATION (ID, PARENT_KEY, CHILD_KEY)
values (
        '4E7BE33753E78225E0630C6F030A1EEE',
        'event',
        'event'
    );

Insert into
    TEMP_ACCESS_ITEMS_RELATION (ID, PARENT_KEY, CHILD_KEY)
values (
        '4E7BE33753E88225E0630C6F030A1EEE',
        'wallet',
        'e_birr'
    );

Insert into
    TEMP_ACCESS_ITEMS_RELATION (ID, PARENT_KEY, CHILD_KEY)
values (
        '4E7BE33753E98225E0630C6F030A1EEE',
        'chat',
        'chat_ethio_telecom_topup'
    );

Insert into
    TEMP_ACCESS_ITEMS_RELATION (ID, PARENT_KEY, CHILD_KEY)
values (
        '4E7BE33753EA8225E0630C6F030A1EEE',
        'quick_pay',
        'quick_pay_ips'
    );

Insert into
    TEMP_ACCESS_ITEMS_RELATION (ID, PARENT_KEY, CHILD_KEY)
values (
        '4E7BE33753EB8225E0630C6F030A1EEE',
        'quick_pay',
        'quick_pay_star_pay'
    );

Insert into
    TEMP_ACCESS_ITEMS_RELATION (ID, PARENT_KEY, CHILD_KEY)
values (
        '4E7BE33753EC8225E0630C6F030A1EEE',
        'quick_pay',
        'quick_pay_cbe'
    );
--------------------------------------------------------
--  DDL for Index PK_TEMP_ACCESS_ITEMS_RELATION
--------------------------------------------------------

CREATE UNIQUE INDEX "PK_TEMP_ACCESS_ITEMS_RELATION" ON "TEMP_ACCESS_ITEMS_RELATION" ("ID") PCTFREE 10 INITRANS 2 MAXTRANS 255 COMPUTE STATISTICS STORAGE (
    INITIAL 65536 NEXT 1048576 MINEXTENTS 1 MAXEXTENTS 2147483645 PCTINCREASE 0 FREELISTS 1 FREELIST GROUPS 1 BUFFER_POOL DEFAULT FLASH_CACHE DEFAULT CELL_FLASH_CACHE DEFAULT
) TABLESPACE "USERS";
--------------------------------------------------------
--  DDL for Index IX_TEMP_ACCESS_ITEMS_RELATION_CHILD
--------------------------------------------------------

CREATE INDEX "IX_TEMP_ACCESS_ITEMS_RELATION_CHILD" ON "TEMP_ACCESS_ITEMS_RELATION" (UPPER("CHILD_KEY")) PCTFREE 10 INITRANS 2 MAXTRANS 255 COMPUTE STATISTICS STORAGE (
    INITIAL 65536 NEXT 1048576 MINEXTENTS 1 MAXEXTENTS 2147483645 PCTINCREASE 0 FREELISTS 1 FREELIST GROUPS 1 BUFFER_POOL DEFAULT FLASH_CACHE DEFAULT CELL_FLASH_CACHE DEFAULT
) TABLESPACE "USERS";
--------------------------------------------------------
--  DDL for Index IX_TEMP_ACCESS_ITEMS_RELATION_PARENT
--------------------------------------------------------

CREATE INDEX "IX_TEMP_ACCESS_ITEMS_RELATION_PARENT" ON "TEMP_ACCESS_ITEMS_RELATION" (UPPER("PARENT_KEY")) PCTFREE 10 INITRANS 2 MAXTRANS 255 COMPUTE STATISTICS STORAGE (
    INITIAL 65536 NEXT 1048576 MINEXTENTS 1 MAXEXTENTS 2147483645 PCTINCREASE 0 FREELISTS 1 FREELIST GROUPS 1 BUFFER_POOL DEFAULT FLASH_CACHE DEFAULT CELL_FLASH_CACHE DEFAULT
) TABLESPACE "USERS";
--------------------------------------------------------
--  DDL for Index UX_TEMP_ACCESS_ITEMS_RELATION_PAIR
--------------------------------------------------------

CREATE UNIQUE INDEX "UX_TEMP_ACCESS_ITEMS_RELATION_PAIR" ON "TEMP_ACCESS_ITEMS_RELATION" (
    UPPER("PARENT_KEY"),
    UPPER("CHILD_KEY")
) PCTFREE 10 INITRANS 2 MAXTRANS 255 COMPUTE STATISTICS STORAGE (
    INITIAL 65536 NEXT 1048576 MINEXTENTS 1 MAXEXTENTS 2147483645 PCTINCREASE 0 FREELISTS 1 FREELIST GROUPS 1 BUFFER_POOL DEFAULT FLASH_CACHE DEFAULT CELL_FLASH_CACHE DEFAULT
) TABLESPACE "USERS";
--------------------------------------------------------
--  Constraints for Table TEMP_ACCESS_ITEMS_RELATION
--------------------------------------------------------

ALTER TABLE "TEMP_ACCESS_ITEMS_RELATION"
MODIFY ("ID" NOT NULL ENABLE);

ALTER TABLE "TEMP_ACCESS_ITEMS_RELATION"
MODIFY ("PARENT_KEY" NOT NULL ENABLE);

ALTER TABLE "TEMP_ACCESS_ITEMS_RELATION"
MODIFY ("CHILD_KEY" NOT NULL ENABLE);

ALTER TABLE "TEMP_ACCESS_ITEMS_RELATION"
ADD CONSTRAINT "PK_TEMP_ACCESS_ITEMS_RELATION" PRIMARY KEY ("ID") USING INDEX "PK_TEMP_ACCESS_ITEMS_RELATION" ENABLE;