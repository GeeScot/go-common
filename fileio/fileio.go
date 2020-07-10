package fileio

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

// ReadJSON reads a file and unmarshals in to a struct
func ReadJSON(filename string, jsonStruct interface{}) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}

	json.Unmarshal(data, jsonStruct)
}
