realtime-sign-go
==

Golang HTTP API and WebSocket server for real-time sign updates

Environment variables
==

* `RSIGN_API_KEY`: The authentication secret that should be accepted for the `X-API-Key` header
* `RSIGN_ADDR`: The address to bind the server to (default `:3000`)

Building
==

    dep ensure
    go build main.go
