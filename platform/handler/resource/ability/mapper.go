package ability

import "net/http"

type RouteMapper interface {
	Map(ability string, request []byte) (*http.Request, error)
}

type DefaultRouMapper struct {
}

func (r *DefaultRouMapper) Map(ability string, request []byte) (*http.Request, error) {
	return nil, nil
}
