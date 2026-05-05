# envdiff

> Compare `.env` files across environments and surface missing or mismatched keys.

---

## Installation

```bash
go install github.com/yourusername/envdiff@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/envdiff.git
cd envdiff && go build -o envdiff .
```

---

## Usage

```bash
envdiff [flags] <base> <target>
```

Compare a local `.env` against a production template:

```bash
envdiff .env.example .env.production
```

**Example output:**

```
MISSING in .env.production:
  - DATABASE_URL
  - REDIS_HOST

MISMATCHED values:
  ~ LOG_LEVEL: "debug" → "info"

✔ 12 keys match
```

### Flags

| Flag | Description |
|------|-------------|
| `--keys-only` | Compare key names only, ignore values |
| `--quiet` | Exit with non-zero status on diff, no output |
| `--format json` | Output results as JSON |

---

## Contributing

Pull requests are welcome. Please open an issue first to discuss any significant changes.

---

## License

[MIT](LICENSE) © 2024 yourusername