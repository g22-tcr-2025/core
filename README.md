# Clash Royale Terminal

A terminal-based multiplayer strategy game inspired by Clash Royale, implemented in Go. Players connect to a server, log in, and battle each other using troops and towers, with real-time updates and colorful terminal UI.

## Features

- Multiplayer server-client architecture over TCP
- User authentication and persistent metadata
- Real-time match updates and notifications
- Troop and tower management with leveling and EXP
- Colorful, interactive terminal UI

## Project Structure

```
assets/           # Game data: users, standard troops/towers, player metadata
cmd/
  client/         # Client application entrypoint
  server/         # Server application entrypoint
internal/
  config/         # Game and protocol constants
  data/           # User and metadata management
  game/           # Core game logic (engine, player, troop, tower, etc.)
  logic/          # Matchmaking logic
  network/        # TCP networking and protocol
  ui/             # Terminal UI rendering and color utilities
pkg/
  utils/          # Utility functions (e.g., rune width for UI)
```

## Getting Started

### Prerequisites

- Go 1.24+
- Terminal that supports ANSI escape codes

### Install Dependencies

```sh
go mod tidy
```

### Running the Server

```sh
go run ./cmd/server
```

### Running the Client

Open a new terminal for each player:

```sh
go run ./cmd/client
```

### Default Users

See [`assets/users.json`](assets/users.json):

- Username: `player1`, Password: `1`
- Username: `player2`, Password: `2`

## Gameplay

- Log in with your username and password.
- Wait for matchmaking.
- Use commands to attack:  
  ```
  <troop_index> <tower_index>
  ```
  Example: `0 2` (use troop 0 to attack tower 2)
- Watch for notifications and match updates in real time.

## Customization

- Add new users in [`assets/users.json`](assets/users.json).
- Edit troop/tower stats in [`assets/standard.json`](assets/standard.json).
- Player progress is saved in `assets/metadata/<username>/metadata.json`.