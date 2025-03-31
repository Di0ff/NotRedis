package engine

import (
	"NotRedis/internal/myError"
	"go.uber.org/zap"
	"sync"
)

type MyMap struct {
	mu     sync.RWMutex
	data   map[string]string
	logger *zap.Logger
}

func NewMyMap(cap int, logger *zap.Logger) *MyMap {
	return &MyMap{
		data:   make(map[string]string, cap),
		logger: logger,
	}
}

func (om *MyMap) Set(key, value string) error {
	om.mu.Lock()
	defer om.mu.Unlock()

	if key == "" || value == "" {
		om.logger.Error("Set fail: empty key or value")

		return myError.EmptyKeyOrValue
	}

	om.data[key] = value
	om.logger.Info("Set success")

	return nil
}

func (om *MyMap) Get(key string) (string, error) {
	om.mu.RLock()
	defer om.mu.RUnlock()

	value, ok := om.data[key]
	if !ok {
		om.logger.Error("Get fail: key not found")

		return "", myError.KeyNotFound
	}

	om.logger.Info("Get success")

	return value, nil
}

func (om *MyMap) Delete(key string) error {
	om.mu.Lock()
	defer om.mu.Unlock()

	if _, ok := om.data[key]; !ok {
		om.logger.Error("Delete fail: key not found")

		return myError.KeyNotFound
	}

	delete(om.data, key)
	om.logger.Info("Delete success")

	return nil
}
