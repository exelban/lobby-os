const list = {
  "unifi": {
    "name": "UniFi",
    "color": "#0061ff",
    "icon": "/presets/unifi.png"
  },
  "proxmox": {
    "name": "Proxmox",
    "color": "#d6d6d6",
    "icon": "/presets/proxmox.png"
  },
  "hass": {
    "name": "Home Assistant",
    "icon": "/presets/hass.png"
  },
  "cloudflare": {
    "name": "Cloudflare",
    "icon": "/presets/cloudflare.png"
  },
  "flatnotes": {
    "name": "FlatNotes",
    "color": "#fec700",
    "icon": "/presets/flatnotes.png"
  },
  "gitea": {
    "name": "Gitea",
    "color": "#b1dd8c",
    "icon": "/presets/gitea.png"
  },
  "grafana": {
    "name": "Grafana",
    "icon": "/presets/grafana.png"
  },
  "homebridge": {
    "name": "Homebridge",
    "color": "#232323",
    "icon": "/presets/homebridge.png"
  },
  "influx": {
    "name": "InfluxDB",
    "color": "#0056d6",
    "icon": "/presets/influxdb.png"
  },
  "jackett": {
    "name": "Jackett",
    "icon": "/presets/jackett.png"
  },
  "jam": {
    "name": "JAM",
    "color": "#444444",
    "icon": "/presets/jam.png"
  },
  "jellyfin": {
    "name": "Jellyfin",
    "icon": "/presets/jellyfin.png"
  },
  "kuma": {
    "name": "Kuma",
    "icon": "/presets/kuma.png"
  },
  "lidarr": {
    "name": "Lidarr",
    "icon": "/presets/lidarr.png"
  },
  "myspeed": {
    "name": "MySpeed",
    "color": "#9FADE9",
    "icon": "/presets/myspeed.png"
  },
  "pihole": {
    "name": "Pi-hole",
    "color": "#444444",
    "icon": "/presets/pihole.png"
  },
  "plex": {
    "name": "Plex",
    "icon": "/presets/plex.png"
  },
  "portainer": {
    "name": "Portainer",
    "icon": "/presets/portainer.png"
  },
  "prowlarr": {
    "name": "Prowlarr",
    "icon": "/presets/prowlarr.png"
  },
  "radarr": {
    "name": "Radarr",
    "color": "#BD0E1F",
    "icon": "/presets/radarr.png"
  },
  "scrypted": {
    "name": "Scrypted",
    "color": "#7b219f",
    "icon": "/presets/scrypted.png"
  },
  "sonarr": {
    "name": "Sonarr",
    "color": "#DA93FD",
    "icon": "/presets/sonarr.png"
  },
  "synology": {
    "name": "Synology",
    "icon": "/presets/synology.png"
  },
  "transmission": {
    "name": "Transmission",
    "icon": "/presets/transmission.png"
  },
  "zigbee2mqtt": {
    "name": "Zigbee2MQTT",
    "icon": "/presets/zigbee2mqtt.png"
  },
  "openWebUI": {
    "name": "Open Web UI",
    "color": "#000000",
    "icon": "/presets/openwebui.png"
  },
  "llama": {
    "name": "llama",
    "color": "#FFFFFF",
    "icon": "/presets/llama.png"
  },
  "truenas": {
    "name": "TrueNAS",
    "icon": "/presets/truenas.png"
  },
  "GoogleSearchConsole": {
    "name": "Google Search Console",
    "icon": "/presets/googlesearchconsole.png"
  },
  "stable-diffusion": {
    "name": "Stable Diffusion",
    "icon": "/presets/stablediffusion.png"
  },
  "traefik": {
    "name": "Traefik",
    "color": "#ffffff",
    "icon": "/presets/traefik.png"
  },
  "nut": {
    "name": "NUT",
    "color": "#0067cd",
    "icon": "/presets/nut.png"
  },
  "docmost": {
    "name": "Docmost",
    "color": "#000000",
    "icon": "/presets/docmost.png"
  },
  "filebrowser": {
    "name": "File Browser",
    "icon": "/presets/filebrowser.png"
  },
  "netdetect": {
    "name": "NetDetect",
    "color": "#ffffff",
    "icon": "/presets/netdetect.png"
  },
  "paperless": {
    "name": "Paperless",
    "color": "#2b8a38",
    "icon": "/presets/paperless.png"
  },
  "nutshell": {
    "name": "NutShell",
    "color": "#ffffff",
    "icon": "/presets/nutshell.png"
  },
  "emqx": {
    "name": "EMQX",
    "color": "#ffffff",
    "icon": "/presets/emqx.png"
  },
  "slzb": {
    "name": "SLZB",
    "color": "#212529",
    "icon": "/presets/slzb.png"
  },
  "bitwarden": {
    "name": "Bitwarden",
    "color": "#175DDC",
    "icon": "/presets/bitwarden.png"
  },
  "drawio": {
    "name": "Draw.io",
    "color": "#ffffff",
    "icon": "/presets/drawio.png"
  },
  "frigate": {
    "name": "Frigate",
    "color": "#ffffff",
    "icon": "/presets/frigate.png"
  },
  "netbootxyz": {
    "name": "Netboot.xyz",
    "color": "#ffffff",
    "icon": "/presets/netbootxyz.png"
  },
  "adguard-home": {
    "name": "AdGuard Home",
    "color": "#ffffff",
    "icon": "/presets/adguard-home.png"
  },
  "umami": {
    "name": "Umami",
    "color": "#ffffff",
    "icon": "/presets/umami.png"
  },
}

export default Object.keys(list).sort().reduce((acc, key) => ({ ...acc, [key]: list[key] }), {})