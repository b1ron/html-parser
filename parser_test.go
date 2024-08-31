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
	p := NewParser(strings.NewReader(simple))
	p.parse()
}
