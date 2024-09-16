package parser

// Node is a trivial interface implemented by any node type in the DOM tree
type Node interface {
	Type() nodeType
	Data() string
}

type nodeType int

func (t nodeType) Type() nodeType {
	return t
}

const (
	nodeList nodeType = iota
	nodeText
	nodeElement
	nodeComment
	nodeDocument
)

type listNode struct {
	nodeType
	children  []Node
	nextChild *Node
	pos       int
}

func (t *tree) newList() *listNode {
	return &listNode{nodeType: nodeList, pos: 0}
}

func (l *listNode) append(n Node) {
	l.children = append(l.children, n)
	l.nextChild = &n
	l.pos++
}

func (l *listNode) Next() Node {
	return *l.nextChild
}

func (l *listNode) lastChild() Node {
	return l.children[l.pos-1]
}

type documentElement struct {
	nodeType
	data string
}

func (d *documentElement) Data() string {
	return d.data
}

func (e *documentElement) Type() nodeType {
	return nodeDocument
}

type elementNode struct {
	nodeType
	data string
	attr map[string]string
}

func (e *elementNode) Data() string {
	return e.data
}

func (e *elementNode) Type() nodeType {
	return nodeElement
}

func (e *elementNode) Attr() map[string]string {
	return e.attr
}
