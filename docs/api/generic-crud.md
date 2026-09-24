---
layout: default
title: Generic CRUD
parent: HTTP API toolkit
nav_order: 1
description: "api/generic — list, get, create and delete endpoints for any GORM model."
permalink: /api/generic-crud
last_modified_date: 2026-09-24
---

# Generic CRUD
{: .no_toc }

Turn any GORM model into a REST resource with one call.
{: .fs-6 .fw-300 }

<dl class="glance">
  <dt>Import</dt><dd><code>github.com/HiWay-Media/hwm-go-utils/api/generic</code></dd>
  <dt>Built on</dt><dd>Fiber v2, GORM, validator.v2</dd>
  <dt>Layers</dt><dd><code>Handler[T]</code> → <code>Service[T]</code> → <code>Store[T]</code>, each behind an interface you can replace</dd>
</dl>

1. TOC
{:toc}

---

## Quick start

```go
type Channel struct {
	ID     uint   `json:"id" gorm:"primaryKey"`
	Name   string `json:"name" validate:"nonzero"`
	Region string `json:"region" validate:"regexp=^(eu|us|ap)$"`
}

store, service, handler := generic.SetEndpoints[Channel]("channels", router, database, logger)
```

`SetEndpoints` builds the three layers and registers the routes on `router` (a
`*fiber.App` or any `fiber.Router` group):

| Method | Path | Handler | Success |
|:--|:--|:--|:--|
| `GET` | `/channels?start=&limit=` | `List` | `200` + array |
| `GET` | `/channels/:id` | `Get` | `200` + object |
| `POST` | `/channels` | `Create` | `200` + created object (with its new `id`) |
| `DELETE` | `/channels/:id` | `Delete` | `200` |

All responses use the [standard envelope](responses): `{"response":"OK","data":…}` or
`{"response":"KO","message":"…"}`.

## Listing and paging

| Query | Default | Rules |
|:--|:--|:--|
| `start` | `0` | Offset. Negative or non-numeric → `400 start invalid` |
| `limit` | `0` | Page size. `0` or anything above `generic.MaxListLimit` is clamped to it. Negative → `400 limit invalid` |

```go
// default is 1000 rows per page; raise, lower or disable (0) the cap at startup
generic.MaxListLimit = 200
```

{: .warning }
In v0.6.x a list request without `limit` returned an **empty array** (GORM rendered
`LIMIT 0`). Since v0.7.0 it returns up to `MaxListLimit` rows. Set
`generic.MaxListLimit = 0` for unbounded lists.

## Creating

`POST` parses the JSON body into `T`, runs `validator.Validate`, then calls `Service.Create`.

Before saving, the handler **clears the fields the server owns**:

- the primary key(s) — clients cannot choose or overwrite an `id`;
- relationship fields (`has one`, `has many`, `belongs to`, `many2many`) — a request
  cannot make GORM create or link rows in other tables.

Plain foreign-key columns such as `OwnerID uint` are kept, so you can still link to an
existing row explicitly. Hooks like `BeforeCreate` that generate UUIDs keep working,
because they run after the fields are cleared.

## Getting and deleting

The `:id` path parameter is passed to the store as an `int` when it is numeric and as a
`string` otherwise (UUIDs, slugs). Either way it is **bound as a query parameter against
the model's primary key** — it is never interpolated into SQL.

`Delete` returns `404` when no row matched.

## Errors

| Situation | Status | Body `message` |
|:--|:--|:--|
| Malformed JSON, failed validation | `400` | the parser / validator message |
| Invalid `start`, `limit` or empty `:id` | `400` | `start invalid`, `limit invalid`, `id invalid` |
| `gorm.ErrRecordNotFound` | `404` | `not found` |
| Any other error | `500` | `internal error` |

{: .important }
Database errors are **logged** (method, path and error) through the logger you passed,
and never sent to the client — no table names, SQL fragments or credentials leak.

## Custom routes and layers

Every layer is an interface, so you can wrap or replace one and keep the rest:

```go
type IStore[T any] interface {
	Get(id any) (*T, error)
	Create(obj *T) error
	Delete(id any) error
	List(start, limit int) ([]T, error)
}
```

For example, scope everything to the caller's tenant with your own store, then reuse the
generic service and handler:

```go
type tenantStore struct {
	generic.IStore[Channel]
	db *gorm.DB
}

func (s tenantStore) List(start, limit int) ([]Channel, error) {
	var out []Channel
	err := s.db.Where("tenant_id = ?", currentTenant()).Offset(start).Limit(limit).Find(&out).Error
	return out, err
}

// override Get, Create and Delete the same way, or unscoped calls fall through
// to the embedded generic store

store := tenantStore{IStore: generic.NewStore[Channel](database), db: database}
handler := generic.NewHandler[Channel](generic.NewService[Channel](store, logger), logger)

g := router.Group("/channels")
g.Get("/", handler.List)
g.Get("/:id", handler.Get)
```

## Typed query and path parameters

Two small helpers convert parameters with generics:

```go
page, err := generic.FromQuery[int](c, "page")    // ?page=3 → *int(3)
id, err := generic.FromParam[string](c, "id")     // /:id    → *string
```
