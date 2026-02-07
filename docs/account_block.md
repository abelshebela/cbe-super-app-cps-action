## Module: account_block

**Description:** This feature lets you temporarily block or unblock whole areas of the bank (regions, districts, cities, branches). It is used when certain locations should not be allowed to use specific services, for example during maintenance, security checks, or special campaigns.

Users will mainly use this screen when they need to quickly enable or disable groups of branches or locations, without changing individual customer accounts.

---

### Screen / Form: Search and Filter Locations (Regions / Districts / Cities / Branches)

**Purpose:**  
To find and review locations (regions, districts, cities, branches) using search and filters, and prepare them for bulk enable/disable actions.

**Input Fields**

**Page**

- **Type:** Number
- **Required:** No
- **Description:** The page number you want to see in the list.
- **Rules:** Must be a whole number greater than or equal to 1.

**Items per page**

- **Type:** Number
- **Required:** No
- **Description:** How many records you want to see on one page (for example 10, 20, 50).
- **Rules:** Must be a whole number between 1 and 100.

**Search**

- **Type:** Text
- **Required:** No
- **Description:** A simple search text to quickly find a location by name or code.
- **Rules:** Avoid special characters such as `! @ # $ %`. Use letters, numbers, and spaces only.

**Filter**

- **Type:** Text (advanced filter)
- **Required:** No
- **Description:** An advanced filter used by power users. It allows filtering by things like whether the location is enabled, or which region/district it belongs to.
- **Rules:**  
  - Should follow the format agreed within your team (usually a structured text).  
  - Avoid special characters unless specifically instructed in your internal procedures.

**Example (Correctly Filled Form)**

Page: `1`  
Items per page: `20`  
Search: `Addis`  
Filter: `enabled=true,region=Addis Ababa`

**Common Mistakes**

- Entering `0` or a negative number for **Page** or **Items per page** – this will be rejected; use positive whole numbers only.
- Using many special characters in **Search** or **Filter** (for example `@@@Addis!!!`) – search may fail or return no results.
- Leaving both **Search** and **Filter** empty and expecting very specific results – the system will simply show the default list.

---

### Screen / Form: Enable or Disable Branches

**Purpose:**  
To enable or disable one or more branches at the same time and record a clear reason for the action.

**Input Fields**

**Selected branches**

- **Type:** Multi-select list (you select from the list of branches)
- **Required:** Yes
- **Description:** One or more branches you want to enable or disable.
- **Rules:**  
  - You must select at least one branch.  
  - Only valid branches from the list can be chosen.

**Reason**

- **Type:** Text
- **Required:** Yes
- **Description:** A short explanation for why you are enabling or disabling the selected branches.
- **Rules:**  
  - Should be meaningful and business-oriented (for example “System maintenance”, “Fraud investigation”).  
  - Avoid leaving it blank or using unclear comments like “N/A” or “test”.

**Example (Correctly Filled Form)**

Selected branches:  
- `Addis Ababa Main Branch`  
- `Bole Branch`

Reason: `Planned maintenance on the regional switch from 18:00 to 20:00`

**Common Mistakes**

- Not selecting any branch and trying to submit – the system will refuse the action because at least one branch is required.
- Writing a vague reason like `test` or `ok` – this will cause confusion later during audits and may be rejected by internal processes.

---

### Screen / Form: Enable or Disable Regions

**Purpose:**  
To enable or disable all branches within one or more regions at once.

**Input Fields**

**Selected regions**

- **Type:** Multi-select list
- **Required:** Yes
- **Description:** One or more regions you want to enable or disable.
- **Rules:** At least one region must be selected.

**Reason**

- **Type:** Text
- **Required:** Yes
- **Description:** Why this region is being enabled or disabled.
- **Rules:** Same as for branches – should be clear, short, and understandable later.

**Example (Correctly Filled Form)**

Selected regions:  
- `Addis Ababa Region`

Reason: `Security review completed; restoring normal operations`

**Common Mistakes**

- Selecting a region by mistake and not reviewing the list before submitting – always double-check the list of regions because this affects many branches at once.

---

### Screen / Form: Enable or Disable Districts

**Purpose:**  
To enable or disable all branches in one or more districts.

**Input Fields**

**Selected districts**

- **Type:** Multi-select list
- **Required:** Yes
- **Description:** One or more districts whose branches you want to enable or disable.
- **Rules:** At least one district is required.

**Reason**

- **Type:** Text
- **Required:** Yes
- **Description:** Explanation of the business reason.

**Example (Correctly Filled Form)**

Selected districts:  
- `Arada District`

Reason: `Temporary shutdown due to network outage in the area`

**Common Mistakes**

- Using the district action instead of branch-level action when only one branch should be affected. Make sure you are on the correct level before confirming.

---

### Screen / Form: Enable or Disable Cities

**Purpose:**  
To enable or disable all branches in one or more cities.

**Input Fields**

**Selected cities**

- **Type:** Multi-select list
- **Required:** Yes
- **Description:** One or more cities you want to enable or disable.
- **Rules:** At least one city is required.

**Reason**

- **Type:** Text
- **Required:** Yes
- **Description:** Explanation of why all branches in these cities are being affected.

**Example (Correctly Filled Form)**

Selected cities:  
- `Dire Dawa`

Reason: `City-wide public holiday; disabling services for the day`

**Common Mistakes**

- Applying a city-wide change for a minor issue that only affects one or two branches. Always check the impact before submitting.
