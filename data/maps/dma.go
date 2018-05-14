package maps

import (
	"context"
	"fmt"
	"googlemaps.github.io/maps"
	"log"
	"net/http"
	"time"
)

type DMA interface {
	GetShortestPath(origin maps.LatLng, destinations string) (Destinations, error)
}

type RouteResponse struct {
	Pos  maps.LatLng
	Dist int
	Dur  float64
}

type Destinations []RouteResponse

func (s Destinations) Len() int {
	return len(s)
}
func (s Destinations) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s Destinations) Less(i, j int) bool {
	return s[i].Dist < s[j].Dist
}
func (s Destinations) TotalDistance() int {
	total := 0
	for _, v := range s {
		total += v.Dist
	}
	return total
}
func (s Destinations) TotalTime() float64 {
	total := float64(0)
	for _, v := range s {
		total += v.Dur
	}
	return total
}
func (s Destinations) Path() []maps.LatLng {
	var res []maps.LatLng
	for _, v := range s {
		res = append(res, v.Pos)
	}
	return res
}

type dmaImpl struct{}

var (
	h = &http.Client{
		Timeout: time.Second * 5,
	}
	dma DMA
)

const (
	key = `AIzaSyCeSzsBN7NJ2wmIBQG7KEiHUCetwBXcbPA`
)

func Init(d DMA) {
	if nil == d {
		dma = &dmaImpl{}
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

func (i *dmaImpl) GetShortestPath(origin maps.LatLng, destinations string) (Destinations, error) {
	c, err := maps.NewClient(maps.WithAPIKey(key))
	if nil != err {
		return nil, err
	}

	res, err := c.DistanceMatrix(context.Background(), &maps.DistanceMatrixRequest{
		Origins:      []string{fmt.Sprintf("%f,%f", origin.Lat, origin.Lng)},
		Destinations: []string{destinations},
	})
	if nil != err {
		return nil, err
	}
	d := Destinations{}
	for i, v := range maps.DecodePolyline(destinations) {
		o := RouteResponse{
			Pos:  v,
			Dist: res.Rows[0].Elements[i].Distance.Meters,
			Dur:  res.Rows[0].Elements[i].Duration.Seconds(),
		}
		d = append(d, o)
	}

	return d, nil
}
