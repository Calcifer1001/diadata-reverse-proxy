# RPC REVERSE PROXY

This is a reverse proxy design to get data from different rpc servers from the ethereum blockchain

## How to run

To run the proxy:

```
cd wsServer/
go run main.go decisionFlow.go
```

After this you need to send a json message as the rpc server expects. Most of the examples are at https://eth.wiki/json-rpc/API at the JSON-RPC methods section

One way of sending the messages is using [websocat](https://github.com/vi/websocat) on linux.

For example if you run the proxy and the websocket communication locally, and want to get the method web3_clientVersion:

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

## PROXY RESPONSES
The proxy divides each service into several categories:
 - Doesn't exist
 - Send first response
 - Compare Services
 - Send only to primary
 - Deprecated

### Doesn't exist
If you send a method that is not in neither of the other categories, you will get a message informing that the method doesn't exist

### Send first response
Your request will be sent to all of the rpc servers and as soon as one of them returns a response, it will be sent to the client. The following messages received concerning this request, will be discarded

### Compare services
When the proxy receives a response, it will store it. As soon as it gets a second response, if boht messages are the same, one of them will be sent to the client. Otherwise, it will wait until the third rpc returns a response and again try to compare the messages and send on a match. If there is no match, you will receive a message with the message "No match".

### Send only to primary
These services will be sent only to the primary server

### Deprecated
Methods that are not longer supported. You will get a different error message than the "doesn't exist" case.

## Subscription
This case is completely special. When you get any rpc response from the method "eth_subscribe", the proxy will send you the first subscription hash received. The proxy will store all three subscription hashes, but will always send you any of the eth_subcription responses with the first subscription hash received. 

The comparison of the received messages will be done concerning this 3 hashes and you might not get a transaction's information that happened, only in the case that the proxy doesn't receive two equal responses. This should be completely rare. 

This text might not be entirely clear, but we did our best. Feel free to contact us in case you don't understand it. :)

### Testing

There are two test files in wsServer/tests witch allows to verify the websocket connection with the servers for the supported RPC methods, one with the first response methods and other with the comparison of responses. 
To run the tests simply go to the wsServer/tests directory and run in a terminal the command "go test" 
