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
}

func (t *tree) newList() *listNode {
	return &listNode{nodeType: nodeList}
}

func (l *listNode) append(n Node) {
	l.children = append(l.children, n)
	l.nextChild = &n
}

func (l *listNode) Next() Node {
	return *l.nextChild
}

type elementNode struct {
	nodeType
	data string
	attr map[string]string
}

func (e *elementNode) Data() string {
	return e.data
}

func (e *elementNode) Attr() map[string]string {
	return e.attr
}
