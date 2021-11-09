package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
	"websocket"
)

var upgrader = websocket.Upgrader{}

var rpcList []RpcInfo
var primaryUrl = "wss://mainnet.infura.io/ws/v3/9bdd9b1d1270497795af3f522ad85091"
var secondaryUrl = "wss://ws-nd-986-369-125.p2pify.com/c669411d9bcc43aa0519602a30346446"
var tertiaryUrl = "wss://eth-mainnet.alchemyapi.io/v2/v1bo6tRKiraJ71BVGKmCtWVedAzzNTd6"

var currentId = 10000

var clientList []Client
var pendingResponses map[float64]*Client
var rpcMessageHandler chan MessageData

type Client struct {
	Id                   int
	Connection           *websocket.Conn
	SubscriptionHashList []string
	PendingResponses     map[float64]int
	ResponseHandler      chan MessageData
}

type MessageData struct {
	TypeMessage           int
	Message               []byte
	MessageAsJson         map[string]interface{}
	Name                  string
	Err                   error
	TimeReceivedInSeconds int64
}

type RpcInfo struct {
	Name       string
	Connection *websocket.Conn
	Responses  []MessageData
}

var urlList = [3]string{
	"wss://mainnet.infura.io/ws/v3/9bdd9b1d1270497795af3f522ad85091",
	"wss://ws-nd-986-369-125.p2pify.com/c669411d9bcc43aa0519602a30346446",
	"wss://eth-mainnet.alchemyapi.io/v2/v1bo6tRKiraJ71BVGKmCtWVedAzzNTd6",
}

func init() {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
}

func main() {
	clientList = make([]Client, 0)
	pendingResponses = make(map[float64]*Client)
	rpcMessageHandler = make(chan MessageData)

	initializeRpcList()
	initializeServiceLists()

	proxy := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		clientConnection, _ := upgrader.Upgrade(rw, req, nil)
		// defer clientConnection.Close()

		var client Client = createClient(clientConnection)
		clientList = append(clientList, client)

		go readMessageFromClientAndHandle(&client)

		for _, rpc := range rpcList {
			go readMessagesFromRPC(rpc)
		}

		go handleMessageReceivedFromRpc()

	})
	http.ListenAndServe(":8080", proxy)
}

func readMessageFromClientAndHandle(client *Client) {
	for {
		typeMessage, message, err := client.Connection.ReadMessage()
		fmt.Println("Message received from client")
		if err != nil {
			fmt.Println("Error reading message from client")
			return
		}
		// client.Connection.WriteMessage(typeMessage, message)

		handleMessageReceivedFromClient(client, typeMessage, message)
	}
}

func handleMessageReceivedFromRpc() {
	for {
		var message MessageData = <-rpcMessageHandler
		// In all the methods except when you are already subscribed, the message contains an Id
		// When you are already subscribed, you don't get an Id parameter, but a subscription hash instead.
		if _, ok := message.MessageAsJson["id"]; ok {
			handleMessageWithId(message)
		} else {

		}
	}
}

func handleMessageWithId(message MessageData) {
	var id = message.MessageAsJson["id"].(float64)
	if _, ok := pendingResponses[id]; !ok {
		fmt.Println("The message has already been sent")
		return
	}
	var client *Client = pendingResponses[id]
	var action int = getActionFromPendingResponses(client, id)
	switch action {
	case ACTION_SEND_FIRST_RESPONSE:
		client.Connection.WriteMessage(message.TypeMessage, message.Message)
		delete(pendingResponses, id)
		delete(client.PendingResponses, id)
	case ACTION_COMPARE_SERVICES:

	case ACTION_SEND_ONLY_ONCE:

	case ACTION_SEND_ONLY_TO_PRIMARY:

	}

}

func getActionFromPendingResponses(client *Client, id float64) int {
	return client.PendingResponses[id]
}

func initializeRpcList() {
	rpcList = make([]RpcInfo, 0)
	addRpcInfo(primaryUrl, "primary")
	addRpcInfo(secondaryUrl, "secondary")
	addRpcInfo(tertiaryUrl, "tertiary")
}

func addRpcInfo(url string, name string) {
	var server, _, _ = websocket.DefaultDialer.Dial(url, nil)
	var rpcInfo RpcInfo
	rpcInfo.Name = name
	rpcInfo.Connection = server
	rpcInfo.Responses = make([]MessageData, 0)
	rpcList = append(rpcList, rpcInfo)
}

func createClient(connection *websocket.Conn) Client {
	var client Client
	client.Id = currentId
	currentId = currentId + 1
	client.Connection = connection
	client.SubscriptionHashList = make([]string, 0)
	client.PendingResponses = make(map[float64]int)
	client.ResponseHandler = make(chan MessageData)
	return client
}

func handleMessageReceivedFromClient(client *Client, typeMessage int, message []byte) {
	var id, action = getActionAndIdFromClientMessage(message)
	pendingResponses[id] = client
	client.PendingResponses[id] = action
	switch action {
	case ACTION_SERVICE_DOESNT_EXISTS:
		// Send to the corresponding channel the response
		var result MessageData = createErrorResponse(id, "Service doesn't exists")
		client.ResponseHandler <- result
	case ACTION_SEND_FIRST_RESPONSE:
		fmt.Println("Message to send after first response")
		for _, rpc := range rpcList {
			rpc.Connection.WriteMessage(typeMessage, message)
		}
	case ACTION_COMPARE_SERVICES:
		for _, rpc := range rpcList {
			rpc.Connection.WriteMessage(typeMessage, message)
		}
	case ACTION_SEND_ONLY_ONCE:
		rpcList[0].Connection.WriteMessage(typeMessage, message)
	case ACTION_SEND_ONLY_TO_PRIMARY:
		rpcList[0].Connection.WriteMessage(typeMessage, message)
	case ACTION_DEPRECATED:
		var result MessageData = createErrorResponse(id, "Service deprecated")
		client.ResponseHandler <- result
	}
}

func createErrorResponse(id float64, message string) MessageData {
	var jsonMessage string = fmt.Sprintf("{\"id\":%f,\"result\":%s}", id, message)
	var messageData MessageData
	messageData.TypeMessage = websocket.TextMessage
	messageData.Message = []byte(jsonMessage)
	messageData.Name = "primary"
	messageData.Err = nil
	messageData.TimeReceivedInSeconds = 0
	return messageData
}

func getActionAndIdFromClientMessage(message []byte) (float64, int) {
	var result map[string]interface{}
	json.Unmarshal(message, &result)

	return result["id"].(float64), getAction(result["method"].(string))
}

func readMessagesFromRPC(rpcInfo RpcInfo) {
	for {
		typeMessage, message, err := rpcInfo.Connection.ReadMessage()
		log.Printf("Message received from %s\n", rpcInfo.Name)
		var result map[string]interface{}
		json.Unmarshal(message, &result)
		// log.Printf("%s\n\n", message)
		var res MessageData
		res.TypeMessage = typeMessage
		res.Message = message
		res.MessageAsJson = result
		res.Name = rpcInfo.Name
		res.Err = err
		res.TimeReceivedInSeconds = time.Now().Unix()
		rpcMessageHandler <- res
	}
}
