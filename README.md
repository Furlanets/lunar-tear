# Lunar Tear

Private server research project for a certain discontinued mobile game.
Discord server: https://discord.gg/MZAf5aVkJG

## How To Launch The Server

### Prerequisites

- Go 1.25+
- [goose](https://github.com/pressly/goose) migration tool
- Populated `server/assets/` directory

```bash
go install github.com/pressly/goose/v3/cmd/goose@latest
```

### Quick Start (Wizard)

The interactive wizard walks you through setup with a few simple questions — no flags or networking knowledge needed. It auto-detects the right IP address for your emulator or phone and launches all services.

```bash
cd server
go run ./cmd/wizard
```

Your choices are saved so next time you just press Enter to relaunch with the same settings.

#### Custom Ports

By default the wizard uses ports 8003 (gRPC), 8080 (CDN), and 3000 (auth). Override any of them with flags:

```bash
go run ./cmd/wizard --grpc-port 9003 --cdn-port 9080
```

| Flag          | Default | Description      |
| ------------- | ------- | ---------------- |
| `--grpc-port` | `8003`  | gRPC server port |
| `--cdn-port`  | `8080`  | CDN server port  |
| `--auth-port` | `3000`  | Auth server port |

Custom ports are saved to `.wizard.json` alongside your other settings. On the next run the saved ports are reused automatically — no need to pass the flags again. If you later pass different port flags, the wizard warns you that the ports changed and asks for confirmation before continuing.

### Regenerate protobuf stubs

```bash
cd server
make proto
```

### Database

Player state is stored in a SQLite database. Run migrations before starting the server:

```bash
cd server
make migrate
```

Or manually:

```bash
cd server
mkdir -p db
goose -dir migrations sqlite3 db/game.db up
```

### Importing a Snapshot

To import a JSON snapshot into the database, use the import tool. The `--uuid` flag must match the UUID your game client sends during authentication:

```bash
cd server
make import SNAPSHOT=snapshots/scene_1.json UUID=<your-client-uuid>
```

Or directly:

```bash
go run ./cmd/import-snapshot \
  --snapshot snapshots/scene_1.json \
  --uuid <your-client-uuid> \
  --db db/game.db
```

| Flag         | Default      | Description                                   |
| ------------ | ------------ | --------------------------------------------- |
| `--snapshot` | *(required)* | Path to JSON snapshot file                     |
| `--uuid`     | *(required)* | UUID to assign (must match the client's UUID)  |
| `--db`       | `db/game.db` | SQLite database path                           |

### Run

```bash
cd server
go run ./cmd/lunar-tear \
  --host 10.0.2.2 \
  --http-port 8080 \
  --grpc-port 8003
```

The default gRPC port is 443, which requires `sudo` (privileged port). Use `--grpc-port` with a high port to avoid this. If you do need port 443, either use `sudo` or grant the binary the capability on Linux:

```bash
go build -o lunar-tear ./cmd/lunar-tear
sudo setcap cap_net_bind_service=+ep ./lunar-tear
./lunar-tear --host 10.0.2.2 --http-port 8080
```

The CDN can run on a completely separate machine — just set `--octo-url` on the game server and `--public-addr` on the CDN to the externally-reachable address.

### Run With Dev Runner

Instead of manually building binaries, use the dev runner to build and launch `lunar-tear` with one command. No Docker required — works on macOS, Linux, and Windows.

```bash
cd server
make dev
```

Or directly:

```bash
cd server
go run ./cmd/dev
```

Output is prefixed with a colored label (`[grpc]`). Press Ctrl+C to shut everything down.

The dev runner automatically builds the required binaries into `bin/` before launching. This means executables have stable file paths, so **Windows Firewall only prompts once** — subsequent runs reuse the same allowed executable. The wizard performs the same build step transparently.

Override defaults with namespaced flags:

```bash
go run ./cmd/dev --grpc.listen 0.0.0.0:9000 --grpc.public-addr 10.0.2.2:9000 --grpc.octo-url http://192.168.1.50:8080
```

Or via `make`:

```bash
make dev ARGS="--grpc.listen 0.0.0.0:9000 --grpc.public-addr 10.0.2.2:9000"
```

| Flag                  | Default            | Description                              |
| --------------------- | ------------------ | ---------------------------------------- |
| `--grpc.listen`       | `0.0.0.0:8003`     | lunar-tear gRPC listen address           |
| `--grpc.public-addr`  | `10.0.2.2:8003`    | lunar-tear externally-reachable addr     |
| `--grpc.octo-url`     | `http://10.0.2.2:8080` | Octo CDN base URL passed to lunar-tear |
| `--grpc.auth-url`     | `http://localhost:3000` | auth server base URL passed to lunar-tear |
| `--no-color`          | `false`            | disable colored output                   |

### Ports

| Protocol | Port | Notes                                                       |
| -------- | ---- | ----------------------------------------------------------- |
| gRPC     | 443  | default; configurable with `--grpc-port` (requires patched client) |
| HTTP     | 8080 | Octo asset API + game web pages (`--http-port` flag)        |

### Flags

| Flag          | Default      | Description                                          |
| ------------- | ------------ | ---------------------------------------------------- |
| `--host`      | `127.0.0.1`  | hostname/IP given to the client                      |
| `--http-port` | `8080`       | HTTP/Octo server port                                |
| `--grpc-port` | `443`        | gRPC server port (client must be patched to match)   |
| `--db`        | `db/game.db` | SQLite database path                                 |

### Docker

Migrations run automatically on container start.

```bash
cd server
docker compose up -d
```

The `db/` directory is mounted as a volume so the database persists across restarts. Make sure `assets/` is populated before starting.

### Makefile Targets

All targets run from the `server/` directory.

| Target         | Description                                             |
| -------------- | ------------------------------------------------------- |
| `make proto`   | Regenerate protobuf stubs                               |
| `make build`   | Build the game server binary                            |
| `make build-dev` | Build the dev runner binary to `bin/`                  |
| `make build-all` | Build available binaries to `bin/`                     |
| `make build-import` | Build the import-snapshot tool                     |
| `make clean`   | Remove the `bin/` directory                              |
| `make dev`     | Build and run `lunar-tear` via dev runner                |
| `make migrate` | Run goose migrations on `db/game.db`                    |
| `make import`  | Import a snapshot (`SNAPSHOT=... UUID=...` required)     |

## ⚠️ Legal Disclaimer

**Lunar Tear** is a fan-made, non-commercial **preservation and research project** dedicated to keeping a certain discontinued mobile game playable for educational and archival purposes.

- This project is **not affiliated with**, **endorsed by**, or **approved by** the original publisher or any of its subsidiaries.
- All trademarks, copyrights, and intellectual property related to the original game and its associated franchises belong to their respective owners.
- All code in this repository is original work developed through clean-room reverse engineering for interoperability with the game client.
- No copyrighted game assets, binaries, or master data are distributed in this repository.

**Use at your own risk.** The author assumes no liability for any damages or legal consequences that may arise from using this software. By using or contributing to this project, you are solely responsible for ensuring your usage complies with all applicable laws in your jurisdiction.

This project is released under the [MIT License](LICENSE).

**If you are a rights holder with concerns regarding this project**, please contact me directly.
