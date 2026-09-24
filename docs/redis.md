---
layout: default
title: Redis & KeyDB
nav_order: 8
description: "redis and keydb — go-redis cluster and single-node clients."
permalink: /redis
last_modified_date: 2026-09-24
---

# Redis & KeyDB
{: .no_toc }

Ready-made [go-redis v8](https://github.com/redis/go-redis) clients for a Redis cluster
configured from the environment, and for single-node Redis or KeyDB.
{: .fs-6 .fw-300 }

## Redis cluster — `redis`

```go
import hwmredis "github.com/HiWay-Media/hwm-go-utils/redis"

hwmredis.Init()                 // reads REDIS_1 … REDIS_6
rdb := hwmredis.GetClient()     // *redis.ClusterClient
defer rdb.Close()

if err := rdb.Set(ctx, "session:42", payload, time.Hour).Err(); err != nil {
	return err
}
```

| Variable | Meaning |
|:--|:--|
| `REDIS_1` … `REDIS_6` | Cluster node addresses (`host:port`); unset variables are skipped |
| `REDIS_PASSWORD` | Optional password (`requirepass` / ACL default user) |

`Init` fills the exported `Brokers` slice; you can also set `hwmredis.Brokers` yourself
before calling `GetClient`.

## Single node / KeyDB — `keydb`

```go
import "github.com/HiWay-Media/hwm-go-utils/keydb"

rdb := keydb.NewKeyDBClient("keydb.internal:6379", os.Getenv("KEYDB_PASSWORD"), 0)
defer rdb.Close()

hits, err := rdb.Incr(ctx, "stats:hits").Result()
```

| Function | Auth | DB |
|:--|:--|:--|
| `NewKeyDBClient(addr, password, db)` | password | chosen |
| `GetKeyDBClient(addr)` | none | `0` |

Both return a standard `*redis.Client`, so the whole go-redis API is available.
