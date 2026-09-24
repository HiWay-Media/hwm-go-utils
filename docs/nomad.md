---
layout: default
title: Nomad
nav_order: 5
description: "nomad — read job definitions and allocation stats, run, scale, restart and delete Nomad jobs."
permalink: /nomad
last_modified_date: 2026-09-24
---

# Nomad
{: .no_toc }

A focused client for the [HashiCorp Nomad HTTP API](https://developer.hashicorp.com/nomad/api-docs):
inspect, run, scale and stop jobs.
{: .fs-6 .fw-300 }

<dl class="glance">
  <dt>Import</dt><dd><code>github.com/HiWay-Media/hwm-go-utils/nomad</code></dd>
  <dt>API</dt><dd><code>/v1/…</code> endpoints of any Nomad agent (default port 4646)</dd>
  <dt>Auth</dt><dd>Optional ACL token (<code>X-Nomad-Token</code>)</dd>
</dl>

1. TOC
{:toc}

---

## Setup

```go
nc := nomad.NewService(nomad.Options{
	BaseUrl:    "http://nomad.service.consul:4646", // with or without a trailing /v1
	Token:      os.Getenv("NOMAD_TOKEN"),            // optional, for ACL-enabled clusters
	ScaleGroup: "encoder",                           // task group used by ScaleJob / RestartJob
	Logger:     logger,                              // optional, defaults to a no-op logger
	LogLevel:   "info",                              // "debug" logs every HTTP exchange
})
```

| Option | Default | Notes |
|:--|:--|:--|
| `BaseUrl` | — | `http://host:4646` and `http://host:4646/v1` are equivalent |
| `Token` | none | Sent as `X-Nomad-Token` on every request |
| `ScaleGroup` | `nomad.DefaultScaleGroup` (`"restreamer"`) | Task group targeted by scaling |
| `Logger` | no-op | `*zap.SugaredLogger` |
| `LogLevel` | — | `"debug"` enables resty request/response logging |

## Operations

| Method | Nomad endpoint | Returns |
|:--|:--|:--|
| `GetDefinition(jobID, region)` | `GET /v1/job/:id` | `*JobDefinition` |
| `RunJob(definition, region)` | `POST /v1/jobs` | error |
| `ScaleJob(jobID, count, region)` | `POST /v1/job/:id/scale` | error |
| `RestartJob(jobID, region)` | scale to 0, then back to 1 after one second | error (of the first step) |
| `DeleteJob(jobID, region, purge)` | `DELETE /v1/job/:id?purge=` | error |
| `AllocationStats(allocID, region)` | `GET /v1/client/allocation/:id/stats` | `*ResourceUsage` |
| `GetAllocations(nodeID, region)` | `GET /v1/node/:id/allocations` | `*NomadAllocations` |

Job, allocation and node ids are path-escaped, so ids containing spaces or slashes are safe.

Every call passes `region` as a query parameter; use `""` for the agent's default region.

## Examples

### Clone a job with a new name

```go
def, err := nc.GetDefinition("restreamer-template", "eu")
if err != nil {
	return err
}
def.ID, def.Name = "restreamer-match-42", "restreamer-match-42"
return nc.RunJob(*def, "eu")
```

### Scale out, then back

```go
if err := nc.ScaleJob("restreamer-match-42", 3, "eu"); err != nil {
	return err
}
// …
return nc.ScaleJob("restreamer-match-42", 1, "eu")
```

### List what runs on a node

```go
allocs, err := nc.GetAllocations(nodeID, "eu")
if err != nil {
	return err
}
for _, a := range allocs.NomadAllocations {
	fmt.Printf("%s  %s/%s\n", a.ID, a.JobID, a.TaskGroup)
}
```

### Read live resource usage

```go
usage, err := nc.AllocationStats(allocID, "eu")
if err != nil {
	return err
}
fmt.Printf("RSS %d MiB, CPU %.1f%%\n", usage.MemoryStats.RSS>>20, usage.CpuStats.Percent)
```

{: .note }
`RestartJob` returns as soon as the job is scaled to 0; the scale back to 1 happens in the
background one second later and its error is only logged. Call `ScaleJob` twice yourself
if you need to wait for, or check, the second step.
