package main

import (
	"fmt"
	"reflect"
	"testing"
	"websocket"
)

//We are using two times infura because chainstack doesn't support testnet
var PrimaryUrl = "wss://kovan.infura.io/ws/v3/a0bfa51a18b24e1fac45a36481bf7f61"
var SecondaryUrl = "wss://eth-kovan.alchemyapi.io/v2/B41RjkzXgvqrWxmYaj0aNiDnNfm_NSO4"
var ThirdUrl = "wss://kovan.infura.io/ws/v3/be1a3f5f45994142bb67759b9fef28c5"
var searchedAttribute = "result"

func Caller(Message string, NodeUrl string) []byte {

	Dial, _, _ := websocket.DefaultDialer.Dial(NodeUrl, nil)
	Dial.WriteMessage(websocket.TextMessage, []byte(Message))
	typeMessage, jsonResponse, _ := Dial.ReadMessage()
	expectedTypeMessage := websocket.TextMessage
	if typeMessage != expectedTypeMessage {
		fmt.Printf("Expected: %v, got: %v \n", expectedTypeMessage, typeMessage)
	}
	return jsonResponse
}
func ExtractResult(jsonObject map[string]interface{}, attribute string, t *testing.T) interface{} {
	value, ok := jsonObject[attribute]
	if !ok {
		t.Errorf("Searched attribute %v not found", attribute)
	}
	return value
}
func CompareResults(first, second, third interface{}, t *testing.T) interface{} {
	switch {
	case first == second, first == third:
		t.Log(first)
		return first
	case second == third:
		t.Log(second)
		return second
	default:
		fmt.Printf("No matches between responses: %v, %v, %v", first, second, third)
		return nil
	}
}
func ValidateResult(result interface{}, t *testing.T, dataTypes []reflect.Type) {
	var typeOfValue = reflect.TypeOf(result)
	validationPass := false

	for _, dataType := range dataTypes {
		validationPass = validationPass || dataType == typeOfValue
	}
	if !validationPass {
		t.Errorf("Unespected type of response")
	}
}

func TestEthprotocolVersion(t *testing.T) {
	messageToSend := "{\"method\":\"eth_protocolVersion\",\"params\":[],\"jsonrpc\":\"2.0\",\"id\":67}"

	jsonResponse1 := Caller(messageToSend, PrimaryUrl)
	jsonResponse2 := Caller(messageToSend, SecondaryUrl)
	jsonResponse3 := Caller(messageToSend, ThirdUrl)

	var jsonObject1 map[string]interface{} = convertToObject(jsonResponse1)
	var jsonObject2 map[string]interface{} = convertToObject(jsonResponse2)
	var jsonObject3 map[string]interface{} = convertToObject(jsonResponse3)

	result1 := ExtractResult(jsonObject1, searchedAttribute, t)
	result2 := ExtractResult(jsonObject2, searchedAttribute, t)
	result3 := ExtractResult(jsonObject3, searchedAttribute, t)

	comparedResult := CompareResults(result1, result2, result3, t)

	typesSlice = nil
	typesSlice = append(typesSlice, typeString)

	ValidateResult(comparedResult, t, typesSlice)
}

func TestNetListening(t *testing.T) {
	messageToSend := "{\"method\":\"net_listening\",\"params\":[],\"jsonrpc\":\"2.0\",\"id\":67}"

	jsonResponse1 := Caller(messageToSend, PrimaryUrl)
	jsonResponse2 := Caller(messageToSend, SecondaryUrl)
	jsonResponse3 := Caller(messageToSend, ThirdUrl)

	var jsonObject1 map[string]interface{} = convertToObject(jsonResponse1)
	var jsonObject2 map[string]interface{} = convertToObject(jsonResponse2)
	var jsonObject3 map[string]interface{} = convertToObject(jsonResponse3)

	result1 := ExtractResult(jsonObject1, searchedAttribute, t)
	result2 := ExtractResult(jsonObject2, searchedAttribute, t)
	result3 := ExtractResult(jsonObject3, searchedAttribute, t)

	comparedResult := CompareResults(result1, result2, result3, t)

	typesSlice = nil
	typesSlice = append(typesSlice, typeBoolean)

	ValidateResult(comparedResult, t, typesSlice)
}
