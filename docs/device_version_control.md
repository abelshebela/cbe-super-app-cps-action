## Module: device_version_control

**Description:** This feature controls which versions of the mobile app are allowed to connect and whether users must update. It is used by operations, release managers, or QA when rolling out new app versions or forcing users to upgrade.

Users will use it when a new mobile app version is released, when a faulty version must be blocked, or when a critical update must be enforced.

---

### Screen / Form: Create Device Version

**Purpose:**  
To register a new mobile app version, indicate whether it is mandatory, and describe what has changed.

**Input Fields**

**Version**

- **Type:** Text
- **Required:** Yes
- **Description:** The app version number as it appears in the mobile app (for example on the app store or build information).
- **Rules:**  
  - Use digits and dots only (for example `1.0.0`, `2.5.3`).  
  - Length should be reasonable (up to about 50 characters).

**Platform**

- **Type:** Dropdown
- **Required:** Yes
- **Description:** Which type of device this version is for.
- **Rules:**  
  - Choose one of: `Android` or `iOS`.  
  - The system will internally treat them as `ANDROID` or `IOS`.

**Force update**

- **Type:** Checkbox / Toggle
- **Required:** Yes
- **Description:** If turned on, users with older versions will be forced to update before they can continue using the app.
- **Rules:**  
  - Turn **on** only when this version is safe and fully tested.  
  - Turn **off** for soft launches or testing.

**Release notes**

- **Type:** Multi-line text
- **Required:** No
- **Description:** A short summary of what changed in this version (new features, fixes, important notes).
- **Rules:**  
  - Maximum around 500 characters.  
  - Use clear, simple language that helps support and QA understand the change.

**Example (Correctly Filled Form)**

Version: `2.3.0`  
Platform: `Android`  
Force update: `On`  
Release notes: `Introduced biometric login for retail customers and fixed timeout issues on bill payments.`

**Common Mistakes**

- Typing letters in **Version** (for example `v2.3.0-beta`) – keep it to numbers and dots only.
- Forgetting to set **Force update** correctly (for example turning it on during initial testing).
- Leaving **Release notes** blank for major releases – this makes troubleshooting harder later.

---

### Screen / Form: Update Device Version

**Purpose:**  
To adjust details of an existing app version, such as changing whether update is forced or improving the release notes.

**Input Fields**

**ID (hidden or read-only)**

- **Type:** Text (usually not editable)
- **Required:** Yes (system-managed)
- **Description:** The internal identifier of the version. You generally do not type this in; it is populated by the system.

**Version**

- **Type:** Text
- **Required:** No
- **Description:** New version value, if you need to correct a typo.
- **Rules:**  
  - Same rules as when creating: digits and dots only.

**Platform**

- **Type:** Dropdown
- **Required:** No
- **Description:** Platform for this version.
- **Rules:**  
  - Use `Android` or `iOS`.  

**Force update**

- **Type:** Checkbox / Toggle
- **Required:** No
- **Description:** Change whether this version should be mandatory.

**Release notes**

- **Type:** Multi-line text
- **Required:** No
- **Description:** Update or expand the description of changes.

**Example (Correctly Filled Form)**

Version: `2.3.0`  
Platform: `Android`  
Force update: `Off`  
Release notes: `Pilot rollout only for staff devices. Full rollout planned next week.`

**Common Mistakes**

- Changing the **Platform** incorrectly (for example marking an iOS build as Android) – always confirm with the build team.
- Making large changes to **Version** instead of creating a new version entry – use update only for corrections, not for new releases.

---

### Screen / Form: Search and Filter Device Versions

**Purpose:**  
To quickly find specific versions, check which ones are enabled, and see history for each platform.

**Input Fields**

**Page**

- **Type:** Number
- **Required:** No
- **Description:** Page number in the results.

**Items per page**

- **Type:** Number
- **Required:** No

**Search**

- **Type:** Text
- **Required:** No
- **Description:** Free-text search across version, platform, and notes.

**Latest version**

- **Type:** Text
- **Required:** No
- **Description:** Used to compare or highlight a specific version string.

**Platform**

- **Type:** Dropdown
- **Required:** No
- **Description:** Filter by device type (Android or iOS).

**Enabled**

- **Type:** Checkbox / Toggle
- **Required:** No
- **Description:** Show only versions that are currently active.

**Force update**

- **Type:** Checkbox / Toggle
- **Required:** No
- **Description:** Show only versions that require users to update.

**Created at / Last modified / Updated at**

- **Type:** Date pickers
- **Required:** No
- **Description:** Restrict results to specific dates or periods.

**Updated by / Created by**

- **Type:** Text
- **Required:** No
- **Description:** Filter by the username or ID of the person who created or last updated the version.

**Filter (advanced)**

- **Type:** Text (advanced filter)
- **Required:** No
- **Description:** Used for more advanced filtering, according to internal rules.

**Example (Correctly Filled Form)**

Page: `1`  
Items per page: `10`  
Search: `2.3`  
Platform: `Android`  
Enabled: `On`  
Force update: `Off`

**Common Mistakes**

- Typing full sentences with many symbols in **Search**, which may not return good results. Keep it short and simple.
- Forgetting to clear filters and wondering why some versions “disappeared”.

---

### Screen / Form: Enable / Disable a Device Version

**Purpose:**  
To turn a given app version on or off for real users.

**Input Fields**

Often this is done via a simple toggle on the version details screen:

**Enabled**

- **Type:** Checkbox / Toggle
- **Required:** Yes
- **Description:** When turned on, this version is allowed to connect. When turned off, users on this version may not be allowed to continue.

**Example (Correctly Filled Form)**

Enabled: `Off` for version `2.1.0` (because of a critical bug)  
Enabled: `On` for version `2.3.0` (the fixed release)

**Common Mistakes**

- Disabling the only working version for a platform – always ensure at least one stable version remains enabled.
- Forgetting to coordinate with customer support before disabling a widely-used version, which can cause a flood of user complaints.

