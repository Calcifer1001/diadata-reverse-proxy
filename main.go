package main

import (
	"fmt"
	"net/url"
	"testing"
	"websocket"
)

var u = url.URL{Scheme: "ws", Host: ":5555", Path: "/reverse"}
var infuraUrl = "wss://mainnet.infura.io/ws/v3/9bdd9b1d1270497795af3f522ad85091"
var upgrader = websocket.Upgrader{}

func main() {
	var urlWs, _ = url.Parse(infuraUrl)
	fmt.Println(urlWs.String())
}
