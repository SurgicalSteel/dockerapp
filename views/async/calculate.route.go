package async

import "github.com/febytanzil/dockerapp/services/route"

func CalculateRoute(token string) {
	route.GetInstance().CalculateRoute(token)
}
