package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	fmt.Println("Hello, world!")
	r := mux.NewRouter()

	//check the connection for debugging
	//http://localhost:8080/test
	r.Handle("/test", http.HandlerFunc(testHandler)).Methods("OPTIONS", "GET")

	//make sure using the correct jobID!
	//Ex: http://localhost:8080/getNewJobId/5
	r.Handle("/getNewJobId/{machine_id}", http.HandlerFunc(allocationHandler)).Methods("OPTIONS", "GET")

	//multiple upload
	//Ex: http://localhost:8080/groupUpload/5/1
	r.Handle("/groupUpload/{machine_id}/{job_id}", http.HandlerFunc(groupUploadHandler)).Methods("OPTIONS", "POST")

	//Get data from
	//Ex: http://localhost:8080/getGroupRange/5/1/?start=0&limit=30
	r.Handle("/getGroupRange/{machine_id}/{job_id}/", http.HandlerFunc(retrievalSpecificHandler)).Methods("OPTIONS", "GET")

	//single upload
	//r.Handle("/upload", http.HandlerFunc(uploadHandler)).Methods("OPTIONS", "POST")

	log.Fatal(http.ListenAndServe("192.168.1.10:8080", r)) //change this to fit to another computer!
}
