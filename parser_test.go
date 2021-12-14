package plugin_log4shell

import (
	"testing"
)

func Test_tokenizer(t *testing.T) {
	testCases := []struct {
		expression string
		expected   []*Token
	}{
		{
			expression: "${b:c}",
			expected: []*Token{
				{Type: "START", Pos: 0, Value: "${"},
				{Type: "CONTENT", Pos: 2, Value: "b"},
				{Type: "SEP", Pos: 3, Value: ":"},
				{Type: "CONTENT", Pos: 4, Value: "c"},
				{Type: "END", Pos: 5, Value: "}"},
			},
		},
		{
			expression: "a${b:c}d",
			expected: []*Token{
				{Type: "CONTENT", Pos: 0, Value: "a"},
				{Type: "START", Pos: 1, Value: "${"},
				{Type: "CONTENT", Pos: 3, Value: "b"},
				{Type: "SEP", Pos: 4, Value: ":"},
				{Type: "CONTENT", Pos: 5, Value: "c"},
				{Type: "END", Pos: 6, Value: "}"},
				{Type: "CONTENT", Pos: 7, Value: "d"},
			},
		},
		{
			expression: "a${b:c}d${e:f}g",
			expected: []*Token{
				{Type: "CONTENT", Pos: 0, Value: "a"},
				{Type: "START", Pos: 1, Value: "${"},
				{Type: "CONTENT", Pos: 3, Value: "b"},
				{Type: "SEP", Pos: 4, Value: ":"},
				{Type: "CONTENT", Pos: 5, Value: "c"},
				{Type: "END", Pos: 6, Value: "}"},
				{Type: "CONTENT", Pos: 7, Value: "d"},
				{Type: "START", Pos: 8, Value: "${"},
				{Type: "CONTENT", Pos: 10, Value: "e"},
				{Type: "SEP", Pos: 11, Value: ":"},
				{Type: "CONTENT", Pos: 12, Value: "f"},
				{Type: "END", Pos: 13, Value: "}"},
				{Type: "CONTENT", Pos: 14, Value: "g"},
			},
		},
		{
			expression: "a${b:c${e:f}g}d",
			expected: []*Token{
				{Type: "CONTENT", Pos: 0, Value: "a"},
				{Type: "START", Pos: 1, Value: "${"},
				{Type: "CONTENT", Pos: 3, Value: "b"},
				{Type: "SEP", Pos: 4, Value: ":"},
				{Type: "CONTENT", Pos: 5, Value: "c"},
				{Type: "START", Pos: 6, Value: "${"},
				{Type: "CONTENT", Pos: 8, Value: "e"},
				{Type: "SEP", Pos: 9, Value: ":"},
				{Type: "CONTENT", Pos: 10, Value: "f"},
				{Type: "END", Pos: 11, Value: "}"},
				{Type: "CONTENT", Pos: 12, Value: "g"},
				{Type: "END", Pos: 13, Value: "}"},
				{Type: "CONTENT", Pos: 14, Value: "d"},
			},
		},
		{
			expression: "a${b${e:f}g:c}d",
			expected: []*Token{
				{Type: "CONTENT", Pos: 0, Value: "a"},
				{Type: "START", Pos: 1, Value: "${"},
				{Type: "CONTENT", Pos: 3, Value: "b"},
				{Type: "START", Pos: 4, Value: "${"},
				{Type: "CONTENT", Pos: 6, Value: "e"},
				{Type: "SEP", Pos: 7, Value: ":"},
				{Type: "CONTENT", Pos: 8, Value: "f"},
				{Type: "END", Pos: 9, Value: "}"},
				{Type: "CONTENT", Pos: 10, Value: "g"},
				{Type: "SEP", Pos: 11, Value: ":"},
				{Type: "CONTENT", Pos: 12, Value: "c"},
				{Type: "END", Pos: 13, Value: "}"},
				{Type: "CONTENT", Pos: 14, Value: "d"},
			},
		},
		{
			expression: "q${::b${c:d}e}${z:y:-j}",
			expected: []*Token{
				{Type: "CONTENT", Pos: 0, Value: "q"},
				{Type: "START", Pos: 1, Value: "${"},
				{Type: "SEP", Pos: 3, Value: ":"},
				{Type: "SEP", Pos: 4, Value: ":"},
				{Type: "CONTENT", Pos: 5, Value: "b"},
				{Type: "START", Pos: 6, Value: "${"},
				{Type: "CONTENT", Pos: 8, Value: "c"},
				{Type: "SEP", Pos: 9, Value: ":"},
				{Type: "CONTENT", Pos: 10, Value: "d"},
				{Type: "END", Pos: 11, Value: "}"},
				{Type: "CONTENT", Pos: 12, Value: "e"},
				{Type: "END", Pos: 13, Value: "}"},
				{Type: "START", Pos: 14, Value: "${"},
				{Type: "CONTENT", Pos: 16, Value: "z"},
				{Type: "SEP", Pos: 17, Value: ":"},
				{Type: "CONTENT", Pos: 18, Value: "y"},
				{Type: "SEP", Pos: 19, Value: ":-"},
				{Type: "CONTENT", Pos: 21, Value: "j"},
				{Type: "END", Pos: 22, Value: "}"},
			},
		},
		{
			expression: "${b${e${g:h}:f}:c}",
			expected: []*Token{
				{Type: "START", Pos: 0, Value: "${"},
				{Type: "CONTENT", Pos: 2, Value: "b"},
				{Type: "START", Pos: 3, Value: "${"},
				{Type: "CONTENT", Pos: 5, Value: "e"},
				{Type: "START", Pos: 6, Value: "${"},
				{Type: "CONTENT", Pos: 8, Value: "g"},
				{Type: "SEP", Pos: 9, Value: ":"},
				{Type: "CONTENT", Pos: 10, Value: "h"},
				{Type: "END", Pos: 11, Value: "}"},
				{Type: "SEP", Pos: 12, Value: ":"},
				{Type: "CONTENT", Pos: 13, Value: "f"},
				{Type: "END", Pos: 14, Value: "}"},
				{Type: "SEP", Pos: 15, Value: ":"},
				{Type: "CONTENT", Pos: 16, Value: "c"},
				{Type: "END", Pos: 17, Value: "}"},
			},
		},
	}

	for _, test := range testCases {
		test := test
		t.Run(test.expression, func(t *testing.T) {
			tokens := tokenizer(test.expression)

			if !equalsTokens(test.expected, tokens) {
				t.Error("tokens parsing")

				t.Log("got")
				for _, token := range tokens {
					t.Logf("%#v\n", token)
				}

				t.Log("want")
				for _, token := range test.expected {
					t.Logf("%#v\n", token)
				}
			}
		})
	}
}

// equalsTokens custom comparison because reflect.DeepEqual doesn't work well with yaegi.
func equalsTokens(a, b []*Token) bool {
	if len(a) != len(b) {
		return false
	}

	for i, token := range a {
		if (token == nil && b[i] != nil) || (token != nil && b[i] == nil) {
			return false
		}

		if token == nil {
			continue
		}

		if token.Type != b[i].Type {
			return false
		}

		if token.Value != b[i].Value {
			return false
		}

		if token.Pos != b[i].Pos {
			return false
		}
	}

	return true
}
