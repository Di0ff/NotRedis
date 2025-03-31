package parser

import (
	"NotRedis/internal/myError"
	"go.uber.org/zap"
	"strings"
)

type MySplitRequest struct {
	Type  string
	Key   string
	Value string
}

type Parser struct {
	logger *zap.Logger
}

func NewParser(logger *zap.Logger) *Parser {
	return &Parser{
		logger: logger,
	}
}

func (p *Parser) Parse(request string) (MySplitRequest, error) {
	f := strings.Fields(request)
	if len(f) == 0 {
		p.logger.Error("Parse error")
		return MySplitRequest{}, myError.EmptyRequest
	}

	switch f[0] {
	case "SET":
		if len(f) != 3 {
			p.logger.Error("Parse fail: SET fail")
			return MySplitRequest{}, myError.SetFail
		}
		p.logger.Info("Parse success")
		return MySplitRequest{Type: f[0], Key: f[1], Value: f[2]}, nil
	case "GET":
		if len(f) != 2 {
			p.logger.Error("Parse fail: GET fail")
			return MySplitRequest{}, myError.GetFail
		}
		p.logger.Info("Parse success")
		return MySplitRequest{Type: f[0], Key: f[1]}, nil
	case "DEL":
		if len(f) != 2 {
			p.logger.Error("Parse fail: DEL fail")
			return MySplitRequest{}, myError.DelFail
		}
		p.logger.Info("Parse success")
		return MySplitRequest{Type: f[0], Key: f[1]}, nil
	default:
		p.logger.Error("Parse fail")
		return MySplitRequest{}, myError.UnknownRequest
	}
}
