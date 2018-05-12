package maps

import (
	"net/http"
	"time"
	"fmt"
	"log"
	"googlemaps.github.io/maps"
	"context"
)

type DMA interface {
	GetShortestPath(startLat, startLong float64, destinations []maps.LatLng) (Destinations, error)
}

type RouteResponse struct {
	Pos maps.LatLng
	Dist int
	Dur float64
}

type Destinations []RouteResponse

type dmaImpl struct {}

var (
	h = &http.Client{
		Timeout: time.Second * 5,
	}
	dma DMA
)

const (
	key = `AIzaSyCeSzsBN7NJ2wmIBQG7KEiHUCetwBXcbPA`
)

func InitDMA(d DMA) {
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

func (i *dmaImpl) GetShortestPath(startLat, startLong float64, destinations []maps.LatLng) (Destinations, error) {
	c, err := maps.NewClient(maps.WithAPIKey(key))
	if nil != err {
		return nil, err
	}

	res, err := c.DistanceMatrix(context.Background(), &maps.DistanceMatrixRequest{
		Origins:[]string{fmt.Sprintf("%f,%f", startLat, startLong)},
		Destinations:[]string{maps.Encode(destinations)},
	})
	if nil != err {
		return nil, err
	}
	d := Destinations{}
	for i, v := range destinations {
		o := RouteResponse{
			Pos: v,
			Dist: res.Rows[0].Elements[i].Distance.Meters,
			Dur:res.Rows[0].Elements[i].Duration.Seconds(),
		}
		d = append(d, o)
	}

	return d, nil
}
