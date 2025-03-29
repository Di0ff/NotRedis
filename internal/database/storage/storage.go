package storage

import (
	"NotRedis/internal/myError"
	"go.uber.org/zap"
)

type Engine interface {
	Set(key, value string) error
	Get(key string) (string, error)
	Del(key string) error
}

type Storage struct {
	engine Engine
	logger *zap.Logger
}

func NewStorage(engine Engine, logger *zap.Logger) (*Storage, error) {
	if engine == nil {
		return nil, myError.EngineNil
	}
	if logger == nil {
		return nil, myError.LoggerNil
	}

	storage := &Storage{
		engine: engine,
		logger: logger,
	}

	return storage, nil
}

func (s *Storage) Set(key, value string) error {
	err := s.engine.Set(key, value)
	if err != nil {
		s.logger.Error("Set command fail")
		return err
	}

	s.logger.Info("Set command success")
	return nil
}

func (s *Storage) Get(key string) (string, error) {
	value, err := s.engine.Get(key)
	if err != nil {
		s.logger.Error("Get command fail")
		return "", err
	}

	s.logger.Info("Get command success")
	return value, nil
}

func (s *Storage) Del(key string) error {
	err := s.engine.Del(key)
	if err != nil {
		s.logger.Error("Del command fail")
		return err
	}

	s.logger.Info("Del command success")
	return nil
}
