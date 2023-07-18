package store

type IDB interface {
	Get(key string) (string, bool)
	Set(key, value string) bool
	Delete(key string) bool
}
