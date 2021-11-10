package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
	"websocket"
)

// Acomodar ids
// Hacer comparaciones servicios con id
// Hacer comparaciones servicios sin id

var upgrader = websocket.Upgrader{}

var rpcList []RpcInfo
var primaryUrl = "wss://mainnet.infura.io/ws/v3/9bdd9b1d1270497795af3f522ad85091"
var secondaryUrl = "wss://ws-nd-986-369-125.p2pify.com/c669411d9bcc43aa0519602a30346446"
var tertiaryUrl = "wss://eth-mainnet.alchemyapi.io/v2/v1bo6tRKiraJ71BVGKmCtWVedAzzNTd6"

var currentId = 10000

var clientList []*Client
var pendingResponses map[float64]*Client
var subscribePendingResponses map[float64]*Client
var rpcMessageHandler chan MessageData

type Client struct {
	Id                      int
	Connection              *websocket.Conn
	SubscriptionHashList    []string
	PendingResponses        map[float64]int
	PendingSubscriptionId   float64
	MessagesAwaitingCompare []MessageData
	ResponseHandler         chan MessageData
}

type MessageData struct {
	TypeMessage           int
	Message               []byte
	MessageAsJson         map[string]interface{}
	RpcName               string
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
	clientList = make([]*Client, 0)
	pendingResponses = make(map[float64]*Client)
	subscribePendingResponses = make(map[float64]*Client)
	rpcMessageHandler = make(chan MessageData)

	initializeRpcList()
	initializeServiceLists()

	proxy := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		clientConnection, _ := upgrader.Upgrade(rw, req, nil)
		// defer clientConnection.Close()

		var client Client = createClient(clientConnection)
		clientList = append(clientList, &client)

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
		if id, ok := message.MessageAsJson["id"]; ok {
			if _, ok2 := subscribePendingResponses[id.(float64)]; ok2 {
				handleSubscriptionMessage(message)
			} else {
				handleMessageWithId(message)
			}
		} else {
			handleSubscribedMessage(message)
		}
	}
}

func handleSubscriptionMessage(message MessageData) {
	var id = message.MessageAsJson["id"].(float64)
	var client *Client = pendingResponses[id]
	var hash string = message.MessageAsJson["result"].(string)
	client.SubscriptionHashList = append(client.SubscriptionHashList, hash)
	pendingResponses[id] = client
	if len(client.SubscriptionHashList) == 1 {
		client.Connection.WriteMessage(message.TypeMessage, message.Message)
	}
	if len(client.SubscriptionHashList) >= 3 {
		delete(pendingResponses, id)
		delete(subscribePendingResponses, id)
		delete(client.PendingResponses, id)
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
	// to avoid tons of repeated code this flag is created and handled right after the switch

	switch action {
	case ACTION_SERVICE_DOESNT_EXISTS,
		ACTION_SEND_FIRST_RESPONSE,
		ACTION_SEND_ONLY_ONCE,
		ACTION_SEND_ONLY_TO_PRIMARY,
		ACTION_DEPRECATED:

		client.Connection.WriteMessage(message.TypeMessage, message.Message)
		delete(pendingResponses, id)
		delete(client.PendingResponses, id)

	case ACTION_COMPARE_SERVICES:
		index, send := compareMessageWithPendingMessages(message, client)
		if send {
			client.Connection.WriteMessage(message.TypeMessage, message.Message)
			delete(pendingResponses, id)
			delete(client.PendingResponses, id)
			// Removes the message from the slice
			client.MessagesAwaitingCompare = append(client.MessagesAwaitingCompare[:index], client.MessagesAwaitingCompare[index+1:]...)
		} else {
			client.MessagesAwaitingCompare = append(client.MessagesAwaitingCompare, message)
		}

	}

}

func handleSubscribedMessage(message MessageData) {
	var hash string = message.MessageAsJson["params"].(map[string]interface{})["subscription"].(string)
	for _, client := range clientList {
		for _, clientHash := range client.SubscriptionHashList {
			if clientHash == hash {
				index, send := shouldSendMessage(message, client)
				if send {
					setHash(&message, client.SubscriptionHashList[0])
					client.Connection.WriteMessage(message.TypeMessage, message.Message)
					// Removes message from expecting messages
					client.MessagesAwaitingCompare = append(client.MessagesAwaitingCompare[:index], client.MessagesAwaitingCompare[index+1:]...)
				} else {
					client.MessagesAwaitingCompare = append(client.MessagesAwaitingCompare, message)
				}
			}
		}
	}
}

func setHash(message *MessageData, hash string) {
	message.MessageAsJson["params"].(map[string]interface{})["subscription"] = hash
	newMessage, err := json.Marshal(message.MessageAsJson)
	if err != nil {
		fmt.Println("Error converting json to string")
	}
	message.Message = newMessage

}

func shouldSendMessage(message MessageData, client *Client) (int, bool) {
	for index, previousMessage := range client.MessagesAwaitingCompare {
		if equals(previousMessage, message) {
			return index, true
		}
	}
	return -1, false
}

func handleSubscribe(client *Client, message []byte) {
	var result map[string]interface{}
	json.Unmarshal(message, &result)
	if result["method"] == "eth_subscribe" {
		var id float64 = result["id"].(float64)
		subscribePendingResponses[id] = client
		client.PendingSubscriptionId = id
	}
}

func compareMessageWithPendingMessages(message MessageData, client *Client) (int, bool) {
	var newMessageResult map[string]interface{} = extractResult(message.MessageAsJson)
	for index, pendingMessage := range client.MessagesAwaitingCompare {
		var pendingMessageResult map[string]interface{} = extractResult(pendingMessage.MessageAsJson)
		if strings.EqualFold(fmt.Sprint(newMessageResult), fmt.Sprint(pendingMessageResult)) {
			return index, true
		}
	}
	return -1, false
}

func equals(message1, message2 MessageData) bool {

	var result1 = message1.MessageAsJson["params"].(map[string]interface{})["result"]
	var result2 = message1.MessageAsJson["params"].(map[string]interface{})["result"]
	// return result1 == result2
	// Cannot use deep equal because some server return message with uppercase and some with lowercase
	return strings.EqualFold(fmt.Sprint(result1), fmt.Sprint(result2))
}

func extractResult(message map[string]interface{}) map[string]interface{} {
	return message["params"].(map[string]interface{})["result"].(map[string]interface{})
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
	client.MessagesAwaitingCompare = make([]MessageData, 0)
	client.ResponseHandler = make(chan MessageData)
	return client
}

func handleMessageReceivedFromClient(client *Client, typeMessage int, message []byte) {
	handleSubscribe(client, message)
	var id, action = getActionAndIdFromClientMessage(message)
	pendingResponses[id] = client
	client.PendingResponses[id] = action
	switch action {
	case ACTION_SERVICE_DOESNT_EXISTS:
		// Send to the corresponding channel the response
		var result MessageData = createErrorResponse(id, "Service doesn't exists")
		rpcMessageHandler <- result
	case ACTION_SEND_FIRST_RESPONSE, ACTION_COMPARE_SERVICES:
		for _, rpc := range rpcList {
			rpc.Connection.WriteMessage(typeMessage, message)
		}
	case ACTION_SEND_ONLY_ONCE, ACTION_SEND_ONLY_TO_PRIMARY:
		rpcList[0].Connection.WriteMessage(typeMessage, message)
	case ACTION_DEPRECATED:
		var result MessageData = createErrorResponse(id, "Service deprecated")
		rpcMessageHandler <- result
	}
}

func createErrorResponse(id float64, message string) MessageData {
	var jsonMessage string = fmt.Sprintf("{\"id\":%.0f,\"result\":\"%s\"}", id, message)
	var messageData MessageData
	messageData.TypeMessage = websocket.TextMessage
	messageData.Message = []byte(jsonMessage)

	var result map[string]interface{}
	json.Unmarshal(messageData.Message, &result)
	messageData.MessageAsJson = result

	messageData.RpcName = "primary"
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
		if err != nil {
			fmt.Println(err)
			continue
		}
		log.Printf("Message received from %s\n", rpcInfo.Name)
		var result map[string]interface{}
		json.Unmarshal(message, &result)
		// log.Printf("%s\n\n", message)
		var res MessageData
		res.TypeMessage = typeMessage
		res.Message = message
		res.MessageAsJson = result
		res.RpcName = rpcInfo.Name
		res.Err = err
		res.TimeReceivedInSeconds = time.Now().Unix()
		rpcMessageHandler <- res
	}
}
