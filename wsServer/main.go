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

var infuraUrl = "wss://mainnet.infura.io/ws/v3/9bdd9b1d1270497795af3f522ad85091"

//var infuraUrl2 = "wss://mainnet.infura.io/ws/v3/1adea96b97c04c1ab7c0efae5a00d840"
// var infuraUrl = "wss://kovan.infura.io/ws/v3/9bdd9b1d1270497795af3f522ad85091"
var chainstackUrl = "wss://ws-nd-986-369-125.p2pify.com/c669411d9bcc43aa0519602a30346446"
var alchemyUrl = "wss://eth-mainnet.alchemyapi.io/v2/v1bo6tRKiraJ71BVGKmCtWVedAzzNTd6"
var globalSubscriptionHash interface{}
var globalResponseTimeoutInSeconds int64 = 15

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

	proxy := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		client, _ := upgrader.Upgrade(rw, req, nil)
		// defer ws.Close()
		infura, _, _ := websocket.DefaultDialer.Dial(infuraUrl, nil)
		chainstack, _, _ := websocket.DefaultDialer.Dial(chainstackUrl, nil)
		alchemy, _, _ := websocket.DefaultDialer.Dial(alchemyUrl, nil)

		messages := make(chan Result)

		go func() {
			for {
				typeMessage, message, _ := client.ReadMessage()
				log.Printf("Message received from client: %s", message)

				infura.WriteMessage(typeMessage, message)
				chainstack.WriteMessage(typeMessage, message)
				alchemy.WriteMessage(typeMessage, message)
			}
		}()

		var wssList = make([]WssInfo, 0)
		wssList = append(wssList, *generateWssInfo("infura", *infura))
		wssList = append(wssList, *generateWssInfo("chainstack", *chainstack))
		wssList = append(wssList, *generateWssInfo("alchemy", *alchemy))
		var sentMessages = make([]Result, 0)

		for _, wss := range wssList {
			go readMessagesFromRPC(wss, messages)
		}

		go func() {
			for true {
				wssList = cleanResponsesByTimeout(wssList)
				time.Sleep(5 * time.Second)
			}
		}()

		for {
			var result = <-messages
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
						}
					}
				}
			}
		}
	})
	http.ListenAndServe(":8080", proxy)
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
