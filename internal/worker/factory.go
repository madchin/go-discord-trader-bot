package worker

import "errors"

var ErrWorkerSpawnFactoryIsFull = errors.New("unable to spawn worker because factory is already full of workers")

type factoryWorkers struct {
	actualCount  int
	maximumCount int
}

func NewFactory(maximumCount int) *factoryWorkers {
	return &factoryWorkers{0, maximumCount}
}
