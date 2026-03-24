# JAD

[![JAD](https://serhiy.s3.eu-central-1.amazonaws.com/Github_repo/JAD/cover.png)](https://github.com/exelban/JAD)

**Just Another Dashboard** — a minimal, self-hosted start page for your homelab.

## Features

- Organize links into custom groups with drag-and-drop reordering
- 50+ built-in presets for popular homelab apps (Proxmox, Home Assistant, Jellyfin, Grafana, Pi-hole, and more)
- Automatic favicon extraction from linked URLs
- Custom colors for each link card
- Dark / Light / System theme
- Import and export links as JSON
- Responsive design for desktop and mobile

## Installation

### Docker

```bash
docker run -d -p 8080:8080 -v jad_data:/srv/data exelban/jad:latest
```

### Docker Compose

```yaml
services:
  jad:
    image: exelban/jad:latest
    restart: always
    ports:
      - "8080:8080"
    volumes:
      - jad_data:/srv/data
    environment:
      - PUID=1000
      - PGID=1000

volumes:
  jad_data:
```

Open `http://localhost:8080` in your browser.

## Configuration

| Variable | Default | Description |
|---|---|---|
| `DATA_PATH` | `/srv/data` | Directory where `links.json` is stored |

You can also pass `--data-path` as a command-line flag.

## License
[MIT License](https://github.com/exelban/JAD/blob/master/LICENSE)
