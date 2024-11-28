package diff

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"golang.org/x/net/html"
	"strings"
	"testing"
)

type DiffTestSuite struct {
	suite.Suite
	service *Service
}

func (s *DiffTestSuite) SetupTest() {
	s.service = NewDiffService()
}

func TestDiffSuite(t *testing.T) {
	suite.Run(t, new(DiffTestSuite))
}

func (s *DiffTestSuite) TestEmptyContent() {
	// Test when previous content is empty
	result, err := s.service.Compare([]byte{}, []byte("new content"))
	assert.NoError(s.T(), err)
	assert.True(s.T(), result.HasChanges)
	assert.Equal(s.T(), 100.0, result.ChangePercent)
	assert.Len(s.T(), result.Changes, 1)
	assert.Equal(s.T(), Modified, result.Changes[0].Type)
	assert.Equal(s.T(), "new content", result.Changes[0].Content)
	// Test when current content is empty
	result, err = s.service.Compare([]byte("old content"), []byte{})
	assert.NoError(s.T(), err)
	assert.True(s.T(), result.HasChanges)
	assert.Equal(s.T(), 100.0, result.ChangePercent)
	assert.Len(s.T(), result.Changes, 1)
}

func (s *DiffTestSuite) TestTextContent() {
	previous := []byte("Hello World")
	current := []byte("Hello Updated World")

	result, err := s.service.Compare(previous, current)
	assert.NoError(s.T(), err)
	assert.True(s.T(), result.HasChanges)
	assert.Greater(s.T(), result.ChangePercent, 0.0)
	assert.NotEmpty(s.T(), result.Changes)
}

func (s *DiffTestSuite) TestHTMLContent() {
	previous := []byte(`<!DOCTYPE html><html><body><h1>Hello</h1></body></html>`)
	current := []byte(`<!DOCTYPE html><html><body><h1>Hello Updated</h1></body></html>`)

	result, err := s.service.Compare(previous, current)
	assert.NoError(s.T(), err)
	assert.True(s.T(), result.HasChanges)
	assert.Greater(s.T(), result.ChangePercent, 0.0)
	assert.NotEmpty(s.T(), result.Changes)
}

func (s *DiffTestSuite) TestNoChanges() {
	content := []byte("Hello World")

	result, err := s.service.Compare(content, content)
	assert.NoError(s.T(), err)
	assert.False(s.T(), result.HasChanges)
	assert.Empty(s.T(), result.Changes)
}

func (s *DiffTestSuite) TestInvalidHTML() {
	previous := []byte(`<!DOCTYPE html><html><body><h1>Hello</h1></body></html>`)
	current := []byte(`<!DOCTYPE html><html><body><h1>Hello</h1><invalid></body></html>`)

	result, err := s.service.Compare(previous, current)
	assert.NoError(s.T(), err)
	assert.True(s.T(), result.HasChanges)
}

func (s *DiffTestSuite) TestAttributeComparison() {
	doc1 := &html.Node{
		Type: html.ElementNode,
		Attr: []html.Attribute{{Key: "class", Val: "test"}},
	}
	doc2 := &html.Node{
		Type: html.ElementNode,
		Attr: []html.Attribute{{Key: "class", Val: "test"}},
	}

	assert.True(s.T(), compareAttributes(doc1, doc2))

	doc2.Attr = []html.Attribute{{Key: "class", Val: "different"}}
	assert.False(s.T(), compareAttributes(doc1, doc2))
}

func (s *DiffTestSuite) TestHashesAreAlwaysSet() {
	// Test text comparison
	result, err := s.service.Compare([]byte("old text"), []byte("new text"))
	assert.NoError(s.T(), err)
	assert.NotEmpty(s.T(), result.PreviousHash)
	assert.NotEmpty(s.T(), result.CurrentHash)

	// Test HTML comparison
	previous := []byte(`<!DOCTYPE html><html><body>old</body></html>`)
	current := []byte(`<!DOCTYPE html><html><body>new</body></html>`)
	result, err = s.service.Compare(previous, current)
	assert.NoError(s.T(), err)
	assert.NotEmpty(s.T(), result.PreviousHash)
	assert.NotEmpty(s.T(), result.CurrentHash)

	// Test empty content
	result, err = s.service.Compare([]byte{}, []byte("new content"))
	assert.NoError(s.T(), err)
	assert.NotEmpty(s.T(), result.PreviousHash)
	assert.NotEmpty(s.T(), result.CurrentHash)
}

func BenchmarkDiffService(b *testing.B) {
	service := NewDiffService()

	b.Run("Text Comparison", func(b *testing.B) {
		previous := []byte("Hello World")
		current := []byte("Hello Updated World")

		b.ResetTimer()
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_, _ = service.Compare(previous, current)
		}
	})

	b.Run("HTML Comparison", func(b *testing.B) {
		previous := []byte(`<!DOCTYPE html><html><body><h1>Hello</h1></body></html>`)
		current := []byte(`<!DOCTYPE html><html><body><h1>Hello Updated</h1></body></html>`)

		b.ResetTimer()
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_, _ = service.Compare(previous, current)
		}
	})

	b.Run("Large HTML Comparison", func(b *testing.B) {
		previous := []byte(`<!DOCTYPE html><html><body>` + strings.Repeat("<div>Hello</div>", 1000) + `</body></html>`)
		current := []byte(`<!DOCTYPE html><html><body>` + strings.Repeat("<div>Hello Updated</div>", 1000) + `</body></html>`)

		b.ResetTimer()
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_, _ = service.Compare(previous, current)
		}
	})

	b.Run("No Changes", func(b *testing.B) {
		content := []byte("Hello World")

		b.ResetTimer()
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_, _ = service.Compare(content, content)
		}
	})
}
