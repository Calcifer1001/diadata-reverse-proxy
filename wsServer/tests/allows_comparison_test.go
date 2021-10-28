package main

// import (
// 	"encoding/json"
// 	"fmt"
// 	"testing"
// 	"websocket"
// )

// //We are using two times infura because chainstack doesn't support testnet
// var PrimaryUrl = "wss://kovan.infura.io/ws/v3/a0bfa51a18b24e1fac45a36481bf7f61"
// var SecondaryUrl = "wss://eth-kovan.alchemyapi.io/v2/B41RjkzXgvqrWxmYaj0aNiDnNfm_NSO4"
// var ThirdUrl = "wss://kovan.infura.io/ws/v3/be1a3f5f45994142bb67759b9fef28c5"

// func Caller(Message string, NodeUrl string) []byte {

// 	Dial, _, _ := websocket.DefaultDialer.Dial(NodeUrl, nil)
// 	Dial.WriteMessage(websocket.TextMessage, []byte(Message))
// 	typeMessage, jsonResponse, _ := Dial.ReadMessage()
// 	expectedTypeMessage := websocket.TextMessage
// 	if typeMessage != expectedTypeMessage {
// 		fmt.Printf("Expected: %v, got: %v \n", expectedTypeMessage, typeMessage)
// 	}
// 	return jsonResponse
// }

// func TestEthprotocolVersion(t *testing.T) {
// 	messageToSend := "{\"method\":\"eth_protocolVersion\",\"params\":[],\"jsonrpc\":\"2.0\",\"id\":67}"

// 	jsonResponse1 := Caller(messageToSend, PrimaryUrl)
// 	jsonResponse2 := Caller(messageToSend, SecondaryUrl)
// 	jsonResponse3 := Caller(messageToSend, ThirdUrl)

// 	var result1, result2, result3 map[string]interface{}
// 	json.Unmarshal(jsonResponse1, &result1)
// 	json.Unmarshal(jsonResponse2, &result2)
// 	json.Unmarshal(jsonResponse3, &result3)

// 	response1 := result1["result"]

// 	// PrintResult(jsonResponse1, t)
// 	// PrintResult(jsonResponse2, t)
// 	// PrintResult(jsonResponse3, t)
// }
