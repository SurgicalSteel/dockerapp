package api

import (
	"net/http"
	"github.com/gorilla/mux"
	"github.com/febytanzil/dockerapp/services/route"
)

func GetRoute(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	token, ok := vars["token"]
	if !ok {

	}
	route.GetInstance().GetShortestRoute(token)
}
