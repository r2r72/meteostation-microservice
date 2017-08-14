package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

//Structure for getting JSON data before insert into DB
type meteo_data struct {
	Humidity      float32
	Temperature   float32
	TempByFeeling float32
	Pressure      int
	PPM           int
}

// Handle HTTP requests
func handleRequests() {

	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/", homePage)
	myRouter.HandleFunc("/meteo/post/hour", postHourData)

	log.Fatal(http.ListenAndServe(":88888", myRouter))
}

// Handling of requests to root
func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Server side micro service for Meteo Project based on Go")
}

// Getting JSON file trogh POST and saving data into DB
func postHourData(w http.ResponseWriter, r *http.Request) {
	var md meteo_data

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&md)
	checkErr(err)

	// Check if POST data return data above 0 (during initialization of sensors some data can be lost and return 0)
	if md.PPM <= 0 || md.Humidity == 0 || md.Pressure == 0 || md.Temperature == 0 || md.TempByFeeling == 0 {
		return
	}

	db, err := sql.Open("mysql", "meteo:meteo22data@/(localhost:3306)MeteoDB?charset=utf8")
	checkErr(err)

	defer db.Close()

	ins, err := db.Prepare("INSERT meteo SET humidity =?, temperature=?, tempbyfeeling=?, pressure=?, ppm=?")
	checkErr(err)

	res, err := ins.Exec(md.Humidity, md.Temperature, md.TempByFeeling, md.Pressure, md.PPM)
	checkErr(err)

	id, err := res.LastInsertId()
	checkErr(err)

	fmt.Println(id)

	fmt.Fprintln(w, "Post data here")
}

// Check error
func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	fmt.Println("Meteo project v.2.1")
	handleRequests()
}
