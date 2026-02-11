## Module: vault_amount_tier

**Description:** This feature defines “amount tiers” for vault products (such as savings or group vaults). Each tier sets a range of amounts and the interest rate that applies within that range.

Users will use this when configuring or adjusting the interest structure for vault products.

---

### Screen / Form: Create Vault Amount Tier

**Purpose:**  
To add a new amount range (tier) and link it to a vault product category, with a specific interest rate.

**Input Fields**

**Vault category**

- **Type:** Dropdown
- **Required:** Yes
- **Description:** The vault product (for example “Personal Vault – Standard”, “Group Vault – Gold”) to which this tier belongs.
- **Rules:**  
  - Choose from existing categories only.  
  - Make sure you select the correct category, as this controls which customers are affected.

**Minimum amount**

- **Type:** Number / Decimal
- **Required:** Yes
- **Description:** The smallest balance for which this tier should apply.
- **Rules:**  
  - Must be greater than 0.  
  - Must not be higher than the **Maximum amount**.

**Maximum amount**

- **Type:** Number / Decimal
- **Required:** Yes
- **Description:** The highest balance for which this tier should apply.
- **Rules:**  
  - Must be greater than 0.  
  - Must be greater than or equal to the **Minimum amount**.

**Interest rate**

- **Type:** Number / Decimal
- **Required:** Yes
- **Description:** The interest rate for balances in this tier.
- **Rules:**  
  - Must be zero or higher.  
  - Typically expressed as a percentage (for example `7.5` for 7.5%).

**Example (Correctly Filled Form)**

Vault category: `Personal Vault – Standard`  
Minimum amount: `1,000.00`  
Maximum amount: `10,000.00`  
Interest rate: `7.5`

**Common Mistakes**

- Setting the **Minimum amount** higher than the **Maximum amount** (for example min `10,000` and max `5,000`).  
- Overlapping ranges between tiers for the same category (for example one tier covers `1,000–10,000` and another covers `8,000–15,000`). This causes confusion on which rate should apply.

---

### Screen / Form: Edit Vault Amount Tier

**Purpose:**  
To correct or change the amount range or interest rate for an existing tier.

**Input Fields**

Same fields as **Create Vault Amount Tier** (usually with the category fixed and not changed):

- Minimum amount (required)  
- Maximum amount (required)  
- Interest rate (required)

**Example (Correctly Filled Form)**

Minimum amount: `10,000.00`  
Maximum amount: `50,000.00`  
Interest rate: `8.25`

**Common Mistakes**

- Lowering the **Maximum amount** so that there is a “gap” which is not covered by any tier.  
- Setting the **Interest rate** to a negative value or a clearly unrealistic number by mistake (for example `-1` or `1000`).

---

### Screen / Form: Search and List Vault Amount Tiers

**Purpose:**  
To view and filter all tiers for review, QA, or product configuration checks.

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
- **Description:** Simple search, typically by category name or other visible text.

**Filter (advanced)**

- **Type:** Text (advanced)  
- **Required:** No  
- **Description:** For more detailed filters (for example only enabled tiers), depending on internal usage.

**Example (Correctly Filled Form)**

Page: `1`  
Items per page: `20`  
Search: `Personal Vault`  
Filter: *(left empty)*

**Common Mistakes**

- Forgetting that filters or page size are still applied and not seeing expected tiers.  
- Using special characters or very long texts in **Search**, which may not be helpful.

---

### Screen / Form: Enable / Disable Vault Amount Tier

**Purpose:**  
To activate or deactivate a specific tier without deleting it.

**Input Fields**

Usually controlled by a simple toggle:

**Active / Enabled**

- **Type:** Checkbox / Toggle  
- **Required:** Yes  
- **Description:** When turned off, this tier is no longer used in calculations, but it stays in the system for reference.

**Example (Correctly Filled Form)**

Active / Enabled: `Off` for a legacy tier  
Active / Enabled: `On` for current tiers used by customers

**Common Mistakes**

- Disabling a tier that is still part of a published product offering without checking with product management.  
- Confusing disabling with deleting – disabling hides the tier from active use but keeps history.

