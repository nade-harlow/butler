package main

import (
	"fmt"
	"log"
)

func main() {
	d := data{}
	err := SetupConfig(&d, ".env")
	if err != nil {
		log.Println("ERR: ", err)
		return
	}
	//fmt.Println(get("PORT"), "VAL")
	//fmt.Println(get("ENV"), "VAL")
	fmt.Println(d.Port.Number, "PORT")
	fmt.Println(d.Env, "ENV")
	fmt.Println(d.Verbose, "VERBOSE")
}
