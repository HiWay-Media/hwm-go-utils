---
layout: default
title: API go
nav_order: 2
description: "HWM Go Utils API Docs"
permalink: /api
last_modified_date: 2026-09-24T12:00:00+0000
---

# HWM API Documentation

## JWT middleware (Keycloak)

`JwtProtected` verifies the bearer token signature (RS256/384/512 only) and
expiry against the realm public key, then stores the claims for later handlers.
The key can be the bare base64 string shown by the Keycloak console or a full
PEM block. Issuer and audience checks are optional but recommended.

```go
import "github.com/HiWay-Media/hwm-go-utils/api/middlewares"

auth := middlewares.JwtProtected(publicKey,
	middlewares.WithIssuer("https://sso.example.com/realms/my-realm"),
	middlewares.WithAudience("my-api"),
)

app.Get("/admin", auth, middlewares.RoleCheck([]string{"admin"}), func(c *fiber.Ctx) error {
	claims, _ := middlewares.GetTokenClaims(c)
	return c.SendString(middlewares.GetCurrentLoggedUserEmailFromJwt(claims))
})
```

Responses: `401` for a missing/invalid token, `403` when none of the roles
passed to `RoleCheck` is in `realm_access.roles`. Claim getters return zero
values when a claim is missing.