package route

type RouteSvc interface {
	SubmitRoute() error
	GetShortestRoute() (error)
}
