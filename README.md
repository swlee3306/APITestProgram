# APITestProgram

Go HTTP testing experiments for local and offline-network development.

## Run and verify

```sh
go build -mod=readonly ./...
go test -mod=readonly ./chain/...
go run .
```

The default program invokes `internal/apiRequest.ApiReqRun`, an interactive request flow. It is not a flag parser implementing all the advanced packages. Only send requests to services you are authorized to test.

## Implementation map

- `internal/apiRequest/`: default interactive request path.
- `chain/`: chain manager, variables and conditional helpers; has unit tests.
- `http/`, `security/`, `logging/`: extended HTTP/TLS/header/logging implementations.
- `cli/`, `history/`, `config/`, `ui/`: supporting components.

The extended packages are not all wired into the default entry point. Their presence does not mean every feature is available from `go run .`. No security certification, complete response validation or production reliability is claimed.

For defect reports, use a local fake server and synthetic payload. Do not include authorization headers, private endpoints, request history or client certificates.

[Project map](https://github.com/swlee3306#selected-work)
