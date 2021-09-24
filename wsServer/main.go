package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"reflect"
	"websocket"
)

var upgrader = websocket.Upgrader{}

//var infuraUrl = "wss://mainnet.infura.io/ws/v3/9bdd9b1d1270497795af3f522ad85091"
//var infuraUrl2 = "wss://mainnet.infura.io/ws/v3/1adea96b97c04c1ab7c0efae5a00d840"
var infuraUrl = "wss://kovan.infura.io/ws/v3/9bdd9b1d1270497795af3f522ad85091"
var infuraUrl2 = "wss://kovan.infura.io/ws/v3/1adea96b97c04c1ab7c0efae5a00d840"

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
		fmt.Println("Hola")
		client, _ := upgrader.Upgrade(rw, req, nil)
		// defer ws.Close()
		infura, _, _ := websocket.DefaultDialer.Dial(infuraUrl, nil)
		infura2, _, _ := websocket.DefaultDialer.Dial(infuraUrl2, nil)
		fmt.Println("1")
		messages := make(chan Result)

		go func() {
			for {
				typeMessage, message, _ := client.ReadMessage()
				log.Printf("Message received from client: %s", message)

				infura.WriteMessage(typeMessage, message)
				infura2.WriteMessage(typeMessage, message)
			}
		}()

		var wssList = make([]WssInfo, 0)
		wssList = append(wssList, *generateWssInfo("infura1", *infura))
		wssList = append(wssList, *generateWssInfo("infura2", *infura2))

		go readMessagesFromRPC(wssList[0], messages)
		go readMessagesFromRPC(wssList[1], messages)

		for {
			var result = <-messages
			for i := 0; i < len(wssList); i++ {
				if result.Name == wssList[i].Name {
					fmt.Printf("Handling response from %s\n", wssList[i].Name)
					wssList[i].Responses = append(wssList[i].Responses, result)
					for j := 0; j < len(wssList); j++ {
						if i == j {
							continue
						}
						if isResultInResultList(result, wssList[j].Responses) {
							client.WriteMessage(result.TypeMessage, result.Message)
							// Borrar los responses de donde llego y de donde estaba
						}
					}
				}
			}
		}
	})
	http.ListenAndServe(":8080", proxy)
}

func isResultInResultList(result Result, responseList []Result) bool {
	for _, otherResponses := range responseList {
		if equals(result.Message, otherResponses.Message) {
			return true
		}
	}
	return false
}

func extractResult(message []byte) map[string]interface{} {
	var result map[string]interface{}
	json.Unmarshal(message, &result)

	var params = result["params"]
	if params == nil {
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
	var result1 = extractResult(message1)
	var result2 = extractResult(message2)
	return reflect.DeepEqual(result1, result2)
}

func readMessagesFromRPC(wssInfo WssInfo, ch chan Result) {
	for {
		typeMessage, message, err := wssInfo.Connection.ReadMessage()
		log.Printf("Message received from %s: %s\n", wssInfo.Name, message)
		var res = new(Result)
		res.TypeMessage = typeMessage
		res.Message = message
		res.Name = wssInfo.Name
		res.Err = err
		ch <- *res
		//client.WriteMessage(typeMessage, message)
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
