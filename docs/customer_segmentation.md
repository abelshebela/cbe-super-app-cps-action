## Module: customer_segmentation

**Description:** This feature groups customers into segments for reporting, pricing, and product rules. Each segmentation links a customer role (for example “Retail”, “Corporate”) to one or more sub‑segments and groups.

Users will use this when defining or adjusting how customers from the core banking system are mapped into business segments in the CPS platform.

---

### Screen / Form: Create Customer Segmentation

**Purpose:**  
To create a new segmentation rule that assigns one customer role to one or more sub‑segments and groups.

**Input Fields**

**Customer role**

- **Type:** Dropdown
- **Required:** Yes
- **Description:** The high-level role of the customer (for example “Retail customer”, “Corporate customer”, “Agent”).
- **Rules:**  
  - Must match an existing role in the CPS system.  
  - Usually provided as a list to select from, not typed manually.

**T24 customer sub‑segments (repeated section)**

- **Type:** List of rows / repeatable group
- **Required:** Yes (at least one row)
- **Description:** One or more mappings to sub‑segments from the core banking system.

For each sub‑segment:

**Name**

- **Type:** Text  
- **Required:** Yes  
- **Description:** A friendly name for the sub‑segment (for example “High value salaried customers”).

**Customer segment (cust_segment)**

- **Type:** Text  
- **Required:** Yes  
- **Description:** The code or description used in the core banking system for the segment.

**Customer group (cust_group)**

- **Type:** Text  
- **Required:** Yes  
- **Description:** The group code or description from the core banking system.

**Example (Correctly Filled Form)**

Customer role: `Retail customer`

T24 customer sub‑segments:

1.  
   - Name: `Salary customers – standard`  
   - Customer segment: `SAL_STD`  
   - Customer group: `RETAIL_SALARY`

2.  
   - Name: `Salary customers – premium`  
   - Customer segment: `SAL_PREM`  
   - Customer group: `RETAIL_PREMIUM`

**Common Mistakes**

- Leaving the sub‑segment list empty – at least one row is required.  
- Entering unclear or inconsistent codes for **Customer segment** or **Customer group**, which makes later mapping and reporting unreliable.  
- Using the same codes in multiple unrelated segmentations, which can cause confusion.

---

### Screen / Form: Edit Customer Segmentation

**Purpose:**  
To update the sub‑segment mappings for an existing customer role, for example when codes change in the core banking system.

**Input Fields**

Same structure as **Create Customer Segmentation**, usually with **Customer role** fixed:

- T24 customer sub‑segments (list; at least one row)  
  - Name (required)  
  - Customer segment (required)  
  - Customer group (required)

**Example (Correctly Filled Form)**

T24 customer sub‑segments:

1.  
   - Name: `Salary customers – premium`  
   - Customer segment: `SAL_PREM_NEW`  
   - Customer group: `RETAIL_PREMIUM_V2`

**Common Mistakes**

- Removing all sub‑segments and saving – this would leave the role without any valid mapping.  
- Typing placeholder values like `test` or `xxx` for codes that are expected to match real codes from the core system.

---

### Screen / Form: Search and List Customer Segmentations

**Purpose:**  
To view all segmentations, filter them, and open them for editing or review.

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
- **Description:** Simple search across role name or sub‑segment names.

**Filter (advanced)**

- **Type:** Text (advanced)  
- **Required:** No  
- **Description:** Additional structured filters according to internal usage.

**Example (Correctly Filled Form)**

Page: `1`  
Items per page: `20`  
Search: `Retail`  
Filter: *(left empty)*

**Common Mistakes**

- Forgetting to clear filters when looking for newly created segmentations.  
- Using very broad search terms that return too many results.

---

### Screen / Form: Enable / Disable Customer Segmentation

**Purpose:**  
To turn a segmentation rule on or off without deleting it, for example when a mapping is no longer valid but should remain in the history.

**Input Fields**

**Active / Enabled**

- **Type:** Checkbox / Toggle  
- **Required:** Yes  
- **Description:** When turned off, this segmentation rule is no longer used in customer classification.

**Example (Correctly Filled Form)**

Active / Enabled: `Off` for an old mapping that uses retired T24 codes  
Active / Enabled: `On` for the latest mapping that matches current T24 configuration

**Common Mistakes**

- Disabling all segmentations for a given **Customer role**, leaving that role with no active mapping.  
- Confusing disabling with deleting – disabling stops usage but keeps the record for reference.

