package main

import "log"

func main() {

	err := portUsb()
	if err != nil {
		log.Println(err)
	}

}
