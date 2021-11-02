package main

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"testing"
	"websocket"
)

// jsonRequest already declared in allows_comparison_test.go
// primaryUrl already declared in allows_comparison_test.go
// var primaryUrl = "wss://kovan.infura.io/ws/v3/a0bfa51a18b24e1fac45a36481bf7f61"

var typeObject = reflect.TypeOf(make(map[string]interface{}))
var typeString = reflect.TypeOf("")
var typeBoolean = reflect.TypeOf(true)
var typeInterface = reflect.TypeOf(make([]interface{}, 0))
var typesSlice []reflect.Type

func printMap(json map[string]interface{}, t *testing.T) {
	for k, v := range json {
		t.Log(k, ":", v)
	}
}
func webSocketCall(message string) []byte {

	primaryDial, _, _ := websocket.DefaultDialer.Dial(primaryUrl, nil)
	primaryDial.WriteMessage(websocket.TextMessage, []byte(message))
	typeMessage, jsonResponse, _ := primaryDial.ReadMessage()
	expectedTypeMessage := websocket.TextMessage
	if typeMessage != expectedTypeMessage {
		fmt.Printf("Expected: %v, got: %v \n", expectedTypeMessage, typeMessage)
	}
	return jsonResponse
}
func basicValidations(result map[string]interface{}, t *testing.T, dataTypes []reflect.Type) {

	printMap(result, t) //This func is only for viewing the method's responses. Could be deleted

	searchedAttribute := "result"
	validationPass := false
	value, ok := result[searchedAttribute]
	var typeOfValue = reflect.TypeOf(value)

	if !ok {
		t.Errorf("Searched attribute %v not found", searchedAttribute)
	}
	for _, dataType := range dataTypes {
		validationPass = validationPass || dataType == typeOfValue
	}
	if !validationPass {
		t.Errorf("Unespected type of response")
	}
}
func convertToObject(response []byte) map[string]interface{} {
	var result map[string]interface{}
	json.Unmarshal(response, &result)

	return result
}
func TestWeb3ClientVersion(t *testing.T) {

	methodName := "web3_clientVersion"
	params := ""
	messageToSend := jsonRequest
	messageToSend = strings.Replace(messageToSend, "[method]", methodName, -1)
	messageToSend = strings.Replace(messageToSend, "[params]", params, -1)

	var jsonResponse []byte = webSocketCall(messageToSend)
	var jsonObject map[string]interface{} = convertToObject(jsonResponse)
	var typesSlice = make([]reflect.Type, 0)
	typesSlice = append(typesSlice, typeString)
	basicValidations(jsonObject, t, typesSlice)
}
func TestWeb3Sha3(t *testing.T) {

	methodName := "web3_sha3"
	params := "\"0x68656c6c6f20776f726c64\""
	messageToSend := jsonRequest
	messageToSend = strings.Replace(messageToSend, "[method]", methodName, -1)
	messageToSend = strings.Replace(messageToSend, "[params]", params, -1)
	var jsonResponse []byte = webSocketCall(messageToSend)
	var jsonObject map[string]interface{} = convertToObject(jsonResponse)
	var typesSlice = make([]reflect.Type, 0)
	typesSlice = append(typesSlice, typeString)
	basicValidations(jsonObject, t, typesSlice)
}
func TestNetVersion(t *testing.T) {

	methodName := "net_version"
	params := ""
	messageToSend := jsonRequest
	messageToSend = strings.Replace(messageToSend, "[method]", methodName, -1)
	messageToSend = strings.Replace(messageToSend, "[params]", params, -1)
	var jsonResponse []byte = webSocketCall(messageToSend)
	var jsonObject map[string]interface{} = convertToObject(jsonResponse)
	var typesSlice = make([]reflect.Type, 0)
	typesSlice = append(typesSlice, typeString)
	basicValidations(jsonObject, t, typesSlice)
}
func TestNetPeerCount(t *testing.T) {

	methodName := "net_peerCount"
	params := ""
	messageToSend := jsonRequest
	messageToSend = strings.Replace(messageToSend, "[method]", methodName, -1)
	messageToSend = strings.Replace(messageToSend, "[params]", params, -1)
	var jsonResponse []byte = webSocketCall(messageToSend)
	var jsonObject map[string]interface{} = convertToObject(jsonResponse)
	var typesSlice = make([]reflect.Type, 0)
	typesSlice = append(typesSlice, typeString)
	basicValidations(jsonObject, t, typesSlice)
}
func TestEthsyncing(t *testing.T) {

	methodName := "eth_syncing"
	params := ""
	messageToSend := jsonRequest
	messageToSend = strings.Replace(messageToSend, "[method]", methodName, -1)
	messageToSend = strings.Replace(messageToSend, "[params]", params, -1)
	var jsonResponse []byte = webSocketCall(messageToSend)
	var jsonObject map[string]interface{} = convertToObject(jsonResponse)
	var typesSlice = make([]reflect.Type, 0)
	typesSlice = append(typesSlice, typeString)
	typesSlice = append(typesSlice, typeBoolean)
	basicValidations(jsonObject, t, typesSlice)
}
func TestEthHashrate(t *testing.T) {

	methodName := "eth_hashrate"
	params := ""
	messageToSend := jsonRequest
	messageToSend = strings.Replace(messageToSend, "[method]", methodName, -1)
	messageToSend = strings.Replace(messageToSend, "[params]", params, -1)
	var jsonResponse []byte = webSocketCall(messageToSend)
	var jsonObject map[string]interface{} = convertToObject(jsonResponse)
	var typesSlice = make([]reflect.Type, 0)
	typesSlice = append(typesSlice, typeString)
	basicValidations(jsonObject, t, typesSlice)
}
func TestEthgasPrice(t *testing.T) {

	methodName := "eth_gasPrice"
	params := ""
	messageToSend := jsonRequest
	messageToSend = strings.Replace(messageToSend, "[method]", methodName, -1)
	messageToSend = strings.Replace(messageToSend, "[params]", params, -1)
	var jsonResponse []byte = webSocketCall(messageToSend)
	var jsonObject map[string]interface{} = convertToObject(jsonResponse)
	var typesSlice = make([]reflect.Type, 0)
	typesSlice = append(typesSlice, typeString)
	basicValidations(jsonObject, t, typesSlice)
}
func TestEthMining(t *testing.T) {

	methodName := "eth_mining"
	params := ""
	messageToSend := jsonRequest
	messageToSend = strings.Replace(messageToSend, "[method]", methodName, -1)
	messageToSend = strings.Replace(messageToSend, "[params]", params, -1)
	var jsonResponse []byte = webSocketCall(messageToSend)
	var jsonObject map[string]interface{} = convertToObject(jsonResponse)
	var typesSlice = make([]reflect.Type, 0)
	typesSlice = append(typesSlice, typeBoolean)
	basicValidations(jsonObject, t, typesSlice)
}
func TestEthaccounts(t *testing.T) {

	methodName := "eth_accounts"
	params := ""
	messageToSend := jsonRequest
	messageToSend = strings.Replace(messageToSend, "[method]", methodName, -1)
	messageToSend = strings.Replace(messageToSend, "[params]", params, -1)
	var jsonResponse []byte = webSocketCall(messageToSend)
	var jsonObject map[string]interface{} = convertToObject(jsonResponse)
	var typesSlice = make([]reflect.Type, 0)
	typesSlice = append(typesSlice, typeInterface)
	basicValidations(jsonObject, t, typesSlice)
}
func TestEthblockNumber(t *testing.T) {

	methodName := "eth_blockNumber"
	params := ""
	messageToSend := jsonRequest
	messageToSend = strings.Replace(messageToSend, "[method]", methodName, -1)
	messageToSend = strings.Replace(messageToSend, "[params]", params, -1)
	var jsonResponse []byte = webSocketCall(messageToSend)
	var jsonObject map[string]interface{} = convertToObject(jsonResponse)
	var typesSlice = make([]reflect.Type, 0)
	typesSlice = append(typesSlice, typeString)
	basicValidations(jsonObject, t, typesSlice)
}
func TestEthTransactionByHash(t *testing.T) {

	methodName := "eth_getTransactionByHash"
	params := "\"0x83ce0345913f2cac30e1e0d04ceb83bc01bd0c7e28219c2df593bfabaf58d68c\""
	messageToSend := jsonRequest
	messageToSend = strings.Replace(messageToSend, "[method]", methodName, -1)
	messageToSend = strings.Replace(messageToSend, "[params]", params, -1)
	var jsonResponse []byte = webSocketCall(messageToSend)
	var jsonObject map[string]interface{} = convertToObject(jsonResponse)
	var typesSlice = make([]reflect.Type, 0)
	typesSlice = append(typesSlice, typeObject)
	basicValidations(jsonObject, t, typesSlice)
}
func TestEthsubmitHashrate(t *testing.T) {

	methodName := "eth_submitHashrate"
	params := "\"0x0000000000000000000000000000000000000000000000000000000000500000\",\"0x59daa26581d0acd1fce254fb7e85952f4c09d0915afd33d3886cd914bc7d283c\""
	messageToSend := jsonRequest
	messageToSend = strings.Replace(messageToSend, "[method]", methodName, -1)
	messageToSend = strings.Replace(messageToSend, "[params]", params, -1)
	var jsonResponse []byte = webSocketCall(messageToSend)
	var jsonObject map[string]interface{} = convertToObject(jsonResponse)
	var typesSlice = make([]reflect.Type, 0)
	typesSlice = append(typesSlice, typeBoolean)
	basicValidations(jsonObject, t, typesSlice)
}
func TestEthsubmitWork(t *testing.T) {

	methodName := "eth_submitWork"
	params := "\"0x0000000000000001\",\"0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef\",\"0xD1FE5700000000000000000000000000D1FE5700000000000000000000000000\""
	messageToSend := jsonRequest
	messageToSend = strings.Replace(messageToSend, "[method]", methodName, -1)
	messageToSend = strings.Replace(messageToSend, "[params]", params, -1)
	var jsonResponse []byte = webSocketCall(messageToSend)
	var jsonObject map[string]interface{} = convertToObject(jsonResponse)
	var typesSlice = make([]reflect.Type, 0)
	typesSlice = append(typesSlice, typeBoolean)
	basicValidations(jsonObject, t, typesSlice)
}
