package main

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
	"websocket"
)

var primaryUrl = "wss://kovan.infura.io/ws/v3/a0bfa51a18b24e1fac45a36481bf7f61"

var typeObject = reflect.TypeOf(make(map[string]interface{}))
var typeString = reflect.TypeOf("")
var typeBoolean = reflect.TypeOf(true)

func printMap(json map[string]interface{}, t *testing.T) {
	for k, v := range json {
		t.Log(k, ":", v)
	}
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

func basicValidations(result map[string]interface{}, t *testing.T, dataType reflect.Type) {
	printMap(result, t)
	var searchedAttribute string = "result"
	value, ok := result[searchedAttribute]
	if !ok {
		t.Errorf("Searched attribute %v not found", searchedAttribute)
	}
	var typeOfValue = reflect.TypeOf(value)
	if dataType != typeOfValue {
		t.Errorf("Value %v is not a %v", value, dataType)
	}

}

func convertToObject(response []byte) map[string]interface{} {
	var result map[string]interface{}
	json.Unmarshal(response, &result)

	return result
}

func TestWeb3ClientVersion(t *testing.T) {
	var messageToSend string = "{\"method\":\"web3_clientVersion\",\"params\":[],\"jsonrpc\":\"2.0\",\"id\":67}"

	var jsonResponse []byte = WebSocketCall(messageToSend)
	var jsonObject map[string]interface{} = convertToObject(jsonResponse)
	basicValidations(jsonObject, t, typeString)
}
func TestWeb3Sha3(t *testing.T) {
	var messageToSend string = "{\"method\":\"web3_sha3\",\"params\":[\"0x68656c6c6f20776f726c64\"],\"jsonrpc\":\"2.0\",\"id\":67}"

	var jsonResponse []byte = WebSocketCall(messageToSend)
	var jsonObject map[string]interface{} = convertToObject(jsonResponse)
	basicValidations(jsonObject, t, typeString)
}
func TestNetVersion(t *testing.T) {

	// messageToSend := "{\"method\":\"net_version\",\"params\":[],\"jsonrpc\":\"2.0\",\"id\":67}"

	// jsonResponse := WebSocketCall(messageToSend)

	// validate(jsonResponse, t)
}
func TestNetPeerCount(t *testing.T) {

	// messageToSend := "{\"method\":\"net_peerCount\",\"params\":[],\"jsonrpc\":\"2.0\",\"id\":67}"

	// jsonResponse := WebSocketCall(messageToSend)

	// validate(jsonResponse, t)
}

// func TestEthsyncing(t *testing.T) {
// //devuelve datos o un bool. ver como chequearlo
// 	messageToSend := "{\"method\":\"eth_syncing\",\"params\":[],\"jsonrpc\":\"2.0\",\"id\":67}"

// 	jsonResponse := WebSocketCall(messageToSend)

// 	PrintResult(jsonResponse, t)
// }
func TestEthHashrate(t *testing.T) {

	// messageToSend := "{\"method\":\"eth_hashrate\",\"params\":[],\"jsonrpc\":\"2.0\",\"id\":67}"

	// jsonResponse := WebSocketCall(messageToSend)

	// validate(jsonResponse, t)
}
func TestEthgasPrice(t *testing.T) {

	// messageToSend := "{\"method\":\"eth_gasPrice\",\"params\":[],\"jsonrpc\":\"2.0\",\"id\":67}"

	// jsonResponse := WebSocketCall(messageToSend)

	// validate(jsonResponse, t)
}

// func TestEthMining(t *testing.T) {
// 	messageToSend := "{\"method\":\"eth_mining\",\"params\":[],\"jsonrpc\":\"2.0\",\"id\":67}"
// 	jsonResponse := WebSocketCall(messageToSend)
// 	PrintResult(jsonResponse, t)
// }

// func TestEthaccounts(t *testing.T) {
// // chequear que puede devolver vacío
// 	messageToSend := "{\"method\":\"eth_accounts\",\"params\":[],\"jsonrpc\":\"2.0\",\"id\":67}"

// 	jsonResponse := WebSocketCall(messageToSend)

// 	PrintResult(jsonResponse, t)
// }

func TestEthblockNumber(t *testing.T) {

	// messageToSend := "{\"method\":\"eth_blockNumber\",\"params\":[],\"jsonrpc\":\"2.0\",\"id\":67}"

	// jsonResponse := WebSocketCall(messageToSend)

	// validate(jsonResponse, t)
}

func TestEthTransactionByHash(t *testing.T) {
	//devuelve otro map, consultar como revisarlo
	var messageToSend string = "{\"method\":\"eth_getTransactionByHash\",\"params\":[\"0x83ce0345913f2cac30e1e0d04ceb83bc01bd0c7e28219c2df593bfabaf58d68c\"],\"jsonrpc\":\"2.0\",\"id\":67}"

	var jsonResponse []byte = WebSocketCall(messageToSend)
	var jsonObject map[string]interface{} = convertToObject(jsonResponse)

	basicValidations(jsonObject, t, typeObject)

}

func TestEthsubmitHashrate(t *testing.T) {
	//devuelve un bool ver como chequearlo
	// messageToSend := "{\"method\":\"eth_submitHashrate\",\"params\":[\"0x0000000000000000000000000000000000000000000000000000000000500000\",\"0x59daa26581d0acd1fce254fb7e85952f4c09d0915afd33d3886cd914bc7d283c\"],\"jsonrpc\":\"2.0\",\"id\":67}"

	// jsonResponse := WebSocketCall(messageToSend)

	// validate(jsonResponse, t, typeBoolean)
}

func TestEthsubmitWork(t *testing.T) {
	//devuelve un bool ver como chequearlo
	// messageToSend := "{\"method\":\"eth_submitWork\",\"params\":[\"0x0000000000000001\",\"0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef\",\"0xD1FE5700000000000000000000000000D1FE5700000000000000000000000000\"],\"jsonrpc\":\"2.0\",\"id\":67}"

	// jsonResponse := WebSocketCall(messageToSend)

	// validate(jsonResponse, t)

}
