package maps

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/febytanzil/dockerapp/framework/config"
	"googlemaps.github.io/maps"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type DMA interface {
	GetShortestPath(origin maps.LatLng, destinations string) (*Destinations, error)
}

type RouteResponse struct {
	Pos  maps.LatLng
	Dist int
	Dur  float64
}

type Destinations struct {
	Path          []maps.LatLng `json:"path"`
	TotalDistance int           `json:"total_distance"`
	TotalTime     float64       `json:"total_time"`
}

// commonResponse contains the common response fields to most API calls inside
// the Google Maps APIs. This is used internally.
type commonResponse struct {
	// Status contains the status of the request, and may contain debugging
	// information to help you track down why the call failed.
	Status string `json:"status"`

	// ErrorMessage is the explanatory field added when Status is an error.
	ErrorMessage string `json:"error_message"`
}

// StatusError returns an error iff this object has a non-OK Status.
func (c *commonResponse) StatusError() error {
	if c.Status != "OK" {
		return fmt.Errorf("maps: %s - %s", c.Status, c.ErrorMessage)
	}
	return nil
}

type response struct {
	Routes            []maps.Route            `json:"routes"`
	GeocodedWaypoints []maps.GeocodedWaypoint `json:"geocoded_waypoints"`
	commonResponse
}

type dmaImpl struct{}
type dmaImpl2 struct{}

var (
	h = &http.Client{
		Timeout: time.Second * 5,
	}
	dma DMA
	key string
)

func Init(d DMA) {
	if nil == d {
		dma = &dmaImpl2{}
		key = config.Get().Maps.Key
	} else {
		dma = d
	}
}

func GetInstance() DMA {
	if nil == dma {
		log.Fatal("db is not initialized")
	}
	return dma
}

func (i *dmaImpl) GetShortestPath(origin maps.LatLng, destinations string) (*Destinations, error) {
	c, err := maps.NewClient(maps.WithAPIKey(key))
	if nil != err {
		log.Println("error get gmaps client", err)
		return nil, err
	}
	destArr, err := maps.DecodePolyline(destinations)
	if nil != err {
		log.Println("error decode polyline", err)
		return nil, err
	}
	wpArr := make([]string, len(destArr))
	for _, v := range destArr {
		wpArr = append(wpArr, v.String())
	}

	routes, _, err := c.Directions(context.Background(), &maps.DirectionsRequest{
		Origin:      origin.String(),
		Destination: origin.String(),
		Waypoints:   wpArr,
		Optimize:    true,
	})
	if nil != err {
		log.Println("error calculate distance", err)
		return nil, err
	}
	if 0 == len(routes) {
		return nil, errors.New("cannot find route")
	}

	route := routes[0]
	d := new(Destinations)
	dests, err := maps.DecodePolyline(destinations)
	if nil != err {
		log.Println("error decode polyline", err)
		return nil, err
	}
	d.Path = append(d.Path, origin)
	for i, v := range route.WaypointOrder {
		d.Path = append(d.Path, dests[v])
		d.TotalDistance += route.Legs[i].Meters
		d.TotalTime += route.Legs[i].Duration.Seconds()
	}

	return d, nil
}

func (i *dmaImpl2) GetShortestPath(origin maps.LatLng, destinations string) (*Destinations, error) {
	req, err := http.NewRequest(http.MethodGet, "https://maps.googleapis.com/maps/api/directions/json", nil)
	if nil != err {
		log.Println("error http request", err)
		return nil, err
	}
	q := req.URL.Query()
	q.Add("key", key)
	q.Add("origin", origin.String())
	q.Add("destination", origin.String())
	q.Add("waypoints", fmt.Sprintf("optimize:true|enc:%s:", destinations))
	req.URL.RawQuery = q.Encode()

	res, err := h.Do(req)
	if nil != err {
		log.Println("error http do", err)
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if nil != err {
		log.Println("error read body", err)
		return nil, err
	}
	result := new(response)
	err = json.Unmarshal(body, result)
	if nil != err {
		log.Println("error unmarshal body", err)
		return nil, err
	}
	if 0 == len(result.Routes) {
		return nil, errors.New("cannot find route")
	}

	route := result.Routes[0]
	d := new(Destinations)
	dests, err := maps.DecodePolyline(destinations)
	if nil != err {
		log.Println("error decode polyline", err)
		return nil, err
	}
	d.Path = append(d.Path, origin)
	for i, v := range route.WaypointOrder {
		d.Path = append(d.Path, dests[v])
		d.TotalDistance += route.Legs[i].Meters
		d.TotalTime += route.Legs[i].Duration.Seconds()
	}

	return d, nil
}
