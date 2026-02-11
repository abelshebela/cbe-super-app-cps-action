## Module: ecommerce_merchant

**Description:** This feature manages ecommerce merchants that process payments through the bank. It stores basic merchant details, their settlement account, and (optionally) branch-level information.

Users will use this when onboarding new ecommerce merchants, updating their details, or reviewing existing merchants.

---

### Screen / Form: Create Ecommerce Merchant

**Purpose:**  
To register a new ecommerce merchant with the bank, including their name, code, settlement account, and related branches.

**Input Fields**

**Merchant name**

- **Type:** Text
- **Required:** Yes
- **Description:** The full business name of the merchant as it should appear in internal systems and reports.
- **Rules:**  
  - Use letters, numbers, and spaces.  
  - Avoid unnecessary special characters.

**Merchant code**

- **Type:** Text
- **Required:** Yes
- **Description:** A short unique code used to identify the merchant (often agreed with the merchant or internal teams).
- **Rules:**  
  - Must be unique across merchants.  
  - Avoid spaces; use a consistent format such as `ABC_SHOP`.

**Settlement account number**

- **Type:** Text (numeric)
- **Required:** Yes
- **Description:** The bank account where the merchant’s settlement funds will be credited.
- **Rules:**  
  - Exactly 13 digits.  
  - Digits only; no spaces or dashes.  
  - Must be a valid account number in the bank’s core system.

**Settlement method**

- **Type:** Dropdown
- **Required:** Yes
- **Description:** How settlements will be handled for this merchant.
- **Rules:**  
  - Values are defined by your operations team (for example `Direct to account`, `GL based`, or multi-account setup).  
  - Choose the option that matches the merchant’s settlement agreement.

**Branches (repeated section – optional)**

Each branch represents one physical or logical location of the merchant. You may add zero or more branches.

For each branch:

**Branch code**

- **Type:** Text
- **Required:** Yes (if a branch is added)
- **Description:** An internal code or identifier for the branch.

**Branch name**

- **Type:** Text
- **Required:** Yes (if a branch is added)
- **Description:** Display name of the branch, such as the shop name or location.

**Branch address**

- **Type:** Text
- **Required:** Yes (if a branch is added)
- **Description:** Physical address or descriptive location of the branch.

**Branch owner / contact person**

- **Type:** Text
- **Required:** Yes (if a branch is added)
- **Description:** Name of the person responsible for the branch (e.g. shop owner or manager).

**Example (Correctly Filled Form)**

Merchant name: `Selam Online Market`  
Merchant code: `SELAM_ONLINE`  
Settlement account number: `1001234567890`  
Settlement method: `Direct to account`

Branches:

1.  
   - Branch code: `SELAM_MAIN`  
   - Branch name: `Selam Online Main Store`  
   - Branch address: `Bole, Addis Ababa – behind Friendship Mall`  
   - Branch owner / contact person: `Hanna Tesfaye`

2.  
   - Branch code: `SELAM_WEST`  
   - Branch name: `Selam Online West Warehouse`  
   - Branch address: `CMC, Addis Ababa – next to CMC Road`  
   - Branch owner / contact person: `Abel Mekonnen`

**Common Mistakes**

- Entering an account number with spaces or dashes (`1001-2345-6789-0`) instead of `1001234567890`.
- Using a merchant code that is already in use for another merchant.
- Adding a branch but leaving required branch fields empty.

---

### Screen / Form: Edit Ecommerce Merchant

**Purpose:**  
To change existing merchant details, such as their name, code, settlement account, or branch information.

**Input Fields**

The fields are the same as **Create Ecommerce Merchant**, but usually shown with existing values filled in:

- Merchant name (optional to change)  
- Merchant code (optional to change)  
- Settlement account number (optional to change)  
- Settlement method (optional to change)  
- Branch list (you can add, remove, or edit branches)

**Example (Correctly Filled Form)**

Merchant name: `Selam Online Market`  
Merchant code: `SELAM_ONLINE`  
Settlement account number: `1001234567890`  
Settlement method: `Direct to account`

Updated branch:

- Branch code: `SELAM_MAIN`  
- Branch name: `Selam Online Flagship Store`  
- Branch address: `New building on Bole Road, Addis Ababa`  
- Branch owner / contact person: `Hanna Tesfaye`

**Common Mistakes**

- Changing the merchant code without informing downstream teams (integration, reporting, etc.).  
- Entering a new settlement account number that does not exist or is not active.

---

### Screen / Form: Search and List Ecommerce Merchants

**Purpose:**  
To view and filter the list of merchants for review, QA, or reporting.

**Input Fields**

**Page**

- **Type:** Number  
- **Required:** No

**Items per page**

- **Type:** Number  
- **Required:** No

**Search**

- **Type:** Text  
- **Required:** No  
- **Description:** Simple keyword search (usually by merchant name or code).

**Filter (advanced)**

- **Type:** Text (advanced)  
- **Required:** No  
- **Description:** Additional filters such as phone number or status, based on internal conventions.

**Example (Correctly Filled Form)**

Page: `1`  
Items per page: `20`  
Search: `Selam`  
Filter: *(left empty)*

**Common Mistakes**

- Forgetting that filters are still applied and thinking merchants “disappeared” – always review filters when results look unexpected.

---

### Screen / Form: Enable / Disable Ecommerce Merchant

**Purpose:**  
To mark a merchant as active or inactive for transactions.

**Input Fields**

Typically this is managed with a simple toggle or button on the merchant details screen:

**Active / Enabled**

- **Type:** Checkbox / Toggle  
- **Required:** Yes  
- **Description:** When turned off, the merchant should no longer be used for new ecommerce transactions.

**Example (Correctly Filled Form)**

Active / Enabled: `Off` for merchant `OldTestShop` (no longer in use)  
Active / Enabled: `On` for merchant `Selam Online Market`

**Common Mistakes**

- Disabling a merchant that is still in use without confirming with business owners.  
- Confusing “disable” with “delete” – disabling should be reversible, and existing records remain in the system.

