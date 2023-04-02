package dummylock

type DummyLock struct {
}

func New() *DummyLock {
	return &DummyLock{}
}

func (dl *DummyLock) Lock() {
}

func (dl *DummyLock) Unlock() {
}
