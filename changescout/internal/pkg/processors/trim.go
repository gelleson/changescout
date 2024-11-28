package processors

import (
	"bytes"
	"github.com/gelleson/changescout/changescout/internal/domain"
)

type TrimProcessor struct {
	conf domain.Setting
}

func NewTrimProcessor(conf domain.Setting) *TrimProcessor {
	return &TrimProcessor{
		conf: conf,
	}
}

func (p *TrimProcessor) Skip() bool {
	return !p.conf.Trim
}

func (p *TrimProcessor) Process(body []byte) []byte {
	if p.Skip() {
		return body
	}

	trimmed := bytes.TrimSpace(body)
	if trimmed == nil {
		return []byte{}
	}
	return trimmed
}
