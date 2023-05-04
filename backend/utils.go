package main

import "encoding/json"

func GetResponseBody(success bool, data interface{}) string {

	res := map[string]interface{}{
		"success": success,
		"data":    data,
	}

	bytes, err := json.Marshal(res)
	if err != nil {
		return "Something went wrong"
	}
	return string(bytes)
}
