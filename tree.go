package parser

// tree represents an HTML document's DOM tree
type tree struct {
	root *node
}

// node represents a node in the DOM tree
type node struct {
	parent   *node
	children []*node
	element  *element
}

// element represents an HTML element
type element struct {
	tagName string
	attrs   map[string]string
}
