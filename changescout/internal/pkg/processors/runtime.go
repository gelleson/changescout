package processors

type Processor interface {
	Skip() bool
	Process(body []byte) []byte
}

type ProcessRunner interface {
	Run(body []byte) []byte
}

type runtimeProcessor struct {
	processors []Processor
}

func New(processors ...Processor) ProcessRunner {
	return &runtimeProcessor{processors: processors}
}

func (p *runtimeProcessor) Run(body []byte) []byte {
	for _, processor := range p.processors {
		if !processor.Skip() {
			body = processor.Process(body)
		}
	}
	return body
}
