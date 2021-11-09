package main

import (
	"fmt"
	"os"
	"reflect"
	"strings"
	"testing"
	"websocket"
)

//We are using two times infura because chainstack doesn't support testnet
var primaryUrl = "wss://kovan.infura.io/ws/v3/a0bfa51a18b24e1fac45a36481bf7f61"
var secondaryUrl = "wss://eth-kovan.alchemyapi.io/v2/B41RjkzXgvqrWxmYaj0aNiDnNfm_NSO4"
var tertiaryUrl = "wss://kovan.infura.io/ws/v3/be1a3f5f45994142bb67759b9fef28c5"
var urlSlice = make([]string, 0)
var searchedAttribute = "result"
var jsonRequest = "{\"method\":\"[method]\",\"params\":[[params]],\"jsonrpc\":\"2.0\",\"id\":67}"

func caller(message string, nodeUrl string) []byte {

	Dial, _, _ := websocket.DefaultDialer.Dial(nodeUrl, nil)
	Dial.WriteMessage(websocket.TextMessage, []byte(message))
	typeMessage, jsonResponse, _ := Dial.ReadMessage()
	expectedTypeMessage := websocket.TextMessage
	if typeMessage != expectedTypeMessage {
		fmt.Printf("Expected: %v, got: %v \n", expectedTypeMessage, typeMessage)
	}
	return jsonResponse
}
func extractResult(jsonObject map[string]interface{}, attribute string, t *testing.T) interface{} {
	value, ok := jsonObject[attribute]
	if !ok {
		t.Errorf("Searched attribute %v not found", attribute)
	}
	return value
}
func compareResults(resultSlice []interface{}, t *testing.T) interface{} {
	for index, result := range resultSlice {
		for index2, result2 := range resultSlice {
			if index == index2 {
				continue
			}
			// if result == result2 {
			if reflect.DeepEqual(result, result2) {
				return result
			}
		}
	}
	fmt.Printf("No matches between responses: %v", resultSlice)
	return nil
}
func validateResult(result interface{}, t *testing.T, dataTypes []reflect.Type) {
	var typeOfValue = reflect.TypeOf(result)
	validationPass := false

	for _, dataType := range dataTypes {
		validationPass = validationPass || dataType == typeOfValue
	}
	if !validationPass {
		t.Errorf("Unespected type of response")
	}
}
func TestMain(m *testing.M) {
	//tests general setup
	urlSlice = append(urlSlice, primaryUrl)
	urlSlice = append(urlSlice, secondaryUrl)
	urlSlice = append(urlSlice, tertiaryUrl)
	os.Exit(m.Run())
}
func TestEthprotocolVersion(t *testing.T) {

	methodName := "eth_protocolVersion"
	params := ""
	messageToSend := jsonRequest
	messageToSend = strings.Replace(messageToSend, "[method]", methodName, -1)
	messageToSend = strings.Replace(messageToSend, "[params]", params, -1)

	resultSlice := make([]interface{}, 0)
	for _, url := range urlSlice {
		var jsonResponse = caller(messageToSend, url)
		var jsonObject map[string]interface{} = convertToObject(jsonResponse)
		result := extractResult(jsonObject, searchedAttribute, t)
		resultSlice = append(resultSlice, result)
	}
	var comparedResult = compareResults(resultSlice, t)
	var typesSlice = append(typesSlice, typeString)
	validateResult(comparedResult, t, typesSlice)
}
func TestNetListening(t *testing.T) {

	methodName := "net_listening"
	params := ""
	messageToSend := jsonRequest
	messageToSend = strings.Replace(messageToSend, "[method]", methodName, -1)
	messageToSend = strings.Replace(messageToSend, "[params]", params, -1)

	resultSlice := make([]interface{}, 0)
	for _, url := range urlSlice {
		var jsonResponse = caller(messageToSend, url)
		var jsonObject map[string]interface{} = convertToObject(jsonResponse)
		result := extractResult(jsonObject, searchedAttribute, t)
		resultSlice = append(resultSlice, result)
	}
	var comparedResult = compareResults(resultSlice, t)
	var typesSlice = append(typesSlice, typeBoolean)
	validateResult(comparedResult, t, typesSlice)
}
func TestEthGetBalance(t *testing.T) {

	methodName := "eth_getBalance"
	params := "\"0xa7719d2eD3849D3CD10991b91f1E8D9a2044eD45\""
	messageToSend := jsonRequest
	messageToSend = strings.Replace(messageToSend, "[method]", methodName, -1)
	messageToSend = strings.Replace(messageToSend, "[params]", params, -1)

	resultSlice := make([]interface{}, 0)
	for _, url := range urlSlice {
		var jsonResponse = caller(messageToSend, url)
		var jsonObject map[string]interface{} = convertToObject(jsonResponse)
		result := extractResult(jsonObject, searchedAttribute, t)
		resultSlice = append(resultSlice, result)
	}
	var comparedResult = compareResults(resultSlice, t)
	var typesSlice = append(typesSlice, typeString)
	validateResult(comparedResult, t, typesSlice)
}
func TestEthGetStorageAt(t *testing.T) {

	methodName := "eth_getStorageAt"
	params := "\"0xb451c6835515f8a08ecc4cbc5c5dcb238a48f7b4\",\"0x0\",\"latest\"" //this is the address of a test token we used
	messageToSend := jsonRequest
	messageToSend = strings.Replace(messageToSend, "[method]", methodName, -1)
	messageToSend = strings.Replace(messageToSend, "[params]", params, -1)

	resultSlice := make([]interface{}, 0)
	for _, url := range urlSlice {
		var jsonResponse = caller(messageToSend, url)
		var jsonObject map[string]interface{} = convertToObject(jsonResponse)
		result := extractResult(jsonObject, searchedAttribute, t)
		resultSlice = append(resultSlice, result)
	}
	var comparedResult = compareResults(resultSlice, t)
	var typesSlice = append(typesSlice, typeString)
	validateResult(comparedResult, t, typesSlice)
}
func TestEthGetTransactionCount(t *testing.T) {

	methodName := "eth_getTransactionCount"
	params := "\"0xa7719d2eD3849D3CD10991b91f1E8D9a2044eD45\",\"latest\""
	messageToSend := jsonRequest
	messageToSend = strings.Replace(messageToSend, "[method]", methodName, -1)
	messageToSend = strings.Replace(messageToSend, "[params]", params, -1)

	resultSlice := make([]interface{}, 0)
	for _, url := range urlSlice {
		var jsonResponse = caller(messageToSend, url)
		var jsonObject map[string]interface{} = convertToObject(jsonResponse)
		result := extractResult(jsonObject, searchedAttribute, t)
		resultSlice = append(resultSlice, result)
	}
	var comparedResult = compareResults(resultSlice, t)
	var typesSlice = append(typesSlice, typeString)
	validateResult(comparedResult, t, typesSlice)
}
func TestEthGetBlockTransactionCountByHash(t *testing.T) {

	methodName := "eth_getBlockTransactionCountByHash"
	params := "\"0xb65b2f91f066fdc47a71652d3f0aed4c95f6f2a82f028cc7d8e6cc8b2c6ec11f\""
	messageToSend := jsonRequest
	messageToSend = strings.Replace(messageToSend, "[method]", methodName, -1)
	messageToSend = strings.Replace(messageToSend, "[params]", params, -1)

	resultSlice := make([]interface{}, 0)
	for _, url := range urlSlice {
		var jsonResponse = caller(messageToSend, url)
		var jsonObject map[string]interface{} = convertToObject(jsonResponse)
		result := extractResult(jsonObject, searchedAttribute, t)
		resultSlice = append(resultSlice, result)
	}
	var comparedResult = compareResults(resultSlice, t)
	var typesSlice = append(typesSlice, typeString)
	validateResult(comparedResult, t, typesSlice)
}
func TestEthGetBlockTransactionCountByNumber(t *testing.T) {

	methodName := "eth_getBlockTransactionCountByNumber"
	params := "\"0x1ACB6E3\""
	messageToSend := jsonRequest
	messageToSend = strings.Replace(messageToSend, "[method]", methodName, -1)
	messageToSend = strings.Replace(messageToSend, "[params]", params, -1)

	resultSlice := make([]interface{}, 0)
	for _, url := range urlSlice {
		var jsonResponse = caller(messageToSend, url)
		var jsonObject map[string]interface{} = convertToObject(jsonResponse)
		result := extractResult(jsonObject, searchedAttribute, t)
		resultSlice = append(resultSlice, result)
	}
	var comparedResult = compareResults(resultSlice, t)
	var typesSlice = append(typesSlice, typeString)
	validateResult(comparedResult, t, typesSlice)
}
func TestEthGetUncleCountByBlockHash(t *testing.T) {

	methodName := "eth_getUncleCountByBlockHash"
	params := "\"0xb65b2f91f066fdc47a71652d3f0aed4c95f6f2a82f028cc7d8e6cc8b2c6ec11f\""
	messageToSend := jsonRequest
	messageToSend = strings.Replace(messageToSend, "[method]", methodName, -1)
	messageToSend = strings.Replace(messageToSend, "[params]", params, -1)

	resultSlice := make([]interface{}, 0)
	for _, url := range urlSlice {
		var jsonResponse = caller(messageToSend, url)
		var jsonObject map[string]interface{} = convertToObject(jsonResponse)
		result := extractResult(jsonObject, searchedAttribute, t)
		resultSlice = append(resultSlice, result)
	}
	var comparedResult = compareResults(resultSlice, t)
	var typesSlice = append(typesSlice, typeString)
	validateResult(comparedResult, t, typesSlice)
}
func TestEthGetUncleCountByBlockNumber(t *testing.T) {

	methodName := "eth_getUncleCountByBlockNumber"
	params := "\"0x1ACB6E3\""
	messageToSend := jsonRequest
	messageToSend = strings.Replace(messageToSend, "[method]", methodName, -1)
	messageToSend = strings.Replace(messageToSend, "[params]", params, -1)

	resultSlice := make([]interface{}, 0)
	for _, url := range urlSlice {
		var jsonResponse = caller(messageToSend, url)
		var jsonObject map[string]interface{} = convertToObject(jsonResponse)
		result := extractResult(jsonObject, searchedAttribute, t)
		resultSlice = append(resultSlice, result)
	}
	var comparedResult = compareResults(resultSlice, t)
	var typesSlice = append(typesSlice, typeString)
	validateResult(comparedResult, t, typesSlice)
}
func TestEthGetCode(t *testing.T) {

	methodName := "eth_getCode"
	params := "\"0xb451c6835515f8a08ecc4cbc5c5dcb238a48f7b4\",\"latest\"" //this is the address of a test token
	messageToSend := jsonRequest
	messageToSend = strings.Replace(messageToSend, "[method]", methodName, -1)
	messageToSend = strings.Replace(messageToSend, "[params]", params, -1)

	resultSlice := make([]interface{}, 0)
	for _, url := range urlSlice {
		var jsonResponse = caller(messageToSend, url)
		var jsonObject map[string]interface{} = convertToObject(jsonResponse)
		result := extractResult(jsonObject, searchedAttribute, t)
		resultSlice = append(resultSlice, result)
	}
	var comparedResult = compareResults(resultSlice, t)
	var typesSlice = append(typesSlice, typeString)
	validateResult(comparedResult, t, typesSlice)
}
func TestEthGetBlockByHash(t *testing.T) {

	methodName := "eth_getBlockByHash"
	params := "\"0xb65b2f91f066fdc47a71652d3f0aed4c95f6f2a82f028cc7d8e6cc8b2c6ec11f\",false"
	messageToSend := jsonRequest
	messageToSend = strings.Replace(messageToSend, "[method]", methodName, -1)
	messageToSend = strings.Replace(messageToSend, "[params]", params, -1)

	resultSlice := make([]interface{}, 0)
	for _, url := range urlSlice {
		var jsonResponse = caller(messageToSend, url)
		var jsonObject map[string]interface{} = convertToObject(jsonResponse)
		result := extractResult(jsonObject, searchedAttribute, t)
		resultSlice = append(resultSlice, result)
	}
	var comparedResult = compareResults(resultSlice, t)
	var typesSlice = append(typesSlice, typeObject)
	validateResult(comparedResult, t, typesSlice)
}
func TestEthGetBlockByNumber(t *testing.T) {

	methodName := "eth_getBlockByNumber"
	params := "\"0x1ACB6E3\",false"
	messageToSend := jsonRequest
	messageToSend = strings.Replace(messageToSend, "[method]", methodName, -1)
	messageToSend = strings.Replace(messageToSend, "[params]", params, -1)

	resultSlice := make([]interface{}, 0)
	for _, url := range urlSlice {
		var jsonResponse = caller(messageToSend, url)
		var jsonObject map[string]interface{} = convertToObject(jsonResponse)
		result := extractResult(jsonObject, searchedAttribute, t)
		resultSlice = append(resultSlice, result)
	}
	var comparedResult = compareResults(resultSlice, t)
	var typesSlice = append(typesSlice, typeObject)
	t.Log(comparedResult)
	validateResult(comparedResult, t, typesSlice)
}
func TestEthGetTransactionByBlockHashAndIndex(t *testing.T) {

	methodName := "eth_getTransactionByBlockHashAndIndex"
	params := "\"0x7cce2f931903be2731dc04bbd49d5b1e7f55972ce3fb2f3983d484f335940ab7\",\"0x1\""
	messageToSend := jsonRequest
	messageToSend = strings.Replace(messageToSend, "[method]", methodName, -1)
	messageToSend = strings.Replace(messageToSend, "[params]", params, -1)

	resultSlice := make([]interface{}, 0)
	for _, url := range urlSlice {
		var jsonResponse = caller(messageToSend, url)
		var jsonObject map[string]interface{} = convertToObject(jsonResponse)
		result := extractResult(jsonObject, searchedAttribute, t)
		resultSlice = append(resultSlice, result)
	}
	var comparedResult = compareResults(resultSlice, t)
	var typesSlice = append(typesSlice, typeObject)
	validateResult(comparedResult, t, typesSlice)
}
func TestEthGetTransactionByBlockNumberAndIndex(t *testing.T) {

	methodName := "eth_getTransactionByBlockNumberAndIndex"
	params := "\"0x1AD1703\",\"0x1\""
	messageToSend := jsonRequest
	messageToSend = strings.Replace(messageToSend, "[method]", methodName, -1)
	messageToSend = strings.Replace(messageToSend, "[params]", params, -1)

	resultSlice := make([]interface{}, 0)
	for _, url := range urlSlice {
		var jsonResponse = caller(messageToSend, url)
		var jsonObject map[string]interface{} = convertToObject(jsonResponse)
		result := extractResult(jsonObject, searchedAttribute, t)
		resultSlice = append(resultSlice, result)
	}
	var comparedResult = compareResults(resultSlice, t)
	var typesSlice = append(typesSlice, typeObject)
	validateResult(comparedResult, t, typesSlice)
}
func TestEthGetTransactionReceipt(t *testing.T) {

	methodName := "eth_getTransactionReceipt"
	params := "\"0xbe16f66b00cd4395646636b1a84d75b16ed8b0d4d055a4921566170a8f5bd1dd\""
	messageToSend := jsonRequest
	messageToSend = strings.Replace(messageToSend, "[method]", methodName, -1)
	messageToSend = strings.Replace(messageToSend, "[params]", params, -1)

	resultSlice := make([]interface{}, 0)
	for _, url := range urlSlice {
		var jsonResponse = caller(messageToSend, url)
		var jsonObject map[string]interface{} = convertToObject(jsonResponse)
		result := extractResult(jsonObject, searchedAttribute, t)
		resultSlice = append(resultSlice, result)
	}
	var comparedResult = compareResults(resultSlice, t)
	var typesSlice = append(typesSlice, typeObject)
	t.Log(comparedResult)
	validateResult(comparedResult, t, typesSlice)
}
func TestEthGetUncleByBlockHashAndIndex(t *testing.T) {
	//Devuelve Nil si el bloque no tiene Uncles, ver de hacer esa excepcion
	methodName := "eth_getUncleByBlockHashAndIndex"
	params := "\"0xa2163d7d18578e0995b1304003b857337eaa4534cbe64905c7bd45a744932f1f\",\"0x0\""
	messageToSend := jsonRequest
	messageToSend = strings.Replace(messageToSend, "[method]", methodName, -1)
	messageToSend = strings.Replace(messageToSend, "[params]", params, -1)

	resultSlice := make([]interface{}, 0)
	for _, url := range urlSlice {
		var jsonResponse = caller(messageToSend, url)
		var jsonObject map[string]interface{} = convertToObject(jsonResponse)
		result := extractResult(jsonObject, searchedAttribute, t)
		resultSlice = append(resultSlice, result)
	}
	var comparedResult = compareResults(resultSlice, t)
	var typesSlice = append(typesSlice, typeObject)
	t.Log(comparedResult)
	validateResult(comparedResult, t, typesSlice)
}
func TestEthGetUncleByBlockNumberAndIndex(t *testing.T) {
	//Devuelve Nil si el bloque no tiene Uncles, ver de hacer esa excepcion
	methodName := "eth_getUncleByBlockNumberAndIndex"
	params := "\"0x4D50E2\",\"0x0\""
	messageToSend := jsonRequest
	messageToSend = strings.Replace(messageToSend, "[method]", methodName, -1)
	messageToSend = strings.Replace(messageToSend, "[params]", params, -1)

	resultSlice := make([]interface{}, 0)
	for _, url := range urlSlice {
		var jsonResponse = caller(messageToSend, url)
		var jsonObject map[string]interface{} = convertToObject(jsonResponse)
		result := extractResult(jsonObject, searchedAttribute, t)
		resultSlice = append(resultSlice, result)
	}
	var comparedResult = compareResults(resultSlice, t)
	var typesSlice = append(typesSlice, typeObject)
	t.Log(comparedResult)
	validateResult(comparedResult, t, typesSlice)
}
