package utils

import (
	"encoding/json"
	"fmt"
	"log"
)

func PrintObj(obj interface{}) {
	indent, err := json.MarshalIndent(obj, " ", "  ")
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(string(indent))
}
