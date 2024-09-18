package parser

// Node is a trivial interface implemented by any node type in the DOM tree
type Node interface {
	Type() string
	Data() string
}

type nodeType int

const (
	nodeList nodeType = iota
	nodeText
	nodeElement
	nodeComment
	nodeDocument
)

var nodeTypeMap = map[nodeType]string{
	nodeList:     "list",
	nodeText:     "text",
	nodeElement:  "element",
	nodeComment:  "comment",
	nodeDocument: "document",
}

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

type documentNode struct {
	data string
}

func (d *documentNode) Data() string {
	return d.data
}

func (d *documentNode) Type() string {
	return nodeTypeMap[nodeDocument]
}

type elementNode struct {
	data string
	attr map[string]string
}

func (e *elementNode) Data() string {
	return e.data
}

func (e *elementNode) Type() string {
	return nodeTypeMap[nodeElement]
}

func (e *elementNode) Attr() map[string]string {
	return e.attr
}

// TODO implement textNode and commentNode
