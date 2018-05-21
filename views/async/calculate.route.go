package async

import (
	"github.com/febytanzil/dockerapp/services/route"
	"log"
)

func CalculateRoute(token string) {
	log.Println("token calculate", token)
	err := route.GetInstance().CalculateRoute(token)
	if nil != err {
		log.Println("error calculate route async", token, err)
	}
}
