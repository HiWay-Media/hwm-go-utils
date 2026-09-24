---
layout: default
title: Contributing
nav_order: 13
description: "Development workflow, test patterns and CI for hwm-go-utils."
permalink: /contributing
last_modified_date: 2026-09-24
---

# Contributing
{: .no_toc }

hwm-go-utils is imported by many services: every change is a change to all of them.
{: .fs-6 .fw-300 }

1. TOC
{:toc}

---

## Workflow

1. Branch from `main` and open a pull request — no direct commits.
2. Keep public APIs backward compatible. Add options as `Options` fields or functional
   options; add new functions instead of changing signatures.
3. If a behaviour change is unavoidable, list it under **Behaviour changes** in the PR.
4. Maintainers tag releases (`vX.Y.Z`) on `main` after merging.

## Before pushing

```shell
gofmt -l .                 # must print nothing
go vet ./...
go mod tidy -diff          # must print nothing
go test -race ./...
```

CI runs the same checks on Go 1.26 and 1.27 for every push and pull request.

## Tests without infrastructure

Every fix comes with a regression test that fails on the old code. None of them needs a real
database, Keycloak or Nomad:

| Target | Technique | Example |
|:--|:--|:--|
| GORM queries | Dry-run MySQL dialector + callback capturing the generated SQL | `api/generic/store_test.go` |
| HTTP handlers | `fiber.App.Test` with a fake service | `api/generic/handlers_test.go` |
| JWT | RSA key generated in the test, tokens signed per case (`alg: none`, HS256…) | `api/middlewares/middlewares_test.go` |
| Keycloak, Nomad | `httptest.NewServer` recording paths or returning fake tokens | `keycloak/token_test.go`, `nomad/service_test.go` |

Tests that need a real service must `t.Skip` when their environment variables are missing.

## Library rules

- No `log.Fatal`, `os.Exit` or `fmt.Println` in library code: return errors and log through
  the caller's `*zap.SugaredLogger`.
- Never log secrets: DSNs, tokens, client secrets — not even at debug level.
- Document public API changes in `docs/`. This site is built from `docs/` with Jekyll and
  [Just the Docs](https://just-the-docs.com/) on every push to `main`.

### Previewing the docs

```shell
docker run --rm -v "$PWD":/github/workspace -e GITHUB_WORKSPACE=/github/workspace \
  -e INPUT_SOURCE=./docs -e INPUT_DESTINATION=./_site -e INPUT_FUTURE=false \
  -e INPUT_VERBOSE=false -e INPUT_TOKEN= -e INPUT_BUILD_REVISION=local \
  ghcr.io/actions/jekyll-build-pages:v1.0.13
# then serve _site under /hwm-go-utils/
```
