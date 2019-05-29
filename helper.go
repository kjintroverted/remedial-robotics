package main

import "fmt"

func errCheck(message string, err error) {
	if err != nil {
		fmt.Println("ERROR", message, err)
	}
}

func log(messages ...interface{}) {
	fmt.Println(messages...)
}
