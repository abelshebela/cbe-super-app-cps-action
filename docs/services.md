## Module: services

**Description:** This feature defines the different CPS services (for example transfer types or payment services) that can be used in the system. It stores a name, code, and optional limits on how much can be transferred per transaction.

Users will mainly use it when adding a new type of service, adjusting limits, or temporarily disabling a service.

---

### Screen / Form: Create Service

**Purpose:**  
To register a new service that can be used within the CPS platform.

**Input Fields**

**Service name**

- **Type:** Text  
- **Required:** Yes  
- **Description:** A clear and readable name of the service as it should appear to internal users (for example “Wallet to Wallet Transfer”).
- **Rules:**  
  - Use letters, numbers, and spaces.  
  - Avoid unnecessary symbols.

**Service key**

- **Type:** Text  
- **Required:** Yes  
- **Description:** A short internal key used by the system to refer to this service.
- **Rules:**  
  - Should be unique.  
  - Typically uses uppercase letters and underscores (for example `WALLET_TO_WALLET`).

**Service code**

- **Type:** Text  
- **Required:** Yes  
- **Description:** A code used for integration, configuration, or reporting.
- **Rules:**  
  - Should follow your internal coding standards.  
  - Must be unique among services.

**CBE GL product account**

- **Type:** Text  
- **Required:** No  
- **Description:** The internal general ledger (GL) product account used for accounting entries related to this service.

**Cap (optional section – transaction limits)**

This section may be collapsed or optional. If used, all fields inside should be completed:

**Single transaction limit**

- **Type:** Number / Decimal  
- **Required:** Yes (if the cap section is used)  
- **Description:** The maximum amount allowed for a single transaction using this service.
- **Rules:**  
  - Must be zero or higher.  
  - Should reflect real business rules for this service.

**Minimum transfer amount**

- **Type:** Number / Decimal  
- **Required:** Yes (if the cap section is used)  
- **Description:** The smallest allowed transaction amount for this service.
- **Rules:**  
  - Must be zero or higher.  
  - Must not be greater than the **Single transaction limit**.

**Enabled**

- **Type:** Checkbox / Toggle  
- **Required:** No  
- **Description:** If turned on, the service is active and can be used.

**Marked as deleted**

- **Type:** Checkbox / Toggle  
- **Required:** No  
- **Description:** For internal use when the service should be treated as removed but still kept for history.

**Example (Correctly Filled Form)**

Service name: `Wallet to Wallet Transfer`  
Service key: `WALLET_TO_WALLET_TRANSFER`  
Service code: `W2W_TRF`  
CBE GL product account: `GL-TRF-001`  
Single transaction limit: `50,000.00`  
Minimum transfer amount: `10.00`  
Enabled: `On`  
Marked as deleted: `Off`

**Common Mistakes**

- Setting the **Minimum transfer amount** higher than the **Single transaction limit**.  
- Reusing an existing **Service key** or **Service code**.  
- Leaving limits blank when the business requires them to be enforced.

---

### Screen / Form: Edit Service

**Purpose:**  
To update the details or limits of an existing service.

**Input Fields**

The same as **Create Service**, but all fields are optional to change:

- Service name (optional)  
- Service key (optional)  
- Service code (optional)  
- CBE GL product account (optional)  
- Single transaction limit (optional, if caps are used)  
- Minimum transfer amount (optional, if caps are used)  
- Enabled (optional)  
- Marked as deleted (optional)

**Example (Correctly Filled Form)**

Service name: `Wallet to Wallet Transfer`  
Single transaction limit: `100,000.00`  
Minimum transfer amount: `10.00`

**Common Mistakes**

- Increasing limits without confirming with risk or compliance teams.  
- Changing the **Service key** or **Service code** without checking other systems that depend on them.

---

### Screen / Form: Search and List Services

**Purpose:**  
To view all services, filter them, and open them for configuration or review.

**Input Fields**

**Page**

- **Type:** Number  
- **Required:** No

**Items per page**

- **Type:** Number  
- **Required:** No

**Service name**

- **Type:** Text  
- **Required:** No  
- **Description:** Filter by service name.

**Service code**

- **Type:** Text  
- **Required:** No  
- **Description:** Filter by service code.

**Service type**

- **Type:** Text  
- **Required:** No  
- **Description:** Filter by service type label (if used in your setup).

**Enabled**

- **Type:** Checkbox / Toggle  
- **Required:** No  
- **Description:** Show only active/inactive services.

**Search**

- **Type:** Text  
- **Required:** No  
- **Description:** Simple search across name, code, and type.

**Filter (advanced)**

- **Type:** Text (advanced)  
- **Required:** No  
- **Description:** Additional structured filters used by advanced users.

**Example (Correctly Filled Form)**

Page: `1`  
Items per page: `20`  
Service name: `Wallet`  
Enabled: `On`  
Search: *(left empty)*

**Common Mistakes**

- Forgetting that a filter is still active and assuming a service has been deleted.  
- Using very generic search terms like `Transfer`, which may return too many results.

---

### Screen / Form: Enable / Disable Service

**Purpose:**  
To quickly activate or deactivate a service.

**Input Fields**

**Enabled**

- **Type:** Checkbox / Toggle  
- **Required:** Yes  
- **Description:** When turned off, the service should no longer be usable for new operations.

**Example (Correctly Filled Form)**

Enabled: `Off` for a legacy service that should not be used anymore  
Enabled: `On` for new services being rolled out

**Common Mistakes**

- Disabling a widely used production service without planning communication and monitoring.  
- Assuming that disabling fully removes the service – it usually stays in lists and history.

