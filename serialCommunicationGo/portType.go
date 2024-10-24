package main

import (
	"fmt"
	"go.bug.st/serial"
	"log"
	"time"
)

func portUsb() error {
	println("USB port reader")
	// Open serial port
	port, err := serial.Open("/dev/ttyACM0", &serial.Mode{
		BaudRate: 9600,
	})
	if err != nil {
		return err
	}
	defer port.Close()

	buffer := make([]byte, 100)
	for {
		// Read from the serial port
		n, err := port.Read(buffer)
		if err != nil {
			log.Fatal(err)
		}
		if n > 0 {
			fmt.Printf("Received: %s\n", string(buffer[:n]))
		}
		time.Sleep(time.Second)
	}
}
