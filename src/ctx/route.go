package ctx

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc interface{}
	Middlewares []interface{}
}
