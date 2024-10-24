# scdlbot

Discord bot to download SoundCloud tracks as Discord embeds. Requires a running [cookie relay server](https://github.com/7x11x13/cookie-relay)

# Setup

## Example `docker-compose.yaml`

```yaml
name: scdlbot

services:
  bot:
    image: ghcr.io/7x11x13/scdlbot:latest
    restart: always
    environment:
      - BOT_TOKEN=xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
      - COOKIE_RELAY_URL=https://cookie-relay.xxxxxxxxxx.ts.net
      - COOKIE_RELAY_API_KEY=xxxxxxxxxxxxxxxxxx
      - SOUNDCLOUD_USER_ID=123456789
```
