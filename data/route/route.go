package route

import (
	"database/sql"
	"github.com/febytanzil/dockerapp/framework/database"
	"log"
	"time"
)

type Route interface {
	InsertRoute(tx *sql.Tx, data *RouteData) error
	GetRouteByID(tokenID int64) (*RouteData, error)
}

type RouteData struct {
	TokenID    int64
	Result     string
	CreateTime time.Time
}

type routeImpl struct{}

var r Route

const (
	insertQuery = `insert into t_route (token_id, result, create_time) values ($1, $2, $3)`
	getQuery    = `select token_id, result, create_time from t_route where token_id = $1`
)

func Init(i Route) {
	if nil == r {
		r = &routeImpl{}
	} else {
		r = i
	}
}

func GetInstance() Route {
	if nil == r {
		log.Fatal("route is not initialized")
	}
	return r
}

func (i *routeImpl) InsertRoute(tx *sql.Tx, data *RouteData) error {
	_, err := tx.Exec(insertQuery, data.TokenID, data.Result, time.Now())
	if nil != err {
		log.Println("error insert route result", err)
	}
	return err
}

func (i *routeImpl) GetRouteByID(tokenID int64) (*RouteData, error) {
	r := new(RouteData)
	err := database.Get().QueryRow(getQuery, tokenID).Scan(&r.TokenID, &r.Result, &r.CreateTime)
	if nil != err {
		if sql.ErrNoRows == err {
			return nil, nil
		}
		return nil, err
	}
	return r, err
}
