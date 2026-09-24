---
layout: default
title: Responses
parent: HTTP API toolkit
nav_order: 4
description: "api/models — the standard OK / KO JSON envelopes."
permalink: /api/responses
last_modified_date: 2026-09-24
---

# Responses
{: .no_toc }

One envelope for every endpoint, so clients always know where to look.
{: .fs-6 .fw-300 }

<dl class="glance">
  <dt>Import</dt><dd><code>github.com/HiWay-Media/hwm-go-utils/api/models</code></dd>
</dl>

## Success

```json
{ "response": "OK", "message": "optional text", "data": { "id": 1 } }
```

| Constructor | Produces |
|:--|:--|
| `ApiDefault()` | `{"response":"OK"}` |
| `ApiDefaultResponse(data)` | `{"response":"OK","data":…}` |
| `ApiDefaultMsgOnly(msg)` | `{"response":"OK","message":"…"}` |
| `ApiDefaultMsgResponse(data, msg)` | `{"response":"OK","message":"…","data":…}` |

## Error

```json
{ "response": "KO", "message": "not found" }
```

| Constructor | Produces |
|:--|:--|
| `ApiDefaultError(msg)` | `{"response":"KO","message":"…"}` |

Empty fields are omitted from the JSON.

## In a handler

```go
func getChannel(c *fiber.Ctx) error {
	ch, err := repo.Find(c.Params("id"))
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return c.Status(fiber.StatusNotFound).JSON(models.ApiDefaultError("not found"))
	}
	if err != nil {
		logger.Errorf("find channel: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(models.ApiDefaultError("internal error"))
	}
	return c.JSON(models.ApiDefaultResponse(ch))
}
```

{: .tip }
Send users a stable, generic message and keep the real error in the logs — the same rule
the [generic CRUD](generic-crud) follows.
