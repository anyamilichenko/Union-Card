package utils

import "log"

func HandleFatalError(error error) {
	if error != nil {
		log.Fatal(error)
	}
}
