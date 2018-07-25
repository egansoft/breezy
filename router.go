package main

type RouteSpec struct {
	CmdRoutes []*CmdRoute
	FsRoutes  []*FsRoute
}

type CmdRoute struct {
}

type FsRoute struct {
	UrlRoot string
	FsRoot  string
}

func NewRouteSpec(routeFile string) *RouteSpec {
	return &RouteSpec{}
}
