package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

// Open loads a config data from the specified file
func Open(filename string, jsonStruct interface{}) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}

	json.Unmarshal(data, &jsonStruct)
}
