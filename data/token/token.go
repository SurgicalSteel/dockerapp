package token

import (
	"database/sql"
	"github.com/febytanzil/dockerapp/framework/database"
	"log"
	"time"
	"googlemaps.github.io/maps"
)

type Token interface {
	InsertToken(token string, origin maps.LatLng, destinations string) error
	GetToken(token string) (*TokenData, error)
}

type TokenData struct {
	ID         int64
	Token      string
	OriginLat float64
	OriginLong float64
	Destinations string
	CreateTime time.Time
}

type tokenImpl struct{}

var t Token

const (
	insertQuery = `insert into t_token (token, origin_lat, origin_long, destinations create_time) values ($1, $2, $3, $4, $5)`
	getQuery    = `select id, token, origin_lat, origin_long, destinations, create_time from t_token where token = $1`
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
	_, err := database.Get().Exec(insertQuery, token, origin.Lat, origin.Lng, destinations, time.Now())
	return err
}

func (i *tokenImpl) GetToken(token string) (*TokenData, error) {
	r := new(TokenData)
	err := database.Get().QueryRow(getQuery, token).Scan(&r.ID, &r.Token, &r.OriginLat, &r.OriginLong, &r.Destinations, &r.CreateTime)
	if nil != err {
		if sql.ErrNoRows == err {
			return nil, nil
		}
		return nil, err
	}
	return r, err
}
