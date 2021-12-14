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

// Token a syntax token.
type Token struct {
	Type  string `json:"type,omitempty"`
	Pos   int    `json:"pos,omitempty"`
	Value string `json:"value,omitempty"`
}

func (t Token) String() string {
	return t.Value
}

// Element a set of nodes.
type Element []fmt.Stringer

func (e Element) String() string {
	var data string
	for _, v := range e {
		data += v.String()
	}
	return data
}

// NodeText a text node.
type NodeText struct {
	Text string
}

func (n NodeText) String() string {
	return n.Text
}

// NodeExpression an expression node.
type NodeExpression struct {
	Key   Element
	Value Element
}

func (n NodeExpression) String() string {
	return n.Value.String()
}

// NodeRoot a root node.
type NodeRoot struct {
	Values Element
}

func (n NodeRoot) String() string {
	return n.Values.String()
}

// Parse naively parses Log4j expression.
// https://logging.apache.org/log4j/2.x/manual/configuration.html#PropertySubstitution
func Parse(value string) *NodeRoot {
	root := &NodeRoot{}

	tree(root, tokenizer(value))

	return root
}

func tokenizer(value string) []*Token {
	var tokens []*Token

	var previous *Token

	for i := 0; i < len(value); i++ {
		v := value[i]
		t := &Token{Pos: i}

		switch {
		case v == '$' && value[i+1] == '{':
			t.Type = Start
			t.Value = "${"
			i++

		case v == '}':
			t.Type = End
			t.Value = "}"

		case v == ':':
			t.Type = Separator
			t.Value = ":"

			if value[i+1] == '-' {
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

func tree(root fmt.Stringer, tokens []*Token) int {
	var sep bool
	for i := 0; i < len(tokens); i++ {
		token := tokens[i]

		switch token.Type {
		case Start:
			exp := &NodeExpression{}

			if v, ok := root.(*NodeRoot); ok {
				v.Values = append(v.Values, exp)
			} else if v, ok := root.(*NodeExpression); ok {
				if sep {
					v.Key = append(v.Value, v.Key...)
					v.Value = []fmt.Stringer{exp}
				} else {
					v.Key = append(v.Key, exp)
				}
			} else {
				panic(fmt.Sprintf("invalid start node: %T", root))
			}

			j := tree(exp, tokens[i+1:])
			if j < 0 {
				return i
			}

			i += j

		case End:
			return i + 1

		case Content:

			if v, ok := root.(*NodeRoot); ok {
				v.Values = append(v.Values, &NodeText{Text: token.Value})
			} else if v, ok := root.(*NodeExpression); ok {
				if sep {
					v.Key = append(v.Value, v.Key...)
					v.Value = []fmt.Stringer{&NodeText{Text: token.Value}}
				} else {
					v.Key = append(v.Key, &NodeText{Text: token.Value})
				}
			} else {
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
