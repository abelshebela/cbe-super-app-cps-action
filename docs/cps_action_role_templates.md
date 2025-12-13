# CPS Action Role Mapping Templates

This doc shows how to represent maker/checker/auditor role assignments per CPS action using Mongo Extended JSON, matching your example structure.

Each action entry uses:
- action_code: matches the RequestAction constant, e.g. CREATE_PUBLIC_NOTIFICATION
- action_name: human readable name
- assigned_makers_roles: array of ObjectIds for maker roles
- assigned_checkers_roles: array of approval levels; each level is an array of role ObjectIds
- assigned_auditor_roles: array of ObjectIds for auditor roles
- approver_count: number of checker approval levels (len(assigned_checkers_roles))
- is_maker_only: true to skip checker approval; when true, assigned_checkers_roles is [] and approver_count is 0
- enabled: whether the action mapping is active
- created_at/updated_at: ISO dates as Mongo Extended JSON

## Generic Template
```json
{
  "action_code": "<REQUEST_ACTION>",
  "action_name": "<Human Friendly Name>",
  "assigned_makers_roles": [
    { "$oid": "<role_id>" }
  ],
  "assigned_checkers_roles": [
    [ { "$oid": "<role_id>" } ],
    [ { "$oid": "<role_id>" }, { "$oid": "<role_id>" } ]
  ],
  "assigned_auditor_roles": [
    { "$oid": "<role_id>" }
  ],
  "approver_count": 2,
  "is_maker_only": false,
  "enabled": true,
  "updated_at": { "$date": "<ISO_DATETIME>" },
  "created_at": { "$date": "<ISO_DATETIME>" }
}
```

## Examples

### Notification: CREATE_PUBLIC_NOTIFICATION
```json
{
  "action_code": "CREATE_PUBLIC_NOTIFICATION",
  "action_name": "Create Public Notification",
  "assigned_makers_roles": [
    { "$oid": "692fe3f605f859e0f6133f01" }
  ],
  "assigned_checkers_roles": [
    [ { "$oid": "692fe3f605f859e0f6133f11" } ],
    [ { "$oid": "692fe3f605f859e0f6133f12" }, { "$oid": "692fe3f605f859e0f6133f13" } ]
  ],
  "assigned_auditor_roles": [
    { "$oid": "692fe3f605f859e0f6133faa" }
  ],
  "approver_count": 2,
  "is_maker_only": false,
  "enabled": true,
  "updated_at": { "$date": "2025-12-10T19:25:43.826Z" },
  "created_at": { "$date": "2025-12-09T09:30:00.000Z" }
}
```

### ServicesCatalog: CREATE_SERVICE
```json
{
  "action_code": "CREATE_SERVICE",
  "action_name": "Create Service",
  "assigned_makers_roles": [
    { "$oid": "692fe3f605f859e0f6134001" }
  ],
  "assigned_checkers_roles": [
    [ { "$oid": "692fe3f605f859e0f6134011" } ]
  ],
  "assigned_auditor_roles": [
    { "$oid": "692fe3f605f859e0f6133fab" }
  ],
  "approver_count": 1,
  "is_maker_only": false,
  "enabled": true,
  "updated_at": { "$date": "2025-12-10T19:25:43.826Z" },
  "created_at": { "$date": "2025-12-09T09:30:00.000Z" }
}
```

### Maker-only example (no checker levels)
```json
{
  "action_code": "MARK_NOTIFICATION_AS_SEEN",
  "action_name": "Mark Notification As Seen",
  "assigned_makers_roles": [
    { "$oid": "692fe3f605f859e0f6134101" }
  ],
  "assigned_checkers_roles": [],
  "assigned_auditor_roles": [],
  "approver_count": 0,
  "is_maker_only": true,
  "enabled": true,
  "updated_at": { "$date": "2025-12-10T19:25:43.826Z" },
  "created_at": { "$date": "2025-12-09T09:30:00.000Z" }
}
```

## Next step (optional)
I can auto-generate a full file with one entry for every RequestAction from code (using current RequestActionGroups), either as:
- JSON array: docs/cps_actions_roles_seed.json
- NDJSON: docs/cps_actions_roles_seed.ndjson (one object per line)

It will include empty role arrays and is_maker_only=false by default, so you can fill role ObjectIds and approval levels.
