/*
 * Copyright (c) 2022, Dana Burkart <dana.burkart@gmail.com>
 *
 * SPDX-License-Identifier: BSD-2-Clause
 */

package parser

import "strings"

type stack []jsxTag

func (s *stack) Empty() bool {
	return len(*s) == 0
}
func (s *stack) Push(j jsxTag) {
	*s = append(*s, j)
}

func (s *stack) Pop() (j jsxTag, empty bool) {
	if s.Empty() {
		j = jsxTag{}
		empty = true
	} else {
		i := len(*s) - 1
		j = (*s)[i]
		*s = (*s)[:i]
		empty = false
	}
	return
}

func (s *stack) Peek() (j jsxTag, empty bool) {
	if s.Empty() {
		j = jsxTag{}
		empty = true
	} else {
		j = (*s)[len(*s)-1]
		empty = false
	}
	return
}

// Really, we should be building an AST since Javadoc can have parameters
// virtually anywhere, but storing token lists in Blocks is simpler for now.
type Text []Token

func (t *Text) Length() int {
	return len(*t)
}

// Given a Text token list, return a string with all the parameters
// evaluated.
func (t *Text) Interpolate(doc *Document, symbols SymbolMap, flowIndent string) string {
	interpolationArray := make([]string, t.Length())
	jsxStack := stack{}

	for i := 0; i < t.Length(); i++ {
		token := (*t)[i]

		switch token.Type {
		case TOK_JDOC_NL:
			interpolationArray[i] = "\n" + flowIndent
		case TOK_JDOC_PARAM:
			str := ""
			if token.Lexeme == "@code" {
				str += "`"
				str += strings.TrimSpace((*t)[i+1].Lexeme)
				str += "`"
				i++
			}

			if token.Lexeme == "@link" {
				target := strings.TrimSpace((*t)[i+1].Lexeme)

				// Handle links local to the current class
				if target[0] == '#' {
					target = doc.Blocks[0].Name + target
				}

				symbol := symbols[target]
				if symbol.Type == SYM_TYPE_INVALID {
					str = "*" + target + "*"
				} else {
					// TODO: The name of the link should be a proper definition
					str += "[" + symbol.Name + "](" + symbol.Location + ")"
				}
				i++
			}

			interpolationArray[i] = str
		case TOK_JSX_O:
			jsxStack.Push(jsxTag{i, token.Lexeme})
			interpolationArray[i] = token.Lexeme
		case TOK_JSX_X:
			current := jsxTag{i, token.Lexeme}

			for {
				next, empty := jsxStack.Pop()

				if empty {
					break
				}

				if next.tagType() == current.tagType() {
					if next.tagType() == "pre" {
						interpolationArray[next.index] = "```java"
						token.Lexeme = "```"
					}
					break
				}

				// Close the tag
				interpolationArray[next.index] = next.close()
			}

			interpolationArray[i] = token.Lexeme
		default:
			interpolationArray[i] = token.Lexeme
		}
	}

	// Close anything still on the stack
	for {
		next, empty := jsxStack.Pop()

		if empty {
			break
		}

		interpolationArray[next.index] = next.close()
	}

	return strings.TrimSpace(strings.Join(interpolationArray, ""))
}
