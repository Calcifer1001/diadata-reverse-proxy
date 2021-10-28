package main

import "encoding/json"

var firstResponseServices []string = make([]string, 0)
var compareServices []string = make([]string, 0)
var sendOnlyOnceServices []string = make([]string, 0)
var sendOnlyToPrimaryServices []string = make([]string, 0)
var deprecatedServices []string = make([]string, 0)

var ACTION_SERVICE_DOESNT_EXISTS = -1
var ACTION_SEND_FIRST_RESPONSE = 0
var ACTION_COMPARE_SERVICES = 1
var ACTION_SEND_ONLY_ONCE = 2
var ACTION_SEND_ONLY_TO_PRIMARY = 3
var ACTION_DEPRECATED = 4

func initializeServiceLists() {
	firstResponseServices = append(firstResponseServices, "web3_clientVersion")
}

/**
* service is the string as written in the rpc server. For example eth_subscribe
 */
func getAction(service string) int {
	if arrayContains(firstResponseServices, service) {
		return ACTION_SEND_FIRST_RESPONSE
	}

	if arrayContains(compareServices, service) {
		return ACTION_COMPARE_SERVICES
	}

	if arrayContains(sendOnlyOnceServices, service) {
		return ACTION_SEND_ONLY_ONCE
	}

	if arrayContains(sendOnlyToPrimaryServices, service) {
		return ACTION_SEND_ONLY_TO_PRIMARY
	}

	if arrayContains(deprecatedServices, service) {
		return ACTION_DEPRECATED
	}

	return ACTION_SERVICE_DOESNT_EXISTS
}

/**
* Returns in the first position the id of the message and in the second the actionId acording to the constants defined on top
 */
func getActionAndIdFromMessage(message []byte) (int, int) {
	var result map[string]interface{}
	json.Unmarshal(message, &result)

	return result["id"].(int), getAction(result["method"].(string))
}

func arrayContains(array []string, element string) bool {
	for _, arrayElement := range array {
		if arrayElement == element {
			return true
		}
	}
	return false
}
