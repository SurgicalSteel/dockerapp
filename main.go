package main

import (
	"fmt"
	"github.com/febytanzil/dockerapp/framework/database"
	"github.com/febytanzil/dockerapp/data/maps"
	"github.com/gorilla/mux"
	"net/http"
	"github.com/febytanzil/dockerapp/views/api"
	"log"
)

func main() {
	fmt.Println("hellow")
	inject()

	r := mux.NewRouter()
	r.HandleFunc("/route", api.SubmitRoute).Methods(http.MethodPost)
	r.HandleFunc("/route/{token}", api.GetRoute).Methods(http.MethodGet)
	log.Fatal(http.ListenAndServe(":9000", r))
}

func inject()  {
	database.InitDB("")
	maps.InitDMA(nil)
}
