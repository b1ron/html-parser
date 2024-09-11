package parser

// Node is a trivial interface implemented by any node type in the DOM tree
type Node interface {
	Type() NodeType
}

type NodeType int

func (t NodeType) Type() NodeType {
	return t
}

const (
	nodeList NodeType = iota
	nodeText
	nodeElement
	nodeComment
	nodeDocument
)

type ListNode struct {
	NodeType
	children []Node
	tr       *tree
}

func (t *tree) newList() *ListNode {
	return &ListNode{tr: t, NodeType: nodeList}
}

func (l *ListNode) append(n Node) {
	l.children = append(l.children, n)
}

type textNode struct {
	NodeType
	text string
}

func (t *textNode) Type() NodeType {
	return nodeText
}

type elementNode struct {
	NodeType
	tagName string
	attr    map[string]string
}

func (e *elementNode) Type() NodeType {
	return nodeElement
}

type commentNode struct {
	NodeType
	data string
}

func (c *commentNode) Type() NodeType {
	return nodeComment
}

type documentNode struct {
	NodeType
	docType string
}

func (d *documentNode) Type() NodeType {
	return nodeDocument
}
