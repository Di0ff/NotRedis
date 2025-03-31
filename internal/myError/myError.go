package myError

import "errors"

var (
	KeyNotFound     = errors.New("key not found")
	EmptyKeyOrValue = errors.New("empty key or value")
	EngineNil       = errors.New("engine is nil")
	LoggerNil       = errors.New("logger is nil")
	EmptyRequest    = errors.New("empty request")
	SetFail         = errors.New("SET fail")
	GetFail         = errors.New("GET fail")
	DelFail         = errors.New("DEL fail")
	UnknownRequest  = errors.New("unknown request")
)
