# consul-drift-check

> CLI tool that detects configuration drift between Consul KV namespaces across datacenters

## Installation

```bash
go install github.com/yourusername/consul-drift-check@latest
```

Or download a pre-built binary from the [releases page](https://github.com/yourusername/consul-drift-check/releases).

## Usage

Compare KV namespaces between two Consul datacenters:

```bash
consul-drift-check \
  --source http://consul-dc1.example.com:8500 \
  --target http://consul-dc2.example.com:8500 \
  --prefix services/config
```

**Example output:**

```
[MISSING]  services/config/app/timeout      (not found in target)
[DRIFT]    services/config/app/max-retries  dc1="5"  dc2="3"
[OK]       services/config/app/log-level

Summary: 1 missing, 1 drifted, 1 in sync
```

### Flags

| Flag | Description | Default |
|------|-------------|---------|
| `--source` | Source Consul address | `http://localhost:8500` |
| `--target` | Target Consul address | *(required)* |
| `--prefix` | KV prefix to compare | `/` |
| `--token` | Consul ACL token | *(optional)* |
| `--output` | Output format: `text`, `json` | `text` |

## Requirements

- Go 1.21+
- Access to both Consul datacenters (HTTP API)

## License

MIT © [yourusername](https://github.com/yourusername)