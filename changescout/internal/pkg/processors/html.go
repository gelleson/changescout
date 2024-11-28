package processors

import (
	"bytes"
	"github.com/PuerkitoBio/goquery"
	"github.com/gelleson/changescout/changescout/internal/domain"
)

type HTMLProcessor struct {
	conf domain.Setting
}

func NewHTMLProcessor(conf domain.Setting) *HTMLProcessor {
	return &HTMLProcessor{
		conf: conf,
	}
}

func (p *HTMLProcessor) Skip() bool {
	return p.conf.Selectors == nil
}

func (p *HTMLProcessor) Process(body []byte) []byte {
	if p.Skip() {
		return body
	}

	var stringBuilder bytes.Buffer

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
	if err != nil {
		return body
	}

	for _, selector := range p.conf.Selectors {
		doc.Find(selector).Each(func(i int, s *goquery.Selection) {
			stringBuilder.WriteString(s.Text())
		})
	}

	return stringBuilder.Bytes()
}
