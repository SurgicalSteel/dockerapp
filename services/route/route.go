package route

import (
	"encoding/json"
	"errors"
	maps2 "github.com/febytanzil/dockerapp/data/maps"
	"github.com/febytanzil/dockerapp/data/route"
	token2 "github.com/febytanzil/dockerapp/data/token"
	"github.com/febytanzil/dockerapp/framework/database"
	"googlemaps.github.io/maps"
	"log"
	"time"
)

type RouteSvc interface {
	SubmitRoute(origin maps.LatLng, destinations []maps.LatLng) (string, error)
	GetShortestRoute(token string) (*maps2.Destinations, error)
	CalculateRoute(token string) error
}

func Init(i RouteSvc, ch chan string) {
	if nil == r {
		r = &svcImpl{
			c: ch,
		}
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

type svcImpl struct {
	c chan string
}

var (
	r            RouteSvc
	ErrNoSubmit  = errors.New("no data submitted")
	ErrCalculate = errors.New("error to calculate paths")
	ErrProgress  = errors.New("route data still in progress")
)

func (i *svcImpl) SubmitRoute(origin maps.LatLng, destinations []maps.LatLng) (string, error) {
	token := i.generateToken(origin, destinations)
	exist, err := token2.GetInstance().GetToken(token)
	if nil != err {
		return "", err
	}
	if nil == exist {
		err = token2.GetInstance().InsertToken(token, origin, maps.Encode(destinations))
	} else if token2.StatusError == exist.Status {
		err = token2.GetInstance().UpdateToken(nil, token, token2.StatusPending)
	}

	go func() {
		i.c <- token
	}()

	return token, err
}

func (i *svcImpl) GetShortestRoute(token string) (*maps2.Destinations, error) {
	exist, err := token2.GetInstance().GetToken(token)
	res := new(maps2.Destinations)
	if nil != err {
		return nil, err
	}
	if nil == exist {
		return nil, ErrNoSubmit
	}
	if token2.StatusPending == exist.Status {
		return nil, ErrProgress
	} else if token2.StatusError == exist.Status {
		return nil, ErrCalculate
	}
	data, err := route.GetInstance().GetRouteByID(exist.ID)
	if nil != err {
		return nil, err
	}
	if nil == data {
		return nil, nil
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

func (i *svcImpl) CalculateRoute(token string) error {
	exist, err := token2.GetInstance().GetToken(token)
	if nil != err {
		return err
	}
	if nil == exist {
		return errors.New("no submitted token exist")
	}
	if token2.StatusSuccess == exist.Status {
		log.Println("already calculate token", token)
		return nil
	}

	result, err := maps2.GetInstance().GetShortestPath(maps.LatLng{
		Lat: exist.OriginLat,
		Lng: exist.OriginLong,
	}, exist.Destinations)
	if nil != err {
		if err = token2.GetInstance().UpdateToken(nil, token, token2.StatusError); nil != err {
			log.Println("failed to set error for token", token)
		}
		return ErrCalculate
	}

	tx, err := database.Get().Begin()
	if nil != err {
		log.Println("failed to get tx", err)
	}
	defer tx.Rollback()

	resStr, _ := json.Marshal(result)
	err = route.GetInstance().InsertRoute(tx, &route.RouteData{
		Result:     string(resStr),
		TokenID:    exist.ID,
		CreateTime: time.Now(),
	})
	if nil != err {
		return err
	}

	err = token2.GetInstance().UpdateToken(tx, token, token2.StatusSuccess)
	if nil == err {
		err = tx.Commit()
		if nil != err {
			log.Println("failed to commit tx", err)
		}
	}

	return err
}
