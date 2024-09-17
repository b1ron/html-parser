package parser

import (
	"strings"
	"testing"
)

var simple = `
<!DOCTYPE html>
  <html>
  </html>
`

func TestParser(t *testing.T) {
	p := newParser(strings.NewReader(simple))
	tr := p.parse()
	for _, n := range tr.root.children {
		t.Logf("node type: %s, node data: %s", n.String(), n.Data())
	}
}
