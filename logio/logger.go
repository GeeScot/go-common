package logio

import (
	"fmt"
	"time"
)

func Println(message interface{}, args ...interface{}) {
	var log interface{}

	switch message.(type) {
	case string:
		log = fmt.Sprintf(message.(string), args...)
		break
	default:
		log = message
		break
	}

	timeStamp := time.Now().Format(time.RFC3339)
	fmt.Printf("[%s] %v\n", timeStamp, log)
}
