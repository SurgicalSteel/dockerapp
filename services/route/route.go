package route

import (
	"googlemaps.github.io/maps"
	"log"
	token2 "github.com/febytanzil/dockerapp/data/token"
	"errors"
	"github.com/febytanzil/dockerapp/data/route"
	maps2 "github.com/febytanzil/dockerapp/data/maps"
	"sort"
	"encoding/json"
	"time"
)

type RouteSvc interface {
	SubmitRoute(origin maps.LatLng, destinations []maps.LatLng) (string, error)
	GetShortestRoute(token string) (*Result, error)
}

func Init(i RouteSvc) {
	if nil == r {
		r = &svcImpl{}
	} else {
		r = i
	}
}

func GetInstance() RouteSvc {
	if nil == r {
		log.Fatal("route svc is not initialized")
	}
	return r
}

type svcImpl struct{}
type Result struct {
	Path []maps.LatLng `json:"path"`
	Distance int `json:"distance"`
	Time float64 `json:"time"`
}
var (
	r RouteSvc
	ErrNoSubmit = errors.New("no data submit")
	ErrCalculate = errors.New("error to calculate paths")
)

func (i *svcImpl) SubmitRoute(origin maps.LatLng, destinations []maps.LatLng) (string, error) {
	token := i.generateToken(origin, destinations)
	exist, err := token2.GetInstance().GetToken(token)
	if nil != err {
		return "", err
	}
	if nil == exist {
		err = token2.GetInstance().InsertToken(token, origin, maps.Encode(destinations))
	}

	return token, err
}

func (i *svcImpl) GetShortestRoute(token string) (*Result, error) {
	exist, err := token2.GetInstance().GetToken(token)
	res := new(Result)
	if nil != err {
		return nil, err
	}
	if nil == exist {
		return nil, ErrNoSubmit
	}
	data, err := route.GetInstance().GetRouteByID(exist.ID)
	if nil != err {
		return nil, err
	}
	if nil == data {
		result, err := maps2.GetInstance().GetShortestPath(maps.LatLng{
			Lat: exist.OriginLat,
			Lng: exist.OriginLong,
		},exist.Destinations)
		if nil != err {
			return nil, ErrCalculate
		}
		sort.Sort(result)
		res.Distance = result.TotalDistance()
		res.Time = result.TotalTime()
		res.Path = result.Path()
		resStr, _ := json.Marshal(res)
		err = route.GetInstance().InsertRoute(&route.RouteData{
			Result: string(resStr),
			TokenID: exist.ID,
			CreateTime: time.Now(),
		})
		if nil != err {
			log.Println("error inserting result", err)
		}
	} else {
		err = json.Unmarshal([]byte(data.Result), res)
		if nil != err {
			return nil, err
		}
	}

	return res, nil
}

func (i *svcImpl) generateToken(origin maps.LatLng, destinations []maps.LatLng) string {
	return maps.Encode([]maps.LatLng{origin}) + maps.Encode(destinations)
}