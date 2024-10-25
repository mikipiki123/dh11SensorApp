package main

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"github.com/davecheney/i2c"
	"go.bug.st/serial"
	"log"
	"net"
	"net/http"
	"time"
)

const url = "http://MacBookPro.lan:8080/post"

type measurement struct {
	Humidity    int `json:"humidity"`
	Temperature int `json:"temperature"`
}

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

func portI2C() error {
	println("I2C port reader")
	bus, err := i2c.New(0x55, 1)
	if err != nil {
		return err
	}
	println("1")
	defer bus.Close()

	buf := make([]byte, 4)

	for {
		_, err := bus.Read(buf)
		if err != nil {
			log.Println(err)
		}
		humidity := int16(binary.LittleEndian.Uint16(buf[0:2]))
		temperature := int16(binary.LittleEndian.Uint16(buf[2:4]))

		data := measurement{Humidity: int(humidity), Temperature: int(temperature)}

		jsonBody, err := json.Marshal(data)
		if err != nil {
			log.Println(err)
		}

		fmt.Println("2")
		if addresses, err := net.LookupHost("MacBookPro.lan"); err != nil {
			log.Println(err)
			continue
		} else {
			println(addresses)
		}
		if r, err := http.Post(url, "application/json", bytes.NewBuffer(jsonBody)); err == nil {
			fmt.Println(r.Status)
		} else {
			log.Println(err)
			continue
		}

		println("-|-")
		fmt.Println("--- Humidity:", humidity, "% ---")
		fmt.Println("--- Temperature:", temperature, "% ---")
		println("-|-")
		time.Sleep(30 * time.Minute)
	}
}
