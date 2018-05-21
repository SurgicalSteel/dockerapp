package api

import (
	"encoding/json"
	"github.com/febytanzil/dockerapp/services/route"
	"github.com/gorilla/mux"
	"github.com/tokopedia/transactionapp/core/log"
	maps2 "googlemaps.github.io/maps"
	"net/http"
	"net/url"
)

type GetResponse struct {
	Status   string         `json:"status"`
	Error    string         `json:"error,omitempty"`
	Distance int            `json:"total_distance,omitempty"`
	Duration float64        `json:"total_time,omitempty"`
	Path     []maps2.LatLng `json:"path,omitempty"`
}

func GetRoute(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	raw, ok := vars["token"]
	if !ok {
		res := &GetResponse{
			Error: "failed to read token",
		}
		resBytes, _ := json.Marshal(res)
		w.WriteHeader(http.StatusOK)
		w.Write(resBytes)
		return
	}
	token, err := url.PathUnescape(raw)
	if nil != err {
		log.Println("error decode token", err)
		res := &GetResponse{
			Error: "failed to decode token",
		}
		resBytes, _ := json.Marshal(res)
		w.WriteHeader(http.StatusOK)
		w.Write(resBytes)
		return
	}
	result, err := route.GetInstance().GetShortestRoute(token)
	if nil != err {
		res := &GetResponse{}
		switch err {
		case route.ErrProgress:
			res.Status = "in progress"
		case route.ErrCalculate:
			res.Status = "failure"
			res.Error = "error calculate path"
		default:
			res.Status = "failure"
			res.Error = err.Error()
		}
		resBytes, _ := json.Marshal(res)
		w.WriteHeader(http.StatusOK)
		w.Write(resBytes)
	} else {
		if nil == result {

		}
		res := &GetResponse{
			Status:   "OK",
			Distance: result.Distance,
			Duration: result.Time,
		}
		resBytes, _ := json.Marshal(res)
		w.WriteHeader(http.StatusOK)
		w.Write(resBytes)
	}
}
