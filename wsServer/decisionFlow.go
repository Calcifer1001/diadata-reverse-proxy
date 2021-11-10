package main

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
	firstResponseServices = append(firstResponseServices, "web3_sha3")
	firstResponseServices = append(firstResponseServices, "net_version")
	firstResponseServices = append(firstResponseServices, "net_peerCount")
	firstResponseServices = append(firstResponseServices, "eth_syncing")
	firstResponseServices = append(firstResponseServices, "eth_mining")
	firstResponseServices = append(firstResponseServices, "eth_hashrate")
	firstResponseServices = append(firstResponseServices, "eth_gasPrice")
	firstResponseServices = append(firstResponseServices, "eth_accounts")
	firstResponseServices = append(firstResponseServices, "eth_blockNumber")
	firstResponseServices = append(firstResponseServices, "eth_getTransactionByHash")
	firstResponseServices = append(firstResponseServices, "eth_submitWork")
	firstResponseServices = append(firstResponseServices, "eth_submitHashrate")
	// eth_subscribe is not a comparison method, because the response contains different hash for every rpc
	// This triggers subsequent responses of eth_subscription messages that should be compared and are treated differently
	firstResponseServices = append(firstResponseServices, "eth_subscribe")

	compareServices = append(compareServices, "net_listening")
	compareServices = append(compareServices, "eth_protocolVersion")
	compareServices = append(compareServices, "eth_getBalance")
	compareServices = append(compareServices, "eth_getStorageAt")
	compareServices = append(compareServices, "eth_getTransactionCount")
	compareServices = append(compareServices, "eth_getBlockTransactionCountByHash")
	compareServices = append(compareServices, "eth_getBlockTransactionCountByNumber")
	compareServices = append(compareServices, "eth_getUncleCountByBlockHash")
	compareServices = append(compareServices, "eth_getUncleCountByBlockNumber")
	compareServices = append(compareServices, "eth_getCode")
	compareServices = append(compareServices, "eth_getBlockByHash")
	compareServices = append(compareServices, "eth_getBlockByNumber")
	compareServices = append(compareServices, "eth_getTransactionByBlockHashAndIndex")
	compareServices = append(compareServices, "eth_getTransactionByBlockNumberAndIndex")
	compareServices = append(compareServices, "eth_getTransactionReceipt")
	compareServices = append(compareServices, "eth_getUncleByBlockHashAndIndex")
	compareServices = append(compareServices, "eth_getUncleByBlockNumberAndIndex")

	sendOnlyToPrimaryServices = append(sendOnlyToPrimaryServices, "eth_newFilter")
	sendOnlyToPrimaryServices = append(sendOnlyToPrimaryServices, "eth_newBlockFilter")
	sendOnlyToPrimaryServices = append(sendOnlyToPrimaryServices, "eth_newPendingTransactionFilter")
	sendOnlyToPrimaryServices = append(sendOnlyToPrimaryServices, "eth_uninstallFilter")
	sendOnlyToPrimaryServices = append(sendOnlyToPrimaryServices, "eth_getFilterChanges")
	sendOnlyToPrimaryServices = append(sendOnlyToPrimaryServices, "eth_getFilterLogs")
	sendOnlyToPrimaryServices = append(sendOnlyToPrimaryServices, "eth_getLogs")

	deprecatedServices = append(deprecatedServices, "eth_coinbase")
	deprecatedServices = append(deprecatedServices, "eth_sign")
	deprecatedServices = append(deprecatedServices, "eth_getCompilers")
	deprecatedServices = append(deprecatedServices, "eth_compileLLL")
	deprecatedServices = append(deprecatedServices, "eth_compileSolidity")
	deprecatedServices = append(deprecatedServices, "eth_compileSerpent")
	deprecatedServices = append(deprecatedServices, "eth_getWork")
	deprecatedServices = append(deprecatedServices, "db_putString")
	deprecatedServices = append(deprecatedServices, "db_getString")
	deprecatedServices = append(deprecatedServices, "db_putHex")
	deprecatedServices = append(deprecatedServices, "db_getHex")
	deprecatedServices = append(deprecatedServices, "shh_post")
	deprecatedServices = append(deprecatedServices, "shh_version")
	deprecatedServices = append(deprecatedServices, "shh_newIdentity")
	deprecatedServices = append(deprecatedServices, "shh_hasIdentity")
	deprecatedServices = append(deprecatedServices, "shh_newGroup")
	deprecatedServices = append(deprecatedServices, "shh_addToGroup")
	deprecatedServices = append(deprecatedServices, "shh_newFilter")
	deprecatedServices = append(deprecatedServices, "shh_uninstallFilter")
	deprecatedServices = append(deprecatedServices, "shh_getFilterChanges")
	deprecatedServices = append(deprecatedServices, "shh_getMessages")

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
// func getActionAndIdFromMessage(message []byte) (float64, int) {
// 	var result map[string]interface{}
// 	json.Unmarshal(message, &result)

// 	return result["id"].(float64), getAction(result["method"].(string))
// }

func arrayContains(array []string, element string) bool {
	for _, arrayElement := range array {
		if arrayElement == element {
			return true
		}
	}
	return false
}
