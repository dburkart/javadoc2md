/*
 * Copyright (c) 2022-2023, Dana Burkart <dana.burkart@gmail.com>
 *
 * SPDX-License-Identifier: BSD-2-Clause
 */

package parser

import (
	"strings"
)

type stack []XMLTag

func (s *stack) Empty() bool {
	return len(*s) == 0
}
func (s *stack) Push(j XMLTag) {
	*s = append(*s, j)
}

func (s *stack) Pop() (j XMLTag, empty bool) {
	if s.Empty() {
		j = XMLTag{}
		empty = true
	} else {
		i := len(*s) - 1
		j = (*s)[i]
		*s = (*s)[:i]
		empty = false
	}
	return
}

func (s *stack) Peek() (j XMLTag, empty bool) {
	if s.Empty() {
		j = XMLTag{}
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
				inPre := false
				// Check if for weird interactions, like being inside a <pre> tag, which
				// cancels us out.
				for _, tag := range jsxStack {
					if tag.Type() == "pre" {
						inPre = true
						break
					}
				}

				if !inPre {
					str += "`"
				}

				str += strings.TrimSpace((*t)[i+1].Lexeme)

				if !inPre {
					str += "`"
				}
				i++
			}

			if token.Lexeme == "@link" || token.Lexeme == "@linkplain" {
				target := strings.ReplaceAll((*t)[i+1].Lexeme, " ", "")
				target = strings.ReplaceAll(target, "\n", "")

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
			jsxStack.Push(XMLTag{i, token.Lexeme})
			interpolationArray[i] = token.Lexeme
		case TOK_JSX_X:
			current := XMLTag{i, token.Lexeme}

			for {
				next, empty := jsxStack.Pop()

				if empty {
					break
				}

				if next.Type() == current.Type() {
					if next.Type() == "pre" || next.Type() == "code" {
						// If there are multiple lines between our tags, consider it a long form code block
						if i-next.Index > 2 || strings.IndexRune(interpolationArray[i-1], '\n') != -1 {
							interpolationArray[next.Index] = "```java"
							token.Lexeme = "```"
						} else { // Otherwise it's inline plaintext
							interpolationArray[next.Index] = "`"
							token.Lexeme = "`"
						}
					}

					if next.Type() == "a" {
						interpolationArray[next.Index] = "["
						token.Lexeme = "](" + next.Attributes()["href"] + ")"
					}
					break
				}

				// Close the Tag
				interpolationArray[next.Index] = next.Close()
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

		interpolationArray[next.Index] = next.Close()
	}

	return strings.TrimSpace(strings.Join(interpolationArray, ""))
}
