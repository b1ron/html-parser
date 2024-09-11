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
