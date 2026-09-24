---
layout: default
title: JWT middleware
parent: HTTP API toolkit
nav_order: 2
description: "api/middlewares — Keycloak JWT verification, role checks and claim getters for Fiber."
permalink: /api/jwt-middleware
last_modified_date: 2026-09-24
---

# JWT middleware
{: .no_toc }

Verify Keycloak access tokens and authorise by realm role — without ever crashing on an
unexpected token.
{: .fs-6 .fw-300 }

<dl class="glance">
  <dt>Import</dt><dd><code>github.com/HiWay-Media/hwm-go-utils/api/middlewares</code></dd>
  <dt>Algorithms</dt><dd>RS256, RS384, RS512 only</dd>
  <dt>Checks</dt><dd>signature · <code>exp</code> (required) · <code>iss</code> and <code>aud</code> (opt-in) · realm roles</dd>
</dl>

1. TOC
{:toc}

---

## Usage

```go
auth := middlewares.JwtProtected(publicKey,
	middlewares.WithIssuer("https://sso.example.com/realms/media"),
	middlewares.WithAudience("media-api"),
)

app.Get("/me", auth, func(c *fiber.Ctx) error {
	claims, _ := middlewares.GetTokenClaims(c)
	return c.JSON(fiber.Map{
		"user":     middlewares.GetCurrentLoggedUserUUIDFromJwt(claims),
		"email":    middlewares.GetCurrentLoggedUserEmailFromJwt(claims),
		"customer": middlewares.GetCustomerIdFromJwt(claims),
	})
})

app.Delete("/channels/:id", auth, middlewares.RoleCheck([]string{"admin", "operator"}), deleteChannel)
```

### The public key

Copy it from *Keycloak → Realm settings → Keys → RS256 → Public key*. Both forms work:

- the bare base64 string shown by the console (`MIIBIjANBgkqh…`);
- a full PEM block (`-----BEGIN PUBLIC KEY----- …`).

The key is parsed **once**, when `JwtProtected` is called. If it is invalid, the problem is
logged at startup and every request is rejected with `401` — the middleware fails closed.

## Options

| Option | Checks | When to use |
|:--|:--|:--|
| `WithIssuer(iss)` | `iss` equals `iss` | Always in production: rejects tokens minted by other realms that share a key |
| `WithAudience(aud)` | `aud` contains `aud` | When the Keycloak client has an *audience* mapper for your API |

{: .tip }
By default Keycloak access tokens often carry only `aud: "account"`. Add an **Audience**
protocol mapper to the client before enabling `WithAudience`, or every token will be refused.

## What gets rejected

| Token | Result |
|:--|:--|
| Missing `Authorization` header, missing `Bearer ` prefix, empty token | `401 malformed token` |
| Bad signature, signed with another key | `401 invalid token` |
| `alg: none`, HMAC algorithms (including HS256 signed with the public key) | `401 invalid token` |
| Expired, or **no `exp` claim at all** | `401 invalid token` |
| Wrong issuer / audience (when configured) | `401 invalid token` |

The reason is written to the server log; the client only sees `invalid token`.

## Roles

`RoleCheck(roles)` lets the request through when the token's `realm_access.roles`
contains **at least one** of `roles`.

| Situation | Status |
|:--|:--|
| At least one role present | handler runs |
| None present, or no `realm_access` claim | `403 Missing role …` |
| `RoleCheck` used without `JwtProtected` before it | `401 unauthorized` |
| Empty role list | handler runs |

## Reading claims

`GetTokenClaims(c)` returns the verified `jwt.MapClaims`. The getters never panic: a missing
or mistyped claim returns the zero value.

| Getter | Claim | Missing → |
|:--|:--|:--|
| `GetCurrentLoggedUserUUIDFromJwt` | `sub` | `""` |
| `GetCurrentLoggedUserEmailFromJwt` | `email` | `""` |
| `GetCustomerNameFromJwt` | `customer_name` | `""` |
| `GetCustomerIdFromJwt` | `customer_id` | `0` |
| `GetClientIdFromJwt` | `customer_name` + `_client` | `"_client"` |
| `GetRolesListFromJwt` | `realm_access.roles` | error |

`customer_name` and `customer_id` are custom claims: add them with user-attribute mappers in
Keycloak.
