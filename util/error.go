package util

import (
	"log"
)

// HandleError shows the error to console if DEBUG is enabled
func HandleError(msg string, err error) {
	if Mode == 0 {
		log.Println("[DEBUG]", msg)
		log.Println(err)
	}
}
