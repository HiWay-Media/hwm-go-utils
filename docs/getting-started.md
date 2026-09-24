---
layout: default
title: Getting started
nav_order: 2
description: "Install hwm-go-utils and build a protected CRUD API in a few minutes."
permalink: /getting-started
last_modified_date: 2026-09-24
---

# Getting started
{: .no_toc }

Build a small product catalogue API: MySQL storage, Keycloak-protected routes and
an admin-only delete — about 40 lines of code.
{: .fs-6 .fw-300 }

1. TOC
{:toc}

---

## Requirements

| | |
|:--|:--|
| **Go** | 1.26 or later (`GOTOOLCHAIN=auto` downloads it for you) |
| **Web framework** | [Fiber v2](https://docs.gofiber.io/) for the `api/*` packages |
| **ORM** | [GORM](https://gorm.io/) with the MySQL driver for `api/generic` and `db` |
| **Identity provider** | Keycloak (or any issuer of RS256 JWTs) for `api/middlewares` |

Every other package — `nomad`, `nats_helper`, `redis`, `log`, `utils/*` — can be used on its own.

## 1. Install

```shell
go get github.com/HiWay-Media/hwm-go-utils@latest
```

Releases are tagged `vX.Y.Z` on `main`; pin one in `go.mod` for reproducible builds.

## 2. Create a logger and a database

```go
logger := log.GetLogger(os.Getenv("LOG_LEVEL")) // debug | info | warn | error — default info

database, err := db.Open(
	os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"),
	os.Getenv("DB_HOST"), 3306, "shop",
	"5",  // max idle connections ("" keeps the database/sql default)
	"20", // max open connections
)
if err != nil {
	logger.Fatalf("database: %v", err)
}
```

{: .important }
`db.Open` returns an error that names user, host and database but **never the password**,
so it is safe to log.

## 3. Describe your resource

The generic CRUD works with any GORM model. Validation uses
[`validator.v2`](https://pkg.go.dev/gopkg.in/validator.v2) tags.

```go
type Product struct {
	ID    uint    `json:"id" gorm:"primaryKey"`
	Name  string  `json:"name" validate:"nonzero"`
	Price float64 `json:"price" validate:"min=0"`
}
```

## 4. Protect the routes

`JwtProtected` verifies the bearer token against your realm's public key
(*Keycloak → Realm settings → Keys → RS256 → Public key*). Add `WithIssuer` so only
tokens from your realm are accepted.

```go
auth := middlewares.JwtProtected(os.Getenv("KEYCLOAK_PUBLIC_KEY"),
	middlewares.WithIssuer("https://sso.example.com/realms/shop"),
)
```

## 5. Register the endpoints

```go
app := fiber.New()
api := app.Group("/api", auth)

generic.SetEndpoints[Product]("products", api, database, logger)

logger.Fatal(app.Listen(":8080"))
```

That single call gives you:

| Method | Path | Behaviour |
|:--|:--|:--|
| `GET` | `/api/products?start=0&limit=50` | Paginated list (page size capped at 1000) |
| `GET` | `/api/products/:id` | One product, `404` if missing |
| `POST` | `/api/products` | Validates and creates; the client cannot choose the `id` |
| `DELETE` | `/api/products/:id` | Deletes, `404` if nothing matched |

## 6. Restrict a route to a role

Need finer control than the generated routes? `SetEndpoints` also returns the store,
service and handler, so you can mount them yourself:

```go
_, _, products := generic.SetEndpoints[Product]("products", api, database, logger)

admin := app.Group("/admin", auth, middlewares.RoleCheck([]string{"catalogue-admin"}))
admin.Delete("/products/:id", products.Delete)
```

A token without any of the listed realm roles gets `403 Forbidden`.

## 7. Try it

```shell
TOKEN=$(curl -s -d grant_type=client_credentials -d client_id=shop-cli -d client_secret=… \
  https://sso.example.com/realms/shop/protocol/openid-connect/token | jq -r .access_token)

curl -s -H "Authorization: Bearer $TOKEN" -H 'Content-Type: application/json' \
  -d '{"name":"Headphones","price":79.9}' http://localhost:8080/api/products
# {"response":"OK","data":{"id":1,"name":"Headphones","price":79.9}}

curl -s -H "Authorization: Bearer $TOKEN" http://localhost:8080/api/products
# {"response":"OK","data":[{"id":1,"name":"Headphones","price":79.9}]}
```

## Next steps

- [Generic CRUD in depth](api/generic-crud) — paging, errors, custom routes
- [JWT middleware](api/jwt-middleware) — options, claims, response codes
- [Keycloak](keycloak) — manage users and realms from Go
- [Security model](security) — what the library guarantees for you
