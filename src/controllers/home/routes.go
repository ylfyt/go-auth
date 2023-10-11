package home

import (
	"go-auth/src/meta"
)

var Routes = []meta.EndPoint{
	{
		Method:      "GET",
		Path:        "/",
		HandlerFunc: home,
	},
	{
		Method:      "GET",
		Path:        "/ping",
		HandlerFunc: ping,
	},
}
