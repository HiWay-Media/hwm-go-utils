---
layout: default
title: Upgrade guide
nav_order: 12
description: "What changed in v0.7.0 and how to upgrade from v0.6.x."
permalink: /changelog
last_modified_date: 2026-09-24
---

# Upgrade guide
{: .no_toc }

Everything you need to move a service from `v0.6.x` to `v0.7.0`.
{: .fs-6 .fw-300 }

1. TOC
{:toc}

---

## v0.7.0

The result of a full security and reliability audit. Most changes are fixes you get for free;
the ones below can change what your service observes.

### Before you upgrade

```shell
go get github.com/HiWay-Media/hwm-go-utils@v0.7.0
go mod tidy
go build ./... && go test ./...
```

{: .warning }
**Go 1.26 or later is required.** Update the `go` directive of your module and the Go version
of your Docker images / CI (or build with `GOTOOLCHAIN=auto`).

### Behaviour changes

| Package | Before (v0.6.x) | Now (v0.7.0) | Action |
|:--|:--|:--|:--|
| `api/generic` | `GET /resource` without `limit` returned `[]` | Returns up to `MaxListLimit` (1000) rows | Pass `limit`, or set `generic.MaxListLimit` |
| `api/generic` | Client could send `id` and nested associations on `POST` | Both are cleared before saving | Link rows through FK columns (`OwnerID`) |
| `api/generic` | `DELETE` on a missing id → `200` | `404 not found` | Handle `404` in clients |
| `api/generic` | `500` bodies contained the database error | `internal error` (error in logs) | Read server logs |
| `api/generic` | Negative `start` / `limit` accepted | `400` | — |
| `api/middlewares` | Missing role → `401` | `403` | Handle `403` in clients |
| `api/middlewares` | Tokens without `exp` accepted | Rejected | — (Keycloak always sets `exp`) |
| `api/middlewares` | `401` body echoed the parser error | `invalid token` | Read server logs |
| `api/middlewares` | Missing claim → panic | Zero value / `403` | — |
| `keycloak` | `SetPassword` ignored its `realm` argument | Uses it (empty = client realm) | Pass `""` or the right realm |
| `utils/strings` | `EncodeURL` double-encoded values | Encoded once | Remove any manual decoding workaround |
| `utils/strings` | `RemoveSlice` reused the input array | Returns a copy | — |
| `utils/file` | `WriteToFile` exited the process on error | Returns `error` | Check the error |
| `log` | Levels had ANSI colour codes | Plain text | Update log parsers matching colour codes |

### New features

- `middlewares.WithIssuer`, `middlewares.WithAudience`, `middlewares.GetTokenClaims`
- `generic.MaxListLimit`
- `db.Open` — returns an error instead of exiting
- `nomad.Options.Token` (ACL) and `nomad.Options.ScaleGroup`
- `keydb.NewKeyDBClient(addr, password, db)`, `REDIS_PASSWORD` for the cluster client

### Fixes

- SQL injection through non-numeric ids in the generic store
- Keycloak admin methods crashing (token never stored); token now refreshed automatically
- Nomad scale/restart calling a wrong URL; mixed `/v1` prefixes; unescaped ids
- Nomad `GetAllocations` failing on every call (Nomad returns an array, the client expected an object)
- NATS connection giving up after 60 reconnect attempts
- Database password in logs; `FileExists` panic on permission errors
- 52 of 53 dependency vulnerabilities (see [Security model](security))

## v0.6.x

Earlier releases are listed on the
[GitHub tags page](https://github.com/HiWay-Media/hwm-go-utils/tags).
