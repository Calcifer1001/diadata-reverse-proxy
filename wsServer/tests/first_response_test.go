package main

import (
	"encoding/json"
	"strings"
	"testing"
	"websocket"
)

var primaryUrl = "wss://kovan.infura.io/ws/v3/a0bfa51a18b24e1fac45a36481bf7f61"

func IsSuperAnimal(animal string) bool {
	return strings.ToLower(animal) == "gopher"
}

func TestIsSuperAnimal(t *testing.T) {
	expected := true
	got := IsSuperAnimal("gopher")
	if got != expected {
		t.Errorf("Expected: %v, got: %v", expected, got)
	}
}

func TestWeb3ClientVersion(t *testing.T) {
	messageToSend := "{\"method\":\"web3_clientVersion\",\"params\":[],\"jsonrpc\":\"2.0\",\"id\":67}"
	primaryDial, _, _ := websocket.DefaultDialer.Dial(primaryUrl, nil)
	primaryDial.WriteMessage(websocket.TextMessage, []byte(messageToSend))
	typeMessage, jsonResponse, _ := primaryDial.ReadMessage()

	expectedTypeMessage := websocket.TextMessage
	if typeMessage != expectedTypeMessage {
		t.Errorf("Expected: %v, got: %v", expectedTypeMessage, typeMessage)
	}

	var result map[string]interface{}
	json.Unmarshal(jsonResponse, &result)
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

func printMap(json map[string]interface{}, t *testing.T) {
	for k, v := range json {
		t.Log(k, ":", v)
	}
}
