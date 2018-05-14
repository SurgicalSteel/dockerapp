package api

import (
	"encoding/json"
	"github.com/febytanzil/dockerapp/services/route"
	"googlemaps.github.io/maps"
	"io/ioutil"
	"net/http"
)

type SubmitResponse struct {
	Token string `json:"token,omitempty"`
	Error string `json:"error,omitempty"`
}

type GetResponse struct {
	Status   string  `json:"status"`
	Error    string  `json:"error,omitempty"`
	Distance int     `json:"total_distance,omitempty"`
	Duration float64 `json:"total_time,omitempty"`
}

func SubmitRoute(w http.ResponseWriter, r *http.Request) {
	var (
		req          [][]float64
		origin       maps.LatLng
		destinations []maps.LatLng
	)
	body, err := ioutil.ReadAll(r.Body)
	if nil != err {
		res := &SubmitResponse{
			Error: "failed to read request body",
		}
		resBytes, _ := json.Marshal(res)
		w.WriteHeader(http.StatusOK)
		w.Write(resBytes)
		return
	}
	err = json.Unmarshal(body, &req)
	if nil != err {
		res := &SubmitResponse{
			Error: "failed to map request body",
		}
		resBytes, _ := json.Marshal(res)
		w.WriteHeader(http.StatusOK)
		w.Write(resBytes)
		return
	}

	for i, v := range req {
		if 0 == i {
			origin.Lat = v[0]
			origin.Lng = v[1]
		} else {
			destinations = append(destinations, maps.LatLng{
				Lat: v[0],
				Lng: v[1],
			})
		}
	}

	token, err := route.GetInstance().SubmitRoute(origin, destinations)
	if nil != err {
		res := &SubmitResponse{
			Error: "failed to save route",
		}
		resBytes, _ := json.Marshal(res)
		w.WriteHeader(http.StatusOK)
		w.Write(resBytes)
	} else {
		res := &SubmitResponse{
			Token: token,
		}
		resBytes, _ := json.Marshal(res)
		w.WriteHeader(http.StatusOK)
		w.Write(resBytes)
	}
}
