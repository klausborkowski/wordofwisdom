package cache

type Cache interface {
	Add(int, int64) error
	Get(int) (bool, error)
	Delete(int)
}
