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
var messagePart1 = "{\"method\":\""
var messagePart2 = "\",\"params\":["
var messagePart3 = "],\"jsonrpc\":\"2.0\",\"id\":67}"

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
		// t.Log(first)
		return first
	case second == third:
		// t.Log(second)
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

	methodName := "eth_protocolVersion"
	params := ""
	messageToSend := messagePart1 + methodName + messagePart2 + params + messagePart3

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

	methodName := "net_listening"
	params := ""
	messageToSend := messagePart1 + methodName + messagePart2 + params + messagePart3

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

	fmt.Printf("value of response: %v \n", comparedResult) //this line is for work-in-progress test. Will be deleted
}

func TestEthGetBalance(t *testing.T) {

	methodName := "eth_getBalance"
	params := "\"0xa7719d2eD3849D3CD10991b91f1E8D9a2044eD45\""
	messageToSend := messagePart1 + methodName + messagePart2 + params + messagePart3

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

	fmt.Printf("value of response: %v \n", comparedResult) //this line is for work-in-progress test. Will be deleted
}
func TestEthGetStorageAt(t *testing.T) {

	methodName := "eth_getStorageAt"
	params := "\"0xb451c6835515f8a08ecc4cbc5c5dcb238a48f7b4\",\"0x0\",\"latest\"" //this is the address of a test token we used
	messageToSend := messagePart1 + methodName + messagePart2 + params + messagePart3

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

	fmt.Printf("value of response: %v \n", comparedResult) //this line is for work-in-progress test. Will be deleted
}

func TestEthGetTransactionCount(t *testing.T) {

	methodName := "eth_getTransactionCount"
	params := "\"0xa7719d2eD3849D3CD10991b91f1E8D9a2044eD45\",\"latest\""
	messageToSend := messagePart1 + methodName + messagePart2 + params + messagePart3

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

	fmt.Printf("value of response: %v \n", comparedResult) //this line is for work-in-progress test. Will be deleted
}

func TestEthGetBlockTransactionCountByHash(t *testing.T) {

	methodName := "eth_getBlockTransactionCountByHash"
	params := "\"0xb65b2f91f066fdc47a71652d3f0aed4c95f6f2a82f028cc7d8e6cc8b2c6ec11f\""
	messageToSend := messagePart1 + methodName + messagePart2 + params + messagePart3

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

	fmt.Printf("value of response: %v \n", comparedResult) //this line is for work-in-progress test. Will be deleted
}

func TestEthGetBlockTransactionCountByNumber(t *testing.T) {

	methodName := "eth_getBlockTransactionCountByNumber"
	params := "\"0x1ACB6E3\""
	messageToSend := messagePart1 + methodName + messagePart2 + params + messagePart3

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

	fmt.Printf("value of response: %v \n", comparedResult) //this line is for work-in-progress test. Will be deleted
}

func TestEthGetUncleCountByBlockHash(t *testing.T) {

	methodName := "eth_getUncleCountByBlockHash"
	params := "\"0xb65b2f91f066fdc47a71652d3f0aed4c95f6f2a82f028cc7d8e6cc8b2c6ec11f\""
	messageToSend := messagePart1 + methodName + messagePart2 + params + messagePart3

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

	fmt.Printf("value of response: %v \n", comparedResult) //this line is for work-in-progress test. Will be deleted
}
func TestEthGetUncleCountByBlockNumber(t *testing.T) {

	methodName := "eth_getUncleCountByBlockNumber"
	params := "\"0x1ACB6E3\""
	messageToSend := messagePart1 + methodName + messagePart2 + params + messagePart3

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

	fmt.Printf("value of response: %v \n", comparedResult) //this line is for work-in-progress test. Will be deleted
}

func TestEthGetCode(t *testing.T) {

	methodName := "eth_getCode"
	params := "\"0xb451c6835515f8a08ecc4cbc5c5dcb238a48f7b4\",\"latest\"" //this is the address of a test token we used
	messageToSend := messagePart1 + methodName + messagePart2 + params + messagePart3

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

	fmt.Printf("value of response: %v \n", comparedResult) //this line is for work-in-progress test. Will be deleted
}
