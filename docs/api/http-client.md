---
layout: default
title: HTTP client
parent: HTTP API toolkit
nav_order: 3
description: "api/client — a minimal resty-based JSON client with bearer authentication."
permalink: /api/http-client
last_modified_date: 2026-09-24
---

# HTTP client
{: .no_toc }

Call another JSON service in two lines, with a bearer token and query parameters.
{: .fs-6 .fw-300 }

<dl class="glance">
  <dt>Import</dt><dd><code>github.com/HiWay-Media/hwm-go-utils/api/client</code></dd>
  <dt>Built on</dt><dd><a href="https://github.com/go-resty/resty">resty v2</a></dd>
</dl>

## Usage

```go
svc := client.NewService(logger, "https://billing.internal/api")
svc.SetAuthorization(accessToken) // sent as "Authorization: Bearer …" on every request

resp, err := svc.Get("/invoices", map[string]string{"customer": "42", "status": "open"})
if err != nil {
	return err // network error
}
if resp.IsError() {
	return fmt.Errorf("billing: %s", resp.Status())
}

var invoices []Invoice
if err := json.Unmarshal(resp.Body(), &invoices); err != nil {
	return err
}
```

## Methods

| Method | Body | Notes |
|:--|:--|:--|
| `Get(route, params)` | — | `params` become the query string |
| `Delete(route, params)` | — | |
| `Post(route, body, params)` | JSON | `Content-Type: application/json` is set when `body != nil` |
| `Put(route, body, params)` | JSON | |
| `SetAuthorization(token)` | — | Pass `""` to stop sending the header |

All methods return the raw `*resty.Response`: check `resp.IsError()` / `resp.StatusCode()`
yourself — an HTTP `4xx`/`5xx` is **not** returned as a Go error.

{: .note }
At `debug` level the client logs request and response bodies. Keep production loggers at
`info` or above if those bodies can contain personal data or secrets.
