---
layout: default
title: Security model
nav_order: 11
description: "What hwm-go-utils guarantees about injection, authentication, error handling and secrets."
permalink: /security
last_modified_date: 2026-09-24
---

# Security model
{: .no_toc }

The guarantees the library gives you, and the ones it leaves to your service.
{: .fs-6 .fw-300 }

1. TOC
{:toc}

---

## Guarantees

| Area | Guarantee | Where |
|:--|:--|:--|
| **SQL injection** | Ids from the URL are bound as parameters against the primary key, never used as SQL | `api/generic` |
| **Mass assignment** | `Create` clears primary keys and relationship fields from request bodies | `api/generic` |
| **Unbounded reads** | List pages are capped (`MaxListLimit`, default 1000) | `api/generic` |
| **Error leakage** | Database errors are logged, clients get `internal error` / `not found` | `api/generic` |
| **Token forgery** | Only RS256/384/512 signatures from the configured key; `alg: none` and HMAC confusion rejected | `api/middlewares` |
| **Token lifetime** | `exp` is required — tokens without it are refused | `api/middlewares` |
| **Cross-realm tokens** | Issuer and audience pinning with `WithIssuer` / `WithAudience` | `api/middlewares` |
| **Crash resistance** | Missing or mistyped claims never panic; invalid keys fail closed with `401` | `api/middlewares` |
| **Secrets in logs** | DSNs and admin tokens are never logged | `db`, `keycloak` |
| **Dependencies** | Fiber, fasthttp, golang-jwt and `golang.org/x/*` on patched versions; checked with `govulncheck` | `go.mod` |

## Your responsibilities

- **Enable `WithIssuer` in production.** Without it, any token signed by the same key is accepted.
- **Authorise, not just authenticate.** `JwtProtected` proves who the caller is; use `RoleCheck`
  or your own checks to decide what they may do. The generic CRUD has no per-row ownership —
  wrap the store for multi-tenant data (see [custom layers](api/generic-crud#custom-routes-and-layers)).
- **Keep debug logging off in production.** At `debug`, HTTP clients log full request and
  response bodies.
- **Terminate TLS** in front of the service, and use TLS URLs for Keycloak, NATS and Redis.

## Audit history

A full audit in September 2026 fixed, among others:

| Severity | Finding | Fixed in |
|:--|:--|:--|
| Critical | Non-numeric ids passed to GORM as raw SQL: `DELETE /resource/<crafted id>` could empty a table | v0.7.0 |
| Critical | Keycloak admin token never stored: every admin call crashed | v0.7.0 |
| High | JWT middleware panicked on tokens without expected claims; no issuer/audience checks | v0.7.0 |
| High | 53 known vulnerabilities in dependencies, 9 in imported packages | v0.7.0 |
| High | NATS connection closed permanently after ~1 minute of server downtime | v0.7.0 |
| Medium | Database password written to logs on connection failure | v0.7.0 |
| Medium | Database errors returned to API clients; client-chosen ids on create | v0.7.0 |

Remaining known advisory: `GO-2026-5932` in `golang.org/x/crypto/openpgp` — no upstream fix, and
the package is not imported by this module.

## Reporting a vulnerability

Please **do not open a public issue**: contact the HiWay Media maintainers privately
(for example through the [organisation profile](https://github.com/HiWay-Media)), with a
description and, if possible, a reproduction.
