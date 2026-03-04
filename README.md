# screenpng

An HTTP server for taking full-page screenshots of web pages in PNG format via Selenium WebDriver.

Built on top of the [pog7x/ssfactory](https://github.com/pog7x/ssfactory) library. Supports Firefox and Chrome in headless mode.

## Stack

- Go 1.19, [gorilla/mux](https://github.com/gorilla/mux), [cobra](https://github.com/spf13/cobra) + [viper](https://github.com/spf13/viper), [zap](https://github.com/uber-go/zap)
- Docker image: [pog7x/gobasebrowser](https://hub.docker.com/repository/docker/pog7x/gobasebrowser)

## Running (Docker)

```bash
docker build -t screenpng .
docker run -v $(pwd):/screenpng -p 8099:8099 -it screenpng
```

## API

### `POST /screenshot`

Accepts a list of URLs and saves screenshots as PNG files.

```json
{
  "items": [
    { "url": "https://example.com", "name": "example.png" },
    { "url": "https://go.dev/dl/", "name": "godev.png" }
  ]
}
```

## Configuration

The configuration file is passed via the `-c` flag (example: `configs/.screenpng-config.dev.yml`).

| Parameter | Description | Default (dev) |
|---|---|---|
| `use_browser` | Browser (`firefox` / `chrome`) | `firefox` |
| `webdriver_port` | WebDriver port | `8089` |
| `server_listen_addr` | Server address | `0.0.0.0:8099` |
| `server_read_timeout` | Read timeout | `3s` |
| `server_write_timeout` | Write timeout | `3s` |

## License

MIT
