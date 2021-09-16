package main

import (
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

func init() {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
}

func main() {
	proxy := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		fmt.Println("Hola")
		ws, _ := upgrader.Upgrade(rw, req, nil)
		// defer ws.Close()
		c, _, _ := websocket.DefaultDialer.Dial(infuraUrl, nil)
		fmt.Println("1")
		go func() {
			for {
				typeMessage, message, _ := ws.ReadMessage()
				log.Printf("Message received from ws: %s", message)
				c.WriteMessage(typeMessage, message)
			}
		}()
		go func() {
			for {
				typeMessage, message, _ := c.ReadMessage()
				log.Printf("Message received from c: %s", message)
				ws.WriteMessage(typeMessage, message)
			}
		}()
		fmt.Println("2")
	})
	http.ListenAndServe(":8080", proxy)
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
