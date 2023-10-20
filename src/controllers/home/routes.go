package home

import "github.com/ylfyt/meta/meta"

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
