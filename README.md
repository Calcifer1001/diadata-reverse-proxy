# RPC REVERSE PROXY

This is a reverse proxy design to get data from different rpc servers from the ethereum blockchain

## How to run

To run the proxy:

```
cd wsServer/
go run main.go decisionFlow.go
```

After this you need to send a json message as the rpc server expects. Most of the examples are at https://eth.wiki/json-rpc/API at the JSON-RPC methods section

One way of sending the messages is using [websocat]https://github.com/vi/websocat on linux.

If you run the proxy and the communication locally, and want to get the web3_clientVersion:

```
websocat --linemode-strip-newlines localhost:8080
{"id":1,"jsonrpc":"2.0","method":"web3_clientVersion","params":[]}
```

And as an response example:
```
{"jsonrpc":"2.0","id":1,"result":"Geth/v1.10.12-stable-6c4dc6c3/linux-amd64/gol.17.3"}
```

## RPC SERVER
In the main.go file, you can see three variables: primaryServer, secondaryServer and tertiaryServer. These are the websocket url where the proxy is heading to get the information. It is currently getting information from the kovan network. Feel free to modify this to go to mainnet.

