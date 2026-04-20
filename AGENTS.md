# ChartMogul Go SDK

Go SDK wrapping the ChartMogul API. Single flat package `chartmogul`. Uses `gorequest` for HTTP with automatic retry and exponential backoff via `cenkalti/backoff`.

**Note:** ChartMogul's nginx is incompatible with Go's HTTP/2. Run with `GODEBUG=http2client=0`.

## Commands

```bash
go test -v -timeout=10m ./...                        # run full test suite
go test -v -run TestRetrieveCustomer ./...            # run single test
go test -v -timeout=10m ./integration_tests/...       # run integration tests only
./bin/golangci-lint run                               # lint
./genmocks.sh                                         # regenerate mocks
```

Version is in `version.go`.

## Architecture

### Flat package structure

All source lives in the root package `chartmogul` - no subdirectories for library code. Each API resource gets its own file pair: `<resource>.go` + `<resource>_test.go`.

### Client setup

```go
type API struct {
    ApiKey string
    Client *http.Client   // nil = default; set for tests/VCR
    Retry  *RetryConfig   // nil = enabled with exponential backoff
}
```

Global helpers: `SetURL()` (test overrides), `Setup()` (timeout).

### Generic CRUD helpers (`generic.go`)

All resource methods delegate to shared helpers that handle retry, serialization, and error wrapping:

- `api.create(path, input, output)` - POST
- `api.retrieve(path, uuid, output)` - GET single (replaces `:uuid` in path)
- `api.list(path, output, query...)` - GET collection with optional query params
- `api.update(path, uuid, input, output)` - PATCH
- `api.putTo(path, uuid, input, output)` - PUT
- `api.add(path, uuid, input, output)` - POST to sub-resource
- `api.delete(path, uuid)` - DELETE
- `api.deleteWhat(path, uuid, input)` - DELETE with query params
- `api.deleteWithData(path, uuid, input)` - DELETE with body
- `api.merge(path, input)` - POST merge

Path placeholders (`:uuid`, `:customerUUID`, `:dataSourceUUID`, etc.) are replaced via `strings.Replace` before the request is sent.

### Resource file pattern

Every resource file follows the same layout:

1. **Endpoint constants** at top:
```go
const (
    customersEndpoint       = "customers"
    singleCustomerEndpoint  = "customers/:uuid"
)
```

2. **Type definitions** with JSON tags:
   - `<Resource>` - API response model (GET/list result)
   - `New<Resource>` - creation input (POST body)
   - `Update<Resource>` - update input (PATCH body), uses `*string`/`*bool` pointers for nullable fields
   - `<Resource>s` - list response wrapper with `Pagination` embed
   - `List<Resource>sParams` / `Filter<Resource>Params` - query parameter structs

3. **API methods** on `API` struct, delegating to generic helpers:
```go
func (api API) CreateCustomer(newCustomer *NewCustomer) (*Customer, error) {
    result := &Customer{}
    return result, api.create(customersEndpoint, newCustomer, result)
}
```

### IApi interface (`chartmogul.go`)

All public API methods are declared in the `IApi` interface. This enables mock generation for consumer tests. When adding a new endpoint, add the method signature here.

### Mocks (`mock_chartmogul/`)

Auto-generated via `./genmocks.sh` using `golang/mock`. Never edit manually.

### Retry and error handling

- Retries on: 429, 500, 502, 503, 504, and network errors (i/o timeout, `net.OpError`)
- Exponential backoff via `cenkalti/backoff/v3`
- Errors wrapped with `pkg/errors`
- `httpError` exposes `StatusCode()` and `Response()` via `HTTPError` interface
- `Errors` type (map) has helpers: `IsAlreadyExists()`, `IsInvoiceAndTransactionAlreadyExist()`, etc.

### Pagination

Two patterns embedded into list response structs:
- `Pagination` - cursor-based (`HasMore`, `Cursor`)
- `MetricsPagination` - page-based (`CurrentPage`, `TotalPages`)

## Testing

### Unit tests (root `*_test.go`)

- Standard `testing` package with `httptest.NewServer` for HTTP mocking
- Inline JSON constants as fixtures (e.g., `const oneInvoiceExample = ...`)
- Pattern: spin up mock server, call `SetURL(server.URL + "/v/%v")`, exercise the API method, assert fields
- Use `reflect.DeepEqual` or `go-test/deep` for struct comparison
- `spew.Dump` for debug output
- Live API tests gated behind `-cm` flag (skipped by default)

### Integration tests (`integration_tests/`)

- VCR-based via `dnaeon/go-vcr`: records/replays HTTP interactions as YAML cassettes in `integration_tests/fixtures/`
- Authorization headers filtered out of recordings
- Skipped in short mode (`-short` flag)

### Adding a new endpoint

1. Create `<resource>.go` with endpoint constants, types, and API methods
2. Create `<resource>_test.go` with inline JSON fixtures and httptest-based tests
3. Add method signatures to `IApi` interface in `chartmogul.go`
4. Run `./genmocks.sh` to regenerate mocks
5. Optionally add integration tests in `integration_tests/` with VCR fixture

## Code style

- Go conventions: `PascalCase` exports, `camelCase` locals, `UPPER_CASE` not used for constants
- Endpoint constants: `camelCase` + `Endpoint` suffix (e.g., `singleCustomerEndpoint`)
- JSON tags: `snake_case` matching the API, with `omitempty` on optional fields
- Pointer types (`*string`, `*bool`) in update structs to distinguish zero-value from absent
- Non-pointer receiver on API methods: `func (api API) CreateCustomer(...)`
- Linting: golangci-lint with config in `.golangci.yml`

## CI

GitHub Actions on push/PR to main. Go version matrix: 1.21, 1.22, 1.23, 1.24. Runs `make test` + golangci-lint + CodeClimate coverage.
