package interfaces

type DependencyInjection interface {
	Get() interface{}
	Return(interface{})
}
