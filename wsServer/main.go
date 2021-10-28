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

var upgrader = websocket.Upgrader{}

var primaryUrl = "wss://mainnet.infura.io/ws/v3/9bdd9b1d1270497795af3f522ad85091"

//var infuraUrl2 = "wss://mainnet.infura.io/ws/v3/1adea96b97c04c1ab7c0efae5a00d840"
// var infuraUrl = "wss://kovan.infura.io/ws/v3/9bdd9b1d1270497795af3f522ad85091"
var secondaryUrl = "wss://ws-nd-986-369-125.p2pify.com/c669411d9bcc43aa0519602a30346446"
var tertiaryUrl = "wss://eth-mainnet.alchemyapi.io/v2/v1bo6tRKiraJ71BVGKmCtWVedAzzNTd6"
var primaryServer *websocket.Conn
var secondaryServer *websocket.Conn
var tertiaryServer *websocket.Conn
var globalSubscriptionHash interface{}
var globalResponseTimeoutInSeconds int64 = 15
var clientPendingResponses map[int]int

// var infuraUrl2 = "wss://kovan.infura.io/ws/v3/1adea96b97c04c1ab7c0efae5a00d840"

func init() {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
}

type Result struct {
	TypeMessage           int
	Message               []byte
	Name                  string
	Err                   error
	TimeReceivedInSeconds int64
}

type WssInfo struct {
	Name       string
	Connection websocket.Conn
	Responses  []Result
}

func main() {
	initializeServiceLists()

	primaryServer, _, _ := websocket.DefaultDialer.Dial(primaryUrl, nil)
	secondaryServer, _, _ := websocket.DefaultDialer.Dial(secondaryUrl, nil)
	tertiaryServer, _, _ := websocket.DefaultDialer.Dial(tertiaryUrl, nil)

	proxy := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		client, _ := upgrader.Upgrade(rw, req, nil)
		defer client.Close()
		messages := make(chan Result)

		var wssList = make([]WssInfo, 0)
		wssList = append(wssList, *generateWssInfo("primary", *primaryServer))
		wssList = append(wssList, *generateWssInfo("secondary", *secondaryServer))
		wssList = append(wssList, *generateWssInfo("tertiary", *tertiaryServer))
		var sentMessages = make([]Result, 0)

		// Receive a message from the client it sends it to the corresponding rpc services
		go func() {
			for {
				typeMessage, message, _ := client.ReadMessage()
				log.Printf("Message received from client: %s", message)

				sendMessage(typeMessage, message, messages)
			}
		}()

		for _, wss := range wssList {
			go readMessagesFromRPC(wss, messages)
		}

		go func() {
			for {
				wssList = cleanResponsesByTimeout(wssList)
				time.Sleep(5 * time.Second)
			}
		}()

		for {
			var result Result = <-messages
			handleMessage(client, wssList, result, &sentMessages)
		}
	})
	http.ListenAndServe(":8080", proxy)
}

func sendMessage(typeMessage int, message []byte, messages chan Result) {
	var id, action = getActionAndIdFromMessage(message)
	clientPendingResponses[id] = action

	switch action {
	case ACTION_SERVICE_DOESNT_EXISTS:
		// Send to the corresponding channel the response
		var result Result = createErrorResult(id, "Service doesn't exists")
		messages <- result
	case ACTION_SEND_FIRST_RESPONSE:
		primaryServer.WriteMessage(typeMessage, message)
		secondaryServer.WriteMessage(typeMessage, message)
		tertiaryServer.WriteMessage(typeMessage, message)
	case ACTION_COMPARE_SERVICES:
		primaryServer.WriteMessage(typeMessage, message)
		secondaryServer.WriteMessage(typeMessage, message)
		tertiaryServer.WriteMessage(typeMessage, message)
	case ACTION_SEND_ONLY_ONCE:
		primaryServer.WriteMessage(typeMessage, message)
	case ACTION_SEND_ONLY_TO_PRIMARY:
		primaryServer.WriteMessage(typeMessage, message)
	case ACTION_DEPRECATED:
		var result Result = createErrorResult(id, "Service deprecated")
		messages <- result
	}
}

func createErrorResult(id int, message string) Result {
	var jsonMessage string = fmt.Sprintf("{\"id\":%d,\"result\":%s}", id, message)
	var result Result
	result.TypeMessage = websocket.TextMessage
	result.Message = []byte(jsonMessage)
	result.Name = "primary"
	result.Err = nil
	result.TimeReceivedInSeconds = 0
	return result
}

func handleMessage(client *websocket.Conn, wssList []WssInfo, result Result, sentMessages []Result) {
	var resultMap map[string]interface{}
	json.Unmarshal(result.Message, &result)
	var id int = resultMap["id"].(int)
	if _, ok := clientPendingResponses[id]; !ok {
		fmt.Println("Message already sent. Discarding")
		return
	}
	var sent = false
	switch clientPendingResponses[id] {
	case ACTION_SEND_FIRST_RESPONSE:
		client.WriteMessage(result.TypeMessage, result.Message)
		sent = true
	case ACTION_COMPARE_SERVICES:
		sent = compareAndSendOnMatch(client, wssList, result, &sentMessages)
	case ACTION_SEND_ONLY_ONCE:
		client.WriteMessage(result.TypeMessage, result.Message)
		sent = true
	case ACTION_SEND_ONLY_TO_PRIMARY:
		client.WriteMessage(result.TypeMessage, result.Message)
		sent = true
	}
	if sent {
		delete(clientPendingResponses, id)
	}
}

func compareAndSendOnMatch(client *websocket.Conn, wssList []WssInfo, result Result, sentMessages []Result) bool {
	for i := 0; i < len(wssList); i++ {
		if result.Name == wssList[i].Name {
			wssList[i].Responses = append(wssList[i].Responses, result)
			for j := 0; j < len(wssList); j++ {
				if i == j {
					continue
				}
				if ok, _ := isResultInResultList(result, sentMessages); ok {
					continue
				}
				if ok, _ := isResultInResultList(result, wssList[j].Responses); ok {
					result.Message = setSubscriptionHash(result.Message)
					client.WriteMessage(result.TypeMessage, result.Message)
					sentMessages = append(sentMessages, result)
					return true
				}
			}
		}
	}
	return false
}

func cleanResponsesByTimeout(wssList []WssInfo) []WssInfo {
	for i := 0; i < len(wssList); i++ {
		var wss = wssList[i]
		log.Printf("Checking wss %s with %d elements\n", wss.Name, len(wss.Responses))
		var indexAux = 0
		for j := 0; j < len(wss.Responses); j++ {
			// If shouldn't be removed, it's set to the first part of the slice.
			if time.Now().Unix()-wss.Responses[j].TimeReceivedInSeconds < globalResponseTimeoutInSeconds {
				wss.Responses[indexAux] = wss.Responses[j]
				indexAux++
			}
		}
		// Removing all the elements that are after the ones that should stay
		wss.Responses = wss.Responses[:indexAux]
		wssList[i] = wss
	}
	fmt.Println()
	return wssList

}

func isResultInResultList(result Result, responseList []Result) (bool, int) {
	for i, otherResponses := range responseList {
		if equals(result.Message, otherResponses.Message) {
			return true, i
		}
	}
	return false, -1
}

func setSubscriptionHash(message []byte) []byte {
	var result map[string]interface{}
	json.Unmarshal(message, &result)

	var params = result["params"]
	if params == nil {
		return message
	}

	params.(map[string]interface{})["subscription"] = globalSubscriptionHash
	var output, err = json.Marshal(result)
	if err != nil {
		fmt.Println("Error seteando subscription hash")
	}
	return output
}

func extractResult(message []byte) map[string]interface{} {
	var result map[string]interface{}
	json.Unmarshal(message, &result)

	var params = result["params"]
	if params == nil {
		if globalSubscriptionHash == nil {
			fmt.Println("Setting global subscription hash")
			globalSubscriptionHash = result["result"]
			fmt.Printf("Subscription hash set to %s\n\n", globalSubscriptionHash)
		}
		return nil
	}
	return params.(map[string]interface{})["result"].(map[string]interface{})
}

func generateWssInfo(name string, connection websocket.Conn) *WssInfo {
	var wssInfo = new(WssInfo)
	wssInfo.Name = name
	wssInfo.Connection = connection
	wssInfo.Responses = make([]Result, 0)
	return wssInfo
}

func equals(message1, message2 []byte) bool {
	var result1 map[string]interface{} = extractResult(message1)
	var result2 map[string]interface{} = extractResult(message2)

	// Cannot use deep equal because some server return message with uppercase and some with lowercase
	return strings.EqualFold(fmt.Sprint(result1), fmt.Sprint(result2))
}

func readMessagesFromRPC(wssInfo WssInfo, ch chan Result) {
	for {
		typeMessage, message, err := wssInfo.Connection.ReadMessage()
		log.Printf("Message received from %s\n", wssInfo.Name)
		// log.Printf("%s\n\n", message)
		var res = new(Result)
		res.TypeMessage = typeMessage
		res.Message = message
		res.Name = wssInfo.Name
		res.Err = err
		res.TimeReceivedInSeconds = time.Now().Unix()
		ch <- *res
	}
}
