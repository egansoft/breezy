package tree

import (
	"bytes"
	"fmt"
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
	Vars     []string
	Residual []string
	Payload  *string
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

func New() *Tree {
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

func (u *Node) insert(path []string, node *Node) {
	if len(path) == 0 {
		u.Children = append(u.Children, node)
		return
	}

	head := path[0]
	tail := path[1:]
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

func (t *Tree) Match(path []string) *Match {
	vars := &[]string{}
	match := t.root.match(path, vars)
	if match == nil {
		return nil
	}
	match.Vars = *vars
	return match
}

func (u *Node) match(path []string, vars *[]string) *Match {
	switch u.Type {
	case RootNode:
		return u.matchChild(path, vars)

	case WordNode:
		if len(path) > 0 {
			head := path[0]
			tail := path[1:]

			if u.Word != head {
				return nil
			}
			return u.matchChild(tail, vars)
		}

	case VarNode:
		if len(path) > 0 {
			head := path[0]
			tail := path[1:]

			*vars = append(*vars, head)
			result := u.matchChild(tail, vars)
			if result == nil {
				*vars = (*vars)[:len(*vars)-1]
			}
			return result
		}

	case CmdNode:
		if len(path) == 0 {
			return &Match{
				Payload: u.Payload,
			}
		}

	case FsNode:
		if len(*vars) == 0 {
			return &Match{
				Payload:  u.Payload,
				Residual: path,
			}
		}
	}

	return nil
}

func (u *Node) matchChild(path []string, vars *[]string) *Match {
	for _, v := range u.Children {
		if result := v.match(path, vars); result != nil {
			return result
		}
	}
	return nil
}

type dfsItem struct {
	node  *Node
	level int
}

func (t *Tree) String() string {
	starter := &dfsItem{
		node:  &t.root,
		level: 0,
	}
	stack := []*dfsItem{starter}
	var buffer bytes.Buffer

	for len(stack) > 0 {
		item := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		for i := 0; i < item.level; i++ {
			buffer.WriteString("\t")
		}
		buffer.WriteString(item.node.String())
		buffer.WriteString("\n")

		for _, child := range item.node.Children {
			next := &dfsItem{
				node:  child,
				level: item.level + 1,
			}
			stack = append(stack, next)
		}
	}
	return buffer.String()
}

func (u *Node) String() string {
	switch u.Type {
	case RootNode:
		return "/"
	case WordNode:
		return fmt.Sprintf("/%s", u.Word)
	case VarNode:
		return fmt.Sprintf("/[%s]", u.Word)
	case CmdNode:
		return fmt.Sprintf("$ %s", *u.Payload)
	case FsNode:
		return fmt.Sprintf(": %s", *u.Payload)
	}
	return ""
}

func (m *Match) String() string {
	if m == nil {
		return "no match"
	}
	return fmt.Sprintf("matched!\npayload=\"%s\"\nvars=%v\nresidual=%v", *m.Payload, m.Vars, m.Residual)
}

func tokenIsVar(token string) bool {
	return varTokenRegexp.MatchString(token)
}
