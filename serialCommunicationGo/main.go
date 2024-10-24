package main

import "log"

func main() {

	err := portI2C()
	if err != nil {
		log.Println(err)
	}

}
