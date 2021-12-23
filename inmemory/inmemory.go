package inmemory

type InMemoryDB interface {
	Set(key, value string) error
	Get(key string) (string, error)
}
