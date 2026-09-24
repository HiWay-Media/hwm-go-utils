---
layout: default
title: Database
nav_order: 7
description: "db — open a MySQL connection with GORM and sensible pool limits."
permalink: /database
last_modified_date: 2026-09-24
---

# Database
{: .no_toc }

Open a MySQL connection through GORM in one call, with connection-pool limits and errors that
are safe to log.
{: .fs-6 .fw-300 }

<dl class="glance">
  <dt>Import</dt><dd><code>github.com/HiWay-Media/hwm-go-utils/db</code></dd>
  <dt>Built on</dt><dd>GORM + <code>gorm.io/driver/mysql</code></dd>
</dl>

## `Open` (recommended)

```go
database, err := db.Open(
	"app",           // user
	os.Getenv("DB_PASSWORD"),
	"mysql.internal", 3306,
	"shop",          // database
	"5",             // max idle connections
	"20",            // max open connections
)
if err != nil {
	logger.Errorf("database unavailable: %v", err)
	return err
}
```

| Parameter | Notes |
|:--|:--|
| pool limits | Strings, to pass environment variables straight through. Empty, non-numeric or `≤ 0` values are ignored and the `database/sql` defaults apply |
| DSN | Built as `user:pass@tcp(host:port)/db?parseTime=true` |

{: .important }
Errors look like `connect to app@mysql.internal:3306/shop: dial tcp …`: the DSN, and
therefore the password, is **never** part of the error, so it is safe to log.

## `InitDB` (legacy)

```go
database := db.InitDB(logger, user, password, host, 3306, "shop", "5", "20")
```

Same as `Open`, but calls `logger.Fatalf` on failure and stops the process. It is kept for
existing services; prefer `Open` in new code so the caller decides how to fail.
