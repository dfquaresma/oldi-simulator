# hedge

Lightweight latency/distribution simulator for testing and performance analysis.

## Quick start

Initialize modules (first run):
```bash
go mod init hedge
echo "replace github.com/dfquaresma/hedge => ./" >> go.mod
go mod tidy
```

Build:
```bash
go build ./...
```

To run the Latency Model:
```bash
cd /latency_model/
go run main.go
```
