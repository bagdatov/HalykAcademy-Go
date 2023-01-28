package main

import (
	"log"
	"net/http"
	_ "net/http/pprof"
	api "web/api"
)

func main() {

	http.HandleFunc("/getone", api.GetOne)
	http.HandleFunc("/saveone", api.SaveOne)
	http.HandleFunc("/getall", api.GetAll)

	log.Println("server is listening on http://localhost:8080")

	log.Println(http.ListenAndServe(":8080", nil))
}
