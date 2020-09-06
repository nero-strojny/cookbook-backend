package util

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

// OpenFile find the connection db string
func OpenFile(filePath string) string {
	jsonFile, err := os.Open(filePath)
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
