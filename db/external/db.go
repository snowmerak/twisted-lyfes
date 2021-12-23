package external

type KVCache interface {
	Set(key, value string) error
	Get(key string) (string, error)
}
