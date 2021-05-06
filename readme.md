
Authorization proxy app for minetest
=================

![](https://github.com/minetest-auth-proxy/auth_proxy_app/workflows/docker/badge.svg)


To be used with https://github.com/minetest-auth-proxy/auth_proxy_mod

# Overview

Lets third-party apps query username and password of ingame players

# Docker

Run with docker:
```sh
docker run --rm -it -p 8080:8080 minetestauth/minetest-auth-proxy
```

# Implementations

## Mediawiki

Thanks to **gpcf** (linuxforks) for his partial Mediawiki implementation!
See the folder: `mediawiki`

## cURL

```bash
curl -X POST -H 'Content-Type: application/json' -i 'http://127.0.0.1:8080/api/login' --data '{"username":"test","password":"enter"}'
```

Returns

On success:
```json
{"success": true, "message": null}
```

On failure:
```json
{"success": false, "message": "Banned!"}
```

Or just:
```json
{"success": false, "message": ""}
```

# Building / Running

A `Dockerfile` is included for container usage.

Otherwise just `go run .` (needs a working go installation)
