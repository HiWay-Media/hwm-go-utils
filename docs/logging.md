---
layout: default
title: Logging
nav_order: 9
description: "log — a zap console logger configured from a level string."
permalink: /logging
last_modified_date: 2026-09-24
---

# Logging
{: .no_toc }

A [zap](https://github.com/uber-go/zap) logger in one line, with the level taken from
configuration.
{: .fs-6 .fw-300 }

<dl class="glance">
  <dt>Import</dt><dd><code>github.com/HiWay-Media/hwm-go-utils/log</code></dd>
  <dt>Returns</dt><dd><code>*zap.SugaredLogger</code>, the type every other package accepts</dd>
</dl>

## Usage

```go
logger := log.GetLogger(os.Getenv("LOG_LEVEL"))

logger.Infow("channel started", "channel", 42, "region", "eu")
logger.Errorf("upstream failed: %v", err)
```

Output (stderr, console format, plain-text levels — no ANSI colour codes, so it reads well in
Loki, CloudWatch or `docker logs`):

```text
2026-09-24T10:17:21.770+0200	INFO	player/main.go:31	channel started	{"channel": 42, "region": "eu"}
2026-09-24T10:17:22.104+0200	ERROR	player/main.go:40	upstream failed: context deadline exceeded
```

## Levels

| `logLevel` (case-insensitive) | Level |
|:--|:--|
| `debug` | Debug |
| `info`, empty or unknown | Info |
| `warn` | Warn |
| `error` | Error |
| `fatal` | Fatal |
| `panic` | Panic |

{: .tip }
Pass the same logger to `generic.SetEndpoints`, `nomad.Options`, `nats_helper` and the HTTP
client, so every package writes to one place with one level.
