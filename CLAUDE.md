# CLAUDE.md — hwm-go-utils

Libreria Go condivisa di HiWay Media (`github.com/HiWay-Media/hwm-go-utils`, modulo omonimo), importata dai servizi Go interni tramite tag `vX.Y.Z`. Package: `api/` (client resty, CRUD generico Fiber+GORM, middleware JWT Keycloak, modelli risposta), `keycloak/` (wrapper gocloak v10), `nomad/` (client API Nomad), `nats_helper/` (connessione NATS + JetStream), `db/` (GORM MySQL), `redis/` + `keydb/` (go-redis v8), `log/` (zap), `utils/*` (helper generici). Docs Jekyll in `docs/` (GitHub Pages).

## Regole di lavoro (SEMPRE)

- **È una libreria**: ogni cambio di firma pubblica o di comportamento rompe i servizi che la importano. Preferire aggiunte retro-compatibili; se un breaking change è inevitabile, dichiararlo esplicitamente nella PR (sezione "Cambi di comportamento").
- **Release = tag `vX.Y.Z` su `main` dopo il merge** (serie attuale `v0.6.x`, ultimo `v0.6.80`). Il tag lo crea/pusha l'utente, mai su branch di feature. `patch` per fix, `minor` per nuovi package/API.
- **Branch + PR**, mai commit diretti su `main`. `git push` solo quando l'utente chiede esplicitamente di aprire la PR. MAI `Co-Authored-By` né footer di attribuzione in commit/PR.
- **Gate prima di ogni commit**: `go build ./... && go vet ./... && go test -race ./...` verdi, `go mod tidy` senza diff. `gofmt` sui file toccati (molti file legacy non sono formattati: non riformattare file interi non correlati, gonfia il diff).
- **Ogni fix ha un test di regressione** che fallisce sul codice vecchio. Niente test che richiedono servizi reali senza `t.Skip` se manca l'env (vedi `keycloak/*_test.go` con `KEYCLOAK_SERVER`).
- Libreria ≠ applicazione: niente `log.Fatal`/`os.Exit`/`fmt.Println` nel codice nuovo — restituire `error` e loggare col `*zap.SugaredLogger` passato dal chiamante.
- **Mai segreti nei log**: DSN con password, token, client secret (anche in debug).

## Pattern per test senza infrastruttura (validato, PR audit 2026-09)

- **GORM**: `mysql.New(mysql.Config{SkipInitializeWithVersion: true})` + `gorm.Config{DryRun: true, DisableAutomaticPing: true, SkipDefaultTransaction: true}`; catturare l'SQL con una callback `After("gorm:query")` / `After("gorm:delete")` e `Dialector.Explain`. Esempio: `api/generic/store_test.go`.
- **HTTP (Keycloak, Nomad)**: `httptest.NewServer` che registra metodo+path o restituisce token finti. Esempi: `keycloak/token_test.go`, `nomad/service_test.go`.

## Trappole note / regole tecniche

- **GORM inline conditions**: `db.First(&t, id)` / `db.Delete(&t, id)` con `id` stringa non numerica = **SQL grezzo** (SQL injection). Usare sempre `Where(clause.Eq{...})` o `Where("col = ?", v)`. `Limit(0)` genera `LIMIT 0`: per "nessun limite" usare `Limit(-1)`.
- **Keycloak**: il token admin (client_credentials) va preso con `g.adminToken()` — cache + refresh 30s prima della scadenza, sotto `Mu`. Mai usare `g.adminJWT` direttamente.
- **Nomad**: tutti i path hanno prefisso `/v1/...`; `apiBaseURL()` toglie `/v1` finale dal base URL, quindi funziona con entrambe le convenzioni. ID nel path sempre con `url.PathEscape`. API `:4646` senza ACL token (non supportato dal client). `ScaleJob` ha il gruppo `"restreamer"` hardcoded.
- **NATS**: `MaxReconnects(-1)` obbligatorio — col default (60) la connessione si chiude per sempre dopo ~1 min di server giù. `nats.EncodedConn` è deprecato (migrazione = breaking change).
- **Fiber**: da v2.50 `c.GetReqHeaders()` restituisce `map[string][]string` → l'upgrade rompe `api/middlewares`.
- **Middleware JWT**: le `Get*FromJwt` e `RoleCheck` fanno type assertion senza check → panic se il claim manca (Fiber non ha recover di default). Non verifica `iss`/`aud`.
- CI GitHub Actions: `go-test.yml` gira solo su `push` e fa `go mod tidy` (muta il repo); matrice Go 1.20–1.22 (EOL). `go.mod` dichiara `go 1.20`.

## Debito noto (dall'audit 2026-09-24, non ancora risolto)

- Dipendenze con CVE (`govulncheck -show verbose ./...`): fiber 2.46, fasthttp, jwt/v4 4.4.3, x/crypto, x/net, x/text — nessuna raggiungibile a livello simboli, ma da aggiornare insieme al middleware (vedi Fiber sopra).
- `api/middlewares`: `iss`/`aud`, type assertion sicure, 403 invece di 401 per ruolo mancante, non esporre `err.Error()` al client.
- `db.InitDB` logga la DSN con password e usa `Fatalf`; `utils/file.FileExists` va in panic su errori ≠ not-exist; `utils/strings.EncodeURL` codifica due volte; `log.GetLogger` emette codici ANSI e ignora l'encoder JSON.
- CRUD generico: nessun limite massimo su `List`, mass assignment su `Create`, errori DB restituiti al client.
- `dependabot.yml` punta a `/tests` (inesistente); `.gitignore` è un template Java.

## Puntatori

- Docs: `docs/` (Jekyll, workflow `jekyll-gh-pages.yml`) — aggiornare `docs/<package>.md` se cambia l'API pubblica.
- Repo affini: `devops_hiway` (infra, NATS/Nomad in produzione — vedi il suo `CLAUDE.md` per regole prod).
