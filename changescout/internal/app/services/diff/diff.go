package diff

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/pmezard/go-difflib/difflib"
	"github.com/sergi/go-diff/diffmatchpatch"
	"golang.org/x/net/html"
	"strings"
	"sync"
	"time"
	"unsafe"
)

type Result struct {
	HasChanges    bool      `json:"has_changes"`
	Changes       []Change  `json:"changes"`
	ChangedAt     time.Time `json:"changed_at"`
	PreviousHash  string    `json:"previous_hash"`
	CurrentHash   string    `json:"current_hash"`
	ChangePercent float64   `json:"change_percent"`
	Diff          string    `json:"diff"`
}

type Change struct {
	Type    ChangeType `json:"type"`
	Content string     `json:"content"`
	Path    string     `json:"path"`
}

type ChangeType string

const (
	Added    ChangeType = "added"
	Removed  ChangeType = "removed"
	Modified ChangeType = "modified"
)

type Service struct {
	dmp         *diffmatchpatch.DiffMatchPatch
	bufferPool  sync.Pool
	changesPool sync.Pool
}

func NewDiffService() *Service {
	return &Service{
		dmp: diffmatchpatch.New(),
		bufferPool: sync.Pool{
			New: func() interface{} {
				return new(bytes.Buffer)
			},
		},
		changesPool: sync.Pool{
			New: func() interface{} {
				return make([]Change, 0, 8)
			},
		},
	}
}

func (s *Service) Compare(previous, current []byte) (Result, error) {
	result := Result{
		ChangedAt:    time.Now(),
		PreviousHash: hash(previous),
		CurrentHash:  hash(current),
	}

	if len(previous) == 0 || len(current) == 0 {
		result.HasChanges = true
		result.ChangePercent = 100
		changes := s.getChanges()
		changes = append(changes[:0], Change{
			Type:    Modified,
			Content: unsafeByteToString(current),
		})
		result.Changes = changes
		s.putChanges(changes)
		result.Diff, _ = difflib.GetUnifiedDiffString(difflib.UnifiedDiff{
			A:        difflib.SplitLines(string(previous)),
			B:        difflib.SplitLines(string(current)),
			FromFile: "Expected",
			FromDate: "",
			ToFile:   "Actual",
			ToDate:   "",
			Context:  100,
		})
		return result, nil
	}

	isHTMLContent := len(previous) > 9 && (bytes.Equal(bytes.ToLower(previous[:9]), []byte("<!doctype")) ||
		bytes.Contains(previous[:min(100, len(previous))], []byte("<html>")))

	if isHTMLContent {
		return s.compareHTML(previous, current)
	}
	return s.compareText(previous, current)
}

func (s *Service) compareText(previous, current []byte) (Result, error) {
	diffs := s.dmp.DiffMain(unsafeByteToString(previous), unsafeByteToString(current), true)
	s.dmp.DiffCleanupSemantic(diffs)

	result := Result{
		ChangedAt:    time.Now(),
		PreviousHash: hash(previous),
		CurrentHash:  hash(current),
	}

	changes := s.getChanges()
	totalLen := float64(len(previous) + len(current))
	changedLen := 0.0

	var buf strings.Builder
	for _, diff := range diffs {
		lines := strings.Split(diff.Text, "\n")
		for _, line := range lines {
			if line == "" {
				continue
			}
			switch diff.Type {
			case diffmatchpatch.DiffInsert:
				buf.WriteString(line)
				changes = append(changes, Change{
					Type:    Added,
					Content: diff.Text,
				})
				changedLen += float64(len(diff.Text))
			case diffmatchpatch.DiffDelete:
				buf.WriteString(line)
				changes = append(changes, Change{
					Type:    Removed,
					Content: diff.Text,
				})
				changedLen += float64(len(diff.Text))
			case diffmatchpatch.DiffEqual:
				buf.WriteString(line)
			}
		}
	}

	result.HasChanges = len(changes) > 0
	result.Changes = changes
	if totalLen > 0 {
		result.ChangePercent = (changedLen / totalLen) * 100
	}

	result.Diff = buf.String()
	result.Diff, _ = difflib.GetUnifiedDiffString(difflib.UnifiedDiff{
		A:        difflib.SplitLines(string(previous)),
		B:        difflib.SplitLines(string(current)),
		FromFile: "Expected",
		FromDate: "",
		ToFile:   "Actual",
		ToDate:   "",
		Context:  100,
	})

	return result, nil
}

func (s *Service) compareHTML(previous, current []byte) (Result, error) {
	prevDoc, err := html.Parse(bytes.NewReader(previous))
	if err != nil {
		return Result{}, err
	}

	currDoc, err := html.Parse(bytes.NewReader(current))
	if err != nil {
		return Result{}, err
	}

	result := Result{
		ChangedAt:    time.Now(),
		PreviousHash: hash(previous),
		CurrentHash:  hash(current),
	}

	changes := s.getChanges()
	changes = s.compareNodes("", prevDoc, currDoc, changes)
	result.HasChanges = len(changes) > 0
	result.Changes = changes

	totalNodes := countNodes(prevDoc) + countNodes(currDoc)
	if totalNodes > 0 {
		result.ChangePercent = (float64(len(changes)) / float64(totalNodes)) * 100
	}

	var buf strings.Builder
	s.renderHTMLDiff(&buf, "", prevDoc, currDoc)
	result.Diff = buf.String()
	result.Diff, _ = difflib.GetUnifiedDiffString(difflib.UnifiedDiff{
		A:        difflib.SplitLines(string(previous)),
		B:        difflib.SplitLines(string(current)),
		FromFile: "Expected",
		FromDate: "",
		ToFile:   "Actual",
		ToDate:   "",
		Context:  100,
	})

	return result, nil
}

func (s *Service) renderHTMLDiff(buf *strings.Builder, indent string, node1, node2 *html.Node) {
	if node1 == nil && node2 == nil {
		return
	}

	if node1 == nil {
		buf.WriteString(indent)
		buf.WriteString(s.formatNode(node2))
		buf.WriteString("\n")
		for child := node2.FirstChild; child != nil; child = child.NextSibling {
			s.renderHTMLDiff(buf, indent+"  ", nil, child)
		}
		return
	}

	if node2 == nil {
		buf.WriteString("- ")
		buf.WriteString(indent)
		buf.WriteString(s.formatNode(node1))
		buf.WriteString("\n")
		for child := node1.FirstChild; child != nil; child = child.NextSibling {
			s.renderHTMLDiff(buf, indent+"  ", child, nil)
		}
		return
	}

	if node1.Type != node2.Type || !compareAttributes(node1, node2) {
		buf.WriteString("- ")
		buf.WriteString(indent)
		buf.WriteString(s.formatNode(node1))
		buf.WriteString("\n")
		buf.WriteString("+ ")
		buf.WriteString(indent)
		buf.WriteString(s.formatNode(node2))
		buf.WriteString("\n")
	} else if node1.Type == html.TextNode && strings.TrimSpace(node1.Data) != strings.TrimSpace(node2.Data) {
		buf.WriteString("- ")
		buf.WriteString(indent)
		buf.WriteString(strings.TrimSpace(node1.Data))
		buf.WriteString("\n")
		buf.WriteString("+ ")
		buf.WriteString(indent)
		buf.WriteString(strings.TrimSpace(node2.Data))
		buf.WriteString("\n")
	} else {
		buf.WriteString("  ")
		buf.WriteString(indent)
		buf.WriteString(s.formatNode(node1))
		buf.WriteString("\n")
	}

	child1 := node1.FirstChild
	child2 := node2.FirstChild
	for child1 != nil || child2 != nil {
		s.renderHTMLDiff(buf, indent+"  ", child1, child2)
		if child1 != nil {
			child1 = child1.NextSibling
		}
		if child2 != nil {
			child2 = child2.NextSibling
		}
	}
}

func (s *Service) formatNode(n *html.Node) string {
	switch n.Type {
	case html.ElementNode:
		var attrs []string
		for _, attr := range n.Attr {
			attrs = append(attrs, fmt.Sprintf(`%s="%s"`, attr.Key, attr.Val))
		}
		if len(attrs) > 0 {
			return fmt.Sprintf("<%s %s>", n.Data, strings.Join(attrs, " "))
		}
		return fmt.Sprintf("<%s>", n.Data)
	case html.TextNode:
		return strings.TrimSpace(n.Data)
	default:
		return n.Data
	}
}

func (s *Service) compareNodes(path string, node1, node2 *html.Node, changes []Change) []Change {
	if node1 == nil && node2 == nil {
		return changes
	}

	if node1 == nil {
		return append(changes, Change{
			Type:    Added,
			Content: s.nodeToString(node2),
			Path:    path,
		})
	}

	if node2 == nil {
		return append(changes, Change{
			Type:    Removed,
			Content: s.nodeToString(node1),
			Path:    path,
		})
	}

	if node1.Type != node2.Type || !compareAttributes(node1, node2) {
		return append(changes, Change{
			Type:    Modified,
			Content: s.nodeToString(node2),
			Path:    path,
		})
	}

	if node1.Type == html.TextNode && strings.TrimSpace(node1.Data) != strings.TrimSpace(node2.Data) {
		return append(changes, Change{
			Type:    Modified,
			Content: strings.TrimSpace(node2.Data),
			Path:    path,
		})
	}

	child1 := node1.FirstChild
	child2 := node2.FirstChild

	for child1 != nil || child2 != nil {
		childPath := path
		if node1.Type == html.ElementNode {
			childPath += "/" + node1.Data
		}

		changes = s.compareNodes(childPath, child1, child2, changes)

		if child1 != nil {
			child1 = child1.NextSibling
		}
		if child2 != nil {
			child2 = child2.NextSibling
		}
	}

	return changes
}

func compareAttributes(node1, node2 *html.Node) bool {
	if len(node1.Attr) != len(node2.Attr) {
		return false
	}

	attrs1 := make(map[string]string, len(node1.Attr))
	for _, attr := range node1.Attr {
		attrs1[attr.Key] = attr.Val
	}

	for _, attr := range node2.Attr {
		if val, ok := attrs1[attr.Key]; !ok || val != attr.Val {
			return false
		}
	}

	return true
}

func (s *Service) nodeToString(n *html.Node) string {
	buf := s.getBuf()
	defer s.putBuf(buf)
	html.Render(buf, n)
	return buf.String()
}

func countNodes(n *html.Node) int {
	count := 1
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		count += countNodes(c)
	}
	return count
}

func hash(content []byte) string {
	h := sha256.New()
	h.Write(content)
	return hex.EncodeToString(h.Sum(nil))
}

func (s *Service) getChanges() []Change {
	return s.changesPool.Get().([]Change)[:0]
}

func (s *Service) putChanges(changes []Change) {
	s.changesPool.Put(changes)
}

func (s *Service) getBuf() *bytes.Buffer {
	buf := s.bufferPool.Get().(*bytes.Buffer)
	buf.Reset()
	return buf
}

func (s *Service) putBuf(buf *bytes.Buffer) {
	s.bufferPool.Put(buf)
}

func unsafeByteToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
