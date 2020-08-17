package util

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

func OpenFile() string {
	jsonFile, err := os.Open("config.json")
	if err != nil {
		fmt.Println(err)
		return "Could not open file"
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	var result map[string]interface{}
	json.Unmarshal([]byte(byteValue), &result)

	if str, ok := result["connectionString"].(string); ok {
		return str
	}

	return "Could not find db connection string"
}
