package main

import (
	"encoding/json"
	"fmt"
	"testing"
	"websocket"
)

var primaryUrl = "wss://kovan.infura.io/ws/v3/a0bfa51a18b24e1fac45a36481bf7f61"

func printMap(json map[string]interface{}, t *testing.T) {
	for k, v := range json {
		t.Log(k, ":", v)
	}
}

func TestWeb3ClientVersion(t *testing.T) {
	messageToSend := "{\"method\":\"web3_clientVersion\",\"params\":[],\"jsonrpc\":\"2.0\",\"id\":67}"

	jsonResponse := WebSocketCall(messageToSend)

	PrintResult(jsonResponse, t)
}

func TestWeb3Sha3(t *testing.T) {
	messageToSend := "{\"method\":\"web3_sha3\",\"params\":[\"0x68656c6c6f20776f726c64\"],\"jsonrpc\":\"2.0\",\"id\":67}"

	jsonResponse := WebSocketCall(messageToSend)

	PrintResult(jsonResponse, t)
}

func TestNetVersion(t *testing.T) {

	messageToSend := "{\"method\":\"net_version\",\"params\":[],\"jsonrpc\":\"2.0\",\"id\":67}"

	jsonResponse := WebSocketCall(messageToSend)

	PrintResult(jsonResponse, t)

}

func WebSocketCall(Message string) []byte {

	primaryDial, _, _ := websocket.DefaultDialer.Dial(primaryUrl, nil)
	primaryDial.WriteMessage(websocket.TextMessage, []byte(Message))
	typeMessage, jsonResponse, _ := primaryDial.ReadMessage()

	expectedTypeMessage := websocket.TextMessage
	if typeMessage != expectedTypeMessage {
		fmt.Printf("Expected: %v, got: %v \n", expectedTypeMessage, typeMessage)
	}

	return jsonResponse
}

func PrintResult(Response []byte, t *testing.T) {
	var result map[string]interface{}
	json.Unmarshal(Response, &result)
	printMap(result, t)
	searchedAttribute := "result"
	value, ok := result["result"]
	if !ok {
		t.Errorf("Searched attribute %v not found", searchedAttribute)
	}
	if len(value.(string)) == 0 {
		t.Errorf("Value %v is not a string", value)
	}

}
