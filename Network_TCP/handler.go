package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

//testHandler: for test the connectivity
func testHandler(w http.ResponseWriter, r *http.Request) {

	fmt.Println("Receive a test request at ", time.Now().Format(time.ANSIC))

	w.Header().Set("Access-Control-Allow-Origin", "*")             //可以所有domain跨域访问
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type") //
	w.Header().Set("Content-Type", "text/plain")

	if r.Method == "OPTIONS" {
		fmt.Println("OPTIONS request detected!")
		return
	}

	err := testConnectToMongoDB()

	if err != nil {
		http.Error(w, "Can't connect to DB! ", http.StatusInternalServerError)
		return
	}

	fmt.Println("GET request detected!")
	responseString := fmt.Sprintf("Test successfully! Current time: %v.", time.Now().Format(time.ANSIC))

	w.Write([]byte(responseString))
}

//allocationhandler: assign the job id for the collection
func allocationHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Receive a ID allocation request at ", time.Now().Format(time.ANSIC))

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		fmt.Println("useOrNotHandler OPTIONS request detected!")
		return
	}

	machineid, err := strconv.Atoi(mux.Vars(r)["machine_id"]) //locate machine ID

	if err != nil {
		fmt.Println("Wrong machine ID format?")
		http.Error(w, "cant recognized machine id!", http.StatusBadRequest)
		return
	}

	empty_collection, err := getEmptyCollectionFromMongoDB(machineid)

	if err != nil {
		fmt.Println("can't get collection arr")
		http.Error(w, "cant get collection arr!", http.StatusInternalServerError)
		return
	}

	w.Write([]byte(fmt.Sprint(empty_collection)))
}

//groupUploadHandler: for group upload data option
func groupUploadHandler(w http.ResponseWriter, r *http.Request) {

	fmt.Println("Receive a group upload request at ", time.Now().Format(time.ANSIC))

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")

	if r.Method == "OPTIONS" {
		fmt.Println("upload group item OPTIONS request detected!")
		return
	}

	machine_id, err := strconv.Atoi(mux.Vars(r)["machine_id"]) //locate machine ID

	if err != nil {
		fmt.Println("error machine_id to int")
		http.Error(w, "machine id is invalid or not existed!", http.StatusBadRequest)
		return
	}

	fmt.Println(machine_id)

	job_id, err := strconv.Atoi(mux.Vars(r)["job_id"]) //locate job ID

	if err != nil {
		fmt.Println("error job_id to int")
		http.Error(w, "job id is invalid or not existed!", http.StatusBadRequest)
		return
	}

	var machine_array Machine_group
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&machine_array)

	//fmt.Println(machine_array)

	if err != nil {
		log.Panic(err)
		http.Error(w, "the requestBody is not machine groups", http.StatusBadRequest)
		return
	}

	err = SaveManyToMongoDB(machine_array.Group_data, machine_id, job_id)

	if err != nil {
		log.Panic(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	responseString := fmt.Sprintf("Upload group successfully! Current time: %v.", time.Now().Format(time.ANSIC))

	w.Write([]byte(responseString)) //return the last V and pos
}

// retrievalSpecificHandler: query the data sequence
func retrievalSpecificHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")

	if r.Method == "OPTIONS" {
		fmt.Println("retrieve OPTIONS request detected!")
		return
	}

	//check all parameters
	id, err := strconv.Atoi(mux.Vars(r)["machine_id"]) //machine ID

	if err != nil {
		log.Panic(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	job_id, err := strconv.Atoi(mux.Vars(r)["job_id"]) //jon_id

	if err != nil {
		log.Panic(err)
		http.Error(w, "Can't find job ID!", http.StatusBadRequest)
		return
	}

	start, err := strconv.Atoi(r.URL.Query().Get("start")) //?start={start}

	if err != nil {
		log.Panic(err)
		http.Error(w, "Can't find head!", http.StatusBadRequest)
		return
	}

	limit, err := strconv.Atoi(r.URL.Query().Get("limit")) //?limit={limit}

	if err != nil {
		log.Panic(err)
		http.Error(w, "Can't find limit!", http.StatusBadRequest)
		return
	}

	fmt.Println("id: ", id)
	fmt.Println("job_id: ", job_id)
	fmt.Println("start: ", start)
	fmt.Println("limit: ", limit)
	fmt.Println("observer: ", false)
	fmt.Println()

	var result []Machine
	result, err = getDataFromMongoDB(id, job_id, start, limit)

	if err != nil {
		log.Panic(err)
		return
	}

	json_array, err := json.Marshal(&result)

	if err != nil {
		log.Panic(err)
		return
	}
	fmt.Println("Convert successfully!")
	w.Write(json_array)

}
