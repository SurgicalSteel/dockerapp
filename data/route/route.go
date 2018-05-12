package route

type Route interface {
	InsertRoute(token string) error
	GetRoute(token string) (*RouteData, error)
}

type RouteData struct {
	Token string
	Status int
	Result string
}

type routeImpl struct{}

func (i *routeImpl) InsertRoute(token string) error {

}

func (i *routeImpl) GetRoute(token string) (*RouteData, error) {

}