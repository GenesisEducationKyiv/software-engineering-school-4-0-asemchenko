package util

import (
	"io"
	"log"
)

func CloseBodyWithErrorHandling(Body io.ReadCloser) {
	err := Body.Close()
	if err != nil {
		log.Fatal(err)
	}
}

func WarnOnError(err error) {
	if err != nil {
		log.Println(err)
	}
}
