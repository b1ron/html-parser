package parser

import (
	"bufio"
	"strings"
	"testing"
)

var simple = `
<!DOCTYPE html>
  <html>
  </html>
`

func TestParser(t *testing.T) {
	b := bufio.NewReader(strings.NewReader(simple))
	p := newParser(*b)
	p.parse()
}
