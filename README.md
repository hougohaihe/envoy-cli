# envoy-cli

> A CLI utility for managing and syncing environment variable sets across local and remote dev environments.

---

## Installation

```bash
pip install envoy-cli
```

Or install from source:

```bash
git clone https://github.com/yourname/envoy-cli.git && cd envoy-cli && pip install .
```

---

## Usage

```bash
# Push your local .env to a remote environment
envoy push --env staging --file .env

# Pull environment variables from a remote environment
envoy pull --env production --output .env.production

# List all available environment sets
envoy list

# Sync variables between two environments
envoy sync --from staging --to production
```

### Example

```bash
$ envoy push --env staging --file .env
✔ Loaded 12 variables from .env
✔ Connected to remote: staging
✔ Synced 12 variables successfully
```

---

## Configuration

`envoy-cli` looks for a config file at `~/.envoy/config.yaml` by default. You can override this with the `--config` flag.

```yaml
default_remote: staging
remotes:
  staging: https://envoy.example.com/staging
  production: https://envoy.example.com/production
```

---

## Contributing

Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

---

## License

[MIT](LICENSE)