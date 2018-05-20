package token

import (
	"database/sql"
	"github.com/febytanzil/dockerapp/framework/database"
	"googlemaps.github.io/maps"
	"log"
	"time"
)

type Token interface {
	InsertToken(token string, origin maps.LatLng, destinations string) error
	GetToken(token string) (*TokenData, error)
	UpdateToken(token string, status int) error
}

type TokenData struct {
	ID           int64
	Token        string
	Status       int
	OriginLat    float64
	OriginLong   float64
	Destinations string
	CreateTime   time.Time
}

type tokenImpl struct{}

var t Token

const (
	insertQuery = `insert into t_token (token, status, origin_lat, origin_long, destinations, create_time) values ($1, $2, $3, $4, $5, $6)`
	getQuery    = `select id, token, status, origin_lat, origin_long, destinations, create_time from t_token where token = $1`
	updateQuery = `update t_token set status = $1 where token = $2`

	StatusPending = 1
	StatusSuccess = 2
	StatusError   = 3
)

func Init(i Token) {
	if nil == t {
		t = &tokenImpl{}
	} else {
		t = i
	}
}

func GetInstance() Token {
	if nil == t {
		log.Fatal("token is not initialized")
	}
	return t
}

func (i *tokenImpl) InsertToken(token string, origin maps.LatLng, destinations string) error {
	_, err := database.Get().Exec(insertQuery, token, StatusPending, origin.Lat, origin.Lng, destinations, time.Now())
	if nil != err {
		log.Println("error insert token", err)
	}
	return err
}

func (i *tokenImpl) GetToken(token string) (*TokenData, error) {
	r := new(TokenData)
	err := database.Get().QueryRow(getQuery, token).Scan(&r.ID, &r.Token, &r.Status, &r.OriginLat, &r.OriginLong, &r.Destinations, &r.CreateTime)
	if nil != err {
		if sql.ErrNoRows == err {
			return nil, nil
		}
		log.Println("error get token", err)
		return nil, err
	}
	return r, err
}

func (i *tokenImpl) UpdateToken(token string, status int) error {
	_, err := database.Get().Exec(updateQuery, status, token)
	if nil != err {
		log.Println("error update token", err)
	}
	return err
}
