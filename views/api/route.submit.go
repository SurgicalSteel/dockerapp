package api

import (
	"encoding/json"
	"github.com/febytanzil/dockerapp/services/route"
	"googlemaps.github.io/maps"
	"io/ioutil"
	"net/http"
	"strconv"
)

type SubmitResponse struct {
	Token string `json:"token,omitempty"`
	Error string `json:"error,omitempty"`
}

func SubmitRoute(w http.ResponseWriter, r *http.Request) {
	var (
		req          [][]string
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
			origin.Lat, err = strconv.ParseFloat(v[0], 64)
			origin.Lng, err = strconv.ParseFloat(v[1], 64)
		} else {
			destinations = append(destinations, origin)
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
