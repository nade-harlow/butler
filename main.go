package main

import (
	"fmt"
	"log"
)

func main() {
	err := load(".env")
	if err != nil {
		log.Println("ERR: ", err)
		return
	}
	fmt.Println(get("PORT"), "VAL")
}
