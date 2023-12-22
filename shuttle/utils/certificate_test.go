package utils

import (
	"fmt"
	"testing"

	"google.dev/google/shuttle/utils/log"
)

func TestGenKeyPair(t *testing.T) {
	cert, key, err := GenKeyPair()
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(string(cert))
	fmt.Println(string(key))
}
