package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"time"
)

type measurement struct {
	Time        time.Time
	Humidity    int `json:"humidity"`
	Temperature int `json:"temperature"`
}

func getDataArray() []measurement {
	DB.mut.Lock()
	defer DB.mut.Unlock()

	fmt.Println("1")
	var count int
	err := DB.db.QueryRow("SELECT COUNT(*) FROM sensor").Scan(&count)
	if err != nil {
		log.Println(err)
	}

	if count == 0 {
		return nil
	}

	fmt.Println("2")
	measures := make([]measurement, count)

	rows, err := DB.db.Query("SELECT * FROM sensor")

	i := 0
	for rows.Next() {
		var singleMeasure measurement
		err = rows.Scan(&singleMeasure.Time, &singleMeasure.Humidity, &singleMeasure.Temperature)
		fmt.Println(singleMeasure)
		if err != nil {
			log.Println(err)
		}
		measures[i] = singleMeasure
		i++
	}
	fmt.Println("3")

	return measures
}

func graph(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("graph.html"))

	// Create a slice of strings
	measureArray := getDataArray()

	type data struct {
		Time        []string
		Humidity    []int
		Temperature []int
	}

	newData := data{Time: make([]string, len(measureArray)), Humidity: make([]int, len(measureArray)), Temperature: make([]int, len(measureArray))}

	for i := 0; i < len(measureArray); i++ {
		newData.Time[i] = measureArray[i].Time.Format(time.TimeOnly)
		newData.Humidity[i] = measureArray[i].Humidity
		newData.Temperature[i] = measureArray[i].Temperature
	}

	// Execute the template, passing the slice as data
	err := tmpl.Execute(w, newData)
	if err != nil {
		log.Println("Error executing template:", err)
	}
}

func dbPage(w http.ResponseWriter, r *http.Request) {

	measureArray := getDataArray()

	tmpl, err := template.ParseFiles("DB.html")
	if err != nil {
		log.Println(err)
	}

	tmpl.Execute(w, measureArray)

}

func postHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
	}

	data := measurement{}

	err = json.Unmarshal(body, &data)
	if err != nil {
		log.Print(err)
	}

	err = DB.save_measurements(data.Humidity, data.Temperature)
	if err != nil {
		log.Println(err)
	}

	//err = data.db.insertData([2]int{data.Humidity, data.Temperature})
	//if err != nil {
	//	log.Print(err)
	//} else {
	//	fmt.Println("data printed")
	//}

	defer r.Body.Close()
	fmt.Println("ping")
	fmt.Println("printed", data.Humidity, data.Temperature)
}

func Handler() {
	http.HandleFunc("/graph", graph)
	http.HandleFunc("/dbPage", dbPage)
	http.HandleFunc("/post", postHandler)
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprint(writer, "Hello Client")
	})
}
