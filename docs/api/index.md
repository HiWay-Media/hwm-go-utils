---
layout: default
title: HTTP API toolkit
nav_order: 3
has_children: true
description: "Fiber + GORM building blocks: generic CRUD, JWT middleware, REST client and response envelopes."
permalink: /api/
last_modified_date: 2026-09-24
---

# HTTP API toolkit

Four packages that cover the usual plumbing of a JSON API built on
[Fiber v2](https://docs.gofiber.io/) and [GORM](https://gorm.io/).
{: .fs-6 .fw-300 }

<div class="cards">
  <a class="card" href="{{ '/api/generic-crud' | relative_url }}">
    <div class="card-icon">⚡</div>
    <p class="card-title">Generic CRUD</p>
    <p class="card-text"><code>api/generic</code> — list, get, create and delete routes for any GORM model, with paging, validation and safe errors.</p>
  </a>
  <a class="card" href="{{ '/api/jwt-middleware' | relative_url }}">
    <div class="card-icon">🔐</div>
    <p class="card-title">JWT middleware</p>
    <p class="card-text"><code>api/middlewares</code> — Keycloak token verification, role checks and claim getters.</p>
  </a>
  <a class="card" href="{{ '/api/http-client' | relative_url }}">
    <div class="card-icon">📡</div>
    <p class="card-title">HTTP client</p>
    <p class="card-text"><code>api/client</code> — a tiny resty wrapper for calling other JSON services with a bearer token.</p>
  </a>
  <a class="card" href="{{ '/api/responses' | relative_url }}">
    <div class="card-icon">📦</div>
    <p class="card-title">Responses</p>
    <p class="card-text"><code>api/models</code> — the <code>OK</code> / <code>KO</code> envelopes every endpoint returns.</p>
  </a>
</div>

## Request flow

```mermaid
sequenceDiagram
    participant C as Client
    participant J as JwtProtected
    participant R as RoleCheck
    participant H as Handler
    participant S as Store
    C->>J: Bearer token
    J->>J: signature · exp · iss/aud
    J--xC: 401 invalid token
    J->>R: verified claims
    R--xC: 403 missing role
    R->>H: next()
    H->>S: Get(id) as bound param
    S-->>H: row · not found · error
    H-->>C: 200 · 404 · 500 internal error
```

Rejections (`--x`) stop the chain. Error details stay in the server log; the client only gets
the status and a generic message.
