package plugin_log4shell

import (
	"fmt"
)

// Token types.
const (
	Start     = "START"
	End       = "END"
	Content   = "CONTENT"
	Separator = "SEP"
)

// Node types.
const (
	Expression = "EXP"
	Text       = "TXT"
	Root       = "ROOT"
)

// Nodes a set of nodes.
type Nodes []*Node

func (e Nodes) String() string {
	var data string
	for _, v := range e {
		data += v.String()
	}
	return data
}

// Node a node.
type Node struct {
	Type string

	Text  string
	Key   Nodes
	Value Nodes
}

func (n Node) String() string {
	switch n.Type {
	case Expression, Root:
		return n.Value.String()
	case Text:
		return n.Text
	default:
		panic(fmt.Sprintf("not supported node type: %s", n.Type))
	}
}

// Token a syntax token.
type Token struct {
	Type  string `json:"type,omitempty"`
	Pos   int    `json:"pos,omitempty"`
	Value string `json:"value,omitempty"`
}

func (t Token) String() string {
	return t.Value
}

// Parse naively parses Log4j expression.
// https://logging.apache.org/log4j/2.x/manual/configuration.html#PropertySubstitution
func Parse(value string) *Node {
	root := &Node{Type: Root}

	buildTree(root, tokenize(value))

	return root
}

func tokenize(value string) []*Token {
	var tokens []*Token

	var previous *Token
	var open int
	length := len(value)

	for i := 0; i < length; i++ {
		v := value[i]
		t := &Token{Pos: i}

		switch {
		case v == '$' && length > i+1 && value[i+1] == '{':
			t.Type = Start
			t.Value = "${"
			i++
			open++

		case v == '}' && open > 0:
			t.Type = End
			t.Value = "}"
			open--

		case v == ':' && open > 0:
			t.Type = Separator
			t.Value = ":"

			if length > i+1 && value[i+1] == '-' {
				t.Value = ":-"
				i++
			}

		default:
			if previous != nil && previous.Type == Content {
				previous.Value += string(v)
				continue
			}

			t.Type = Content
			t.Value = string(v)
		}

		previous = t
		tokens = append(tokens, t)
	}

	return tokens
}

func buildTree(root *Node, tokens []*Token) int {
	var sep bool
	for i := 0; i < len(tokens); i++ {
		token := tokens[i]

		switch token.Type {
		case Start:
			exp := &Node{Type: Expression}

			switch root.Type {
			case Root:
				root.Value = append(root.Value, exp)

			case Expression:
				if sep {
					root.Key = append(root.Value, root.Key...)
					root.Value = []*Node{exp}
				} else {
					root.Key = append(root.Key, exp)
				}

			default:
				panic(fmt.Sprintf("invalid start node: %T", root))
			}

			j := buildTree(exp, tokens[i+1:])
			if j < 0 {
				return i
			}

			i += j

		case End:
			return i + 1

		case Content:
			switch root.Type {
			case Root:
				root.Value = append(root.Value, &Node{Type: Text, Text: token.Value})

			case Expression:
				if sep {
					root.Key = append(root.Value, root.Key...)
					root.Value = []*Node{{Type: Text, Text: token.Value}}
				} else {
					root.Key = append(root.Key, &Node{Type: Text, Text: token.Value})
				}
			default:
				panic(fmt.Sprintf("invalid content node: %T", root))
			}

		case Separator:
			sep = true
			continue

		default:
			panic(fmt.Sprintf("invalid token type: %s", token.Type))
		}
	}

	return -1
}
