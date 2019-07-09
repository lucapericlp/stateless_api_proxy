# Stateless API Proxy

This is a Go implementation of the [magic-github-proxy](https://github.com/theacodes/magic-github-proxy) to become familiar with Go and the ecosystem.

## Explanation and Usage
TODO: More detail and explanation

Make sure environment variables are set:
```bash
export TEST_SERVER_ADDR=":8080"
export TEST_KEY_FILE="$HOME/path/to/server.key"
export TEST_CERT_FILE="$HOME/path/to/server.crt"
export MAGICTOKEN_PUBLIC_KEY="$HOME/path/to/mykey.pub"
export MAGICTOKEN_PRIVATE_KEY="$HOME/path/to/priv.pem"
```

Run proxy server using:

```bash
go run main.go
```

Example for token creation using HTTPie where github\_token is the OAuth2 Github token:
```bash
echo '{"github_token":"abc123","scopes":["GET /user","GET /repos"]}' | http POST https://localhost:8080/create --json --verify no
#verify flag disabled due to self signed cert for dev
```

Example for Github API request where $JWT is the result of the previous call:
```bash
http POST https://localhost:8080/api/repos/something --json --verify no --auth-type=jwt --auth=$JWT
```
