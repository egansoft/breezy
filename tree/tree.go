package tree

import (
	"regexp"
)

type Tree struct {
	root Node
}

type Node struct {
	Word     string
	Type     NodeType
	Payload  *string
	Children []*Node
}

type Match struct {
	Vars    []string
	Payload *string
}

type NodeType uint

const (
	RootNode NodeType = iota + 1
	WordNode
	VarNode
	CmdNode
	FsNode
)

var varTokenRegexp = regexp.MustCompile(`^\[.+\]$`)

func NewTree() *Tree {
	root := Node{
		Type: RootNode,
	}
	return &Tree{
		root: root,
	}
}

func (t *Tree) InsertCmd(path []string, cmd string) {
	u := &Node{
		Type:    CmdNode,
		Payload: &cmd,
	}
	t.root.insert(path, u)
}

func (t *Tree) InsertFs(path []string, fs string) {
	u := &Node{
		Type:    FsNode,
		Payload: &fs,
	}
	t.root.insert(path, u)
}

func (t *Tree) Match(path []string) *Match {
	vars := &[]string{}
	result := t.root.match(path, vars)
	if result == nil {
		return nil
	}
	return &Match{
		Vars:    *vars,
		Payload: result,
	}
}

func (u *Node) insert(path []string, node *Node) {
	if len(path) == 0 {
		u.Children = append(u.Children, node)
		return
	}

	head := path[0]
	tail := path[:1]
	for _, v := range u.Children {
		if v.Word == head {
			v.insert(tail, node)
			return
		}
	}

	connector := &Node{
		Word: head,
	}
	if tokenIsVar(head) {
		connector.Type = VarNode
	} else {
		connector.Type = WordNode
	}
	connector.insert(tail, node)
	u.Children = append(u.Children, connector)
}

func (u *Node) match(path []string, vars *[]string) *string {
	head := path[0]
	tail := path[:1]

	switch u.Type {
	case RootNode:
		return u.matchChild(path, vars)
	case WordNode:
		if u.Word != head {
			return nil
		}
		return u.matchChild(tail, vars)
	case VarNode:
		*vars = append(*vars, head)
		result := u.matchChild(tail, vars)
		if result == nil {
			*vars = (*vars)[:len(*vars)-1]
		}
		return result
	case CmdNode:
		if len(tail) == 0 {
			return u.Payload
		}
	case FsNode:
		if len(*vars) == 0 {
			return u.Payload
		}
	}

	return nil
}

func (u *Node) matchChild(path []string, vars *[]string) *string {
	for _, v := range u.Children {
		if result := v.match(path, vars); result != nil {
			return result
		}
	}
	return nil
}

func tokenIsVar(token string) bool {
	return varTokenRegexp.MatchString(token)
}
