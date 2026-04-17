# Lunar Tear

Private server research project for a certain discontinued mobile game.
Discord server: https://discord.gg/G3anrfcV

## How To Launch The Server

### Prerequisites

- Go 1.25+
- [goose](https://github.com/pressly/goose) migration tool
- Populated `server/assets/` directory

```bash
go install github.com/pressly/goose/v3/cmd/goose@latest
```

### Regenerate protobuf stubs

```bash
cd server
make proto
```

### Database

Player state is stored in a SQLite database. Run migrations before starting the server:

```bash
cd server
mkdir -p db
goose -dir migrations sqlite3 db/game.db up
```

To snapshot your current game state (e.g. at scene 401):

```bash
cp db/game.db db/snapshot_scene_401.db
```

To restore from a snapshot, just point the server at it:

```bash
./lunar-tear --db db/snapshot_scene_401.db
```

### Run

```bash
cd server
sudo go run ./cmd/lunar-tear \
  --host 10.0.2.2 \
  --http-port 8080
```

`sudo` is needed because gRPC binds to port 443 (privileged). On Linux you can use `setcap` instead:

```bash
go build -o lunar-tear ./cmd/lunar-tear
sudo setcap cap_net_bind_service=+ep ./lunar-tear
./lunar-tear --host 10.0.2.2 --http-port 8080
```

### Ports

| Protocol | Port | Notes                                                |
| -------- | ---- | ---------------------------------------------------- |
| gRPC     | 443  | hardcoded by the client, not configurable            |
| HTTP     | 8080 | Octo asset API + game web pages (`--http-port` flag) |

### Flags

| Flag          | Default      | Description                     |
| ------------- | ------------ | ------------------------------- |
| `--host`      | `127.0.0.1`  | hostname/IP given to the client |
| `--http-port` | `8080`       | HTTP/Octo server port           |
| `--db`        | `db/game.db` | SQLite database path            |

## ⚠️ Legal Disclaimer

**Lunar Tear** is a fan-made, non-commercial **preservation and research project** dedicated to keeping a certain discontinued mobile game playable for educational and archival purposes.

- This project is **not affiliated with**, **endorsed by**, or **approved by** the original publisher or any of its subsidiaries.
- All trademarks, copyrights, and intellectual property related to the original game and its associated franchises belong to their respective owners.
- All code in this repository is original work developed through clean-room reverse engineering for interoperability with the game client.
- No copyrighted game assets, binaries, or master data are distributed in this repository.

**Use at your own risk.** The author assumes no liability for any damages or legal consequences that may arise from using this software. By using or contributing to this project, you are solely responsible for ensuring your usage complies with all applicable laws in your jurisdiction.

This project is released under the [MIT License](LICENSE).

**If you are a rights holder with concerns regarding this project**, please contact me directly.
