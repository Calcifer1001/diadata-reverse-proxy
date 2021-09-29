package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"strings"
	"websocket"
)

var upgrader = websocket.Upgrader{}

var infuraUrl = "wss://mainnet.infura.io/ws/v3/9bdd9b1d1270497795af3f522ad85091"

//var infuraUrl2 = "wss://mainnet.infura.io/ws/v3/1adea96b97c04c1ab7c0efae5a00d840"
// var infuraUrl = "wss://kovan.infura.io/ws/v3/9bdd9b1d1270497795af3f522ad85091"
var chainstackUrl = "wss://ws-nd-986-369-125.p2pify.com/c669411d9bcc43aa0519602a30346446"
var alchemyUrl = "wss://eth-mainnet.alchemyapi.io/v2/v1bo6tRKiraJ71BVGKmCtWVedAzzNTd6"
var globalSubscriptionHash interface{}

// var infuraUrl2 = "wss://kovan.infura.io/ws/v3/1adea96b97c04c1ab7c0efae5a00d840"

func init() {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
}

type Result struct {
	TypeMessage int
	Message     []byte
	Name        string
	Err         error
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
		// chainstack, _, _ := websocket.DefaultDialer.Dial(chainstackUrl, nil)
		alchemy, _, _ := websocket.DefaultDialer.Dial(alchemyUrl, nil)

		messages := make(chan Result)

		go func() {
			for {
				typeMessage, message, _ := client.ReadMessage()
				log.Printf("Message received from client: %s", message)

				infura.WriteMessage(typeMessage, message)
				// chainstack.WriteMessage(typeMessage, message)
				alchemy.WriteMessage(typeMessage, message)
			}
		}()

		var wssList = make([]WssInfo, 0)
		wssList = append(wssList, *generateWssInfo("infura", *infura))
		// wssList = append(wssList, *generateWssInfo("chainstack", *chainstack))
		wssList = append(wssList, *generateWssInfo("alchemy", *alchemy))
		var sentResults []Result = make([]Result, 0)

		go readMessagesFromRPC(wssList[0], messages)
		go readMessagesFromRPC(wssList[1], messages)

		for {
			var result = <-messages
			for i := 0; i < len(wssList); i++ {
				if result.Name == wssList[i].Name {
					for j := 0; j < len(wssList); j++ {
						if i == j {
							continue
						}
						if ok, index := isResultInResultList(result, wssList[j].Responses); ok {
							result.Message = setSubscriptionHash(result.Message)
							client.WriteMessage(result.TypeMessage, result.Message)
							sentResults = append(sentResults, result)
							wssList[j].Responses = remove(wssList[j].Responses, index)
							// Borrar los responses de donde llego y de donde estaba
						} else {
							if ok, index := isResultInResultList(result, sentResults); ok {
								sentResults = remove(sentResults, index)
							} else {
								wssList[i].Responses = append(wssList[i].Responses, result)
							}
						}
					}
				}
			}
		}
	})
	http.ListenAndServe(":8080", proxy)
}

func remove(slice []Result, s int) []Result {
	return append(slice[:s], slice[s+1:]...)
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
		log.Printf("Message received from %s: %s\n\n", wssInfo.Name, message)
		var res = new(Result)
		res.TypeMessage = typeMessage
		res.Message = message
		res.Name = wssInfo.Name
		res.Err = err
		ch <- *res
	}
}

func setReq(req *http.Request, url *url.URL) {
	req.Host = url.Host
	req.URL.Host = url.Host
	req.URL.Scheme = url.Scheme
	req.RequestURI = ""
}

func getResponse(url *url.URL) (*http.Response, error) {
	fmt.Println(url.RequestURI())
	var req, err = http.NewRequest("GET", url.RequestURI(), nil)
	setReq(req, url)
	s, _, _ := net.SplitHostPort(req.RemoteAddr)
	req.Header.Set("X-Forwarded-For", s)
	resp, err := http.DefaultClient.Do(req)
	return resp, err
}

func handleUrl(url string) {
	c, _, _ := websocket.DefaultDialer.Dial(url, nil)

	go func() {
		for {
			_, message, _ := c.ReadMessage()
			log.Printf("Message received: %s", message)
		}
	}()
}
