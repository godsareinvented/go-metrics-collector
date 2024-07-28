package interfaces

type Storage interface {
	GetAll() map[string]interface{}
	Get(key string) interface{}
	Set(key string, value interface{})
}
