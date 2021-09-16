package main

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"websocket"
)

var upgrader = websocket.Upgrader{}
var infuraUrl = "wss://mainnet.infura.io/ws/v3/9bdd9b1d1270497795af3f522ad85091"
var infuraUrl2 = "wss://mainnet.infura.io/ws/v3/1adea96b97c04c1ab7c0efae5a00d840"

func init() {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
}

type Result struct {
	TypeMessage int
	Message     []byte
	Name        string
	Err         error
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
		go readMessagesFromRPC(infura, messages, "infura")
		go readMessagesFromRPC(infura2, messages, "infura2")
		fmt.Println("2")
		var infura1Responses = make([]Result, 0)
		var infura2Responses = make([]Result, 0)
		for {
			var result = <-messages

			if result.Name == "infura" {
				fmt.Println("Handling response from infura")
				infura1Responses = append(infura1Responses, result)
				for _, infura2Result := range infura2Responses {
					if bytes.Compare(result.Message, infura2Result.Message) == 0 {
						client.WriteMessage(result.TypeMessage, result.Message)
					}
				}
			}

			if result.Name == "infura2" {
				fmt.Println("Handling response from infura2")
				infura2Responses = append(infura2Responses, result)
				for _, infura1Result := range infura1Responses {
					if bytes.Compare(result.Message, infura1Result.Message) == 0 {
						client.WriteMessage(result.TypeMessage, result.Message)
					}
				}
			}

		}
	})
	http.ListenAndServe(":8080", proxy)
}

func readMessagesFromRPC(conn *websocket.Conn, ch chan Result, name string) {
	for {
		typeMessage, message, err := conn.ReadMessage()
		log.Printf("Message received from %s: %s\n", name, message)
		var res = new(Result)
		res.TypeMessage = typeMessage
		res.Message = message
		res.Name = name
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
