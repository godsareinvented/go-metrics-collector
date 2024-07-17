package interfaces

type Storage interface {
	Get(key string) interface{}
	Set(key string, value interface{})
}
