/*
 * Copyright (c) 2022, Dana Burkart <dana.burkart@gmail.com>
 *
 * SPDX-License-Identifier: BSD-2-Clause
 */

package parser

import (
	"os"
	"path/filepath"
	"strings"
	"unicode"
)

type VisitorConfigOptions struct {
	OutputDirectory string
	SkipPrivateDefs bool
}

func VisitDocuments(options *VisitorConfigOptions, docs chan *Document) {
	var documents []*Document

	// The symbol visitor is special in that we want to visit _every_ document
	// with this visitor before proceeding
	symbolVisitor := SymbolVisitor{Symbols: make(map[string]Symbol)}
	for {
		doc, ok := <-docs
		if !ok {
			break
		}

		symbolVisitor.visit(doc)
		documents = append(documents, doc)
	}

	visitors := []Visitor{
		&MarkdownVisitor{
			OutputDirectory: options.OutputDirectory,
			SkipPrivateDefs: options.SkipPrivateDefs,
			Symbols:         symbolVisitor.Symbols,
		},
	}

	for _, v := range visitors {
		for _, d := range documents {
			if len(d.Blocks) == 0 {
				continue
			}

			v.visit(d)
		}
	}
}

type Visitor interface {
	visit(*Document) (bool, string)
}

type SymbolVisitor struct {
	Symbols map[string]Symbol
}

func (v *SymbolVisitor) visit(doc *Document) (err bool, description string) {
	err = false
	description = ""

	for i, block := range doc.Blocks {
		symbol := Symbol{Type: block.Type, Package: doc.Package, Name: block.Name}
		if i == 0 {
			symbol.Location = block.Name
			v.Symbols[block.Name] = symbol
			v.Symbols[doc.Package+"."+block.Name] = symbol
		} else {
			symbolName := doc.Blocks[0].Name + "#" + block.Name
			symbol.Location = symbolName
			v.Symbols[symbolName] = symbol
			v.Symbols[doc.Package+"."+symbolName] = symbol
		}
	}

	return
}

// The MarkdownVisitor is responsible for emitting a markdown document for
// each Document.
type MarkdownVisitor struct {
	OutputDirectory string
	SkipPrivateDefs bool
	Symbols         map[string]Symbol
}

type jsxTag struct {
	index int
	tag   string
}

func (j *jsxTag) tagType() string {
	start, end := 1, 0

	for i, value := range j.tag {
		if unicode.IsSpace(value) {
			end = i - 1
			break
		}

		if value == '/' && i == len(j.tag)-1 {
			end = i - 1
			break
		} else if value == '/' {
			start = i + 1
		}

		if value == '>' {
			end = i
			break
		}
	}

	return j.tag[start:end]
}

func (j *jsxTag) close() string {
	isClosed := false
	for i, value := range j.tag {
		if value == '>' && j.tag[i-1] == '/' {
			isClosed = true
		}
	}

	if !isClosed {
		j.tag = j.tag[:len(j.tag)-1] + "/>"
	}

	return j.tag
}

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

// Given a Text token list, return a string with all the parameters
// evaluated.
func interpolateText(tokens Text, doc *Document, symbols SymbolMap, flowIndent string) string {
	interpolationArray := make([]string, len(tokens))
	jsxStack := stack{}

	for i := 0; i < len(tokens); i++ {
		token := tokens[i]

		switch token.Type {
		case TOK_JDOC_NL:
			interpolationArray[i] = "\n" + flowIndent
		case TOK_JDOC_PARAM:
			str := ""
			if token.Lexeme == "@code" {
				str += "`"
				str += strings.TrimSpace(tokens[i+1].Lexeme)
				str += "`"
				i++
			}

			if token.Lexeme == "@link" {
				target := strings.TrimSpace(tokens[i+1].Lexeme)

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

func (m *MarkdownVisitor) visit(doc *Document) (err bool, description string) {
	err = false
	description = ""
	needs_newline := false

	if m.SkipPrivateDefs && doc.Blocks[0].Attributes["visibility"] == "private" {
		return
	}

	f, createErr := os.Create(filepath.Join(m.OutputDirectory, doc.Blocks[0].Name+".md"))
	if createErr != nil {
		err = true
		description = createErr.Error()
		return
	}
	defer f.Close()

	for i, v := range doc.Blocks {
		heading := "### "
		sectionName := "`" + v.Definition + "` {#" + v.Name + "}"

		if i == 0 {
			heading = "# "
			sectionName = v.Name
		}

		f.WriteString(heading + sectionName + "\n\n")

		// Write out the definition separately if this is the first block
		if i == 0 {
			if doc.Package != "" {
				f.WriteString("```java\n")
				f.WriteString("import " + doc.Package + "." + v.Name + "\n```\n\n")
			}

			f.WriteString("## Definition\n\n")
			f.WriteString("```java\n" + v.Definition + "\n```\n\n")

			f.WriteString("## Overview\n\n")
		}

		// Before writing out content, write out any deprecated admonitions
		if ret, found := v.Tags["@deprecated"]; found {
			f.WriteString(":::caution Deprecated\n\n")
			f.WriteString(interpolateText(ret, doc, m.Symbols, "") + "\n\n")
			f.WriteString(":::\n\n")
		}

		f.WriteString(interpolateText(v.Text, doc, m.Symbols, ""))
		f.WriteString("\n\n")

		if len(v.Arguments) > 0 {
			f.WriteString("**Parameters:**" + "\n\n")
			needs_newline = true
		}

		// TODO: I'm not sure what the best way to ensure we get everything
		//       from Javadocs written out since it may contain more than we
		//       captured.
		for _, value := range v.Arguments {
			if description, found := v.Params[value.Name]; found {
				f.WriteString("* `" + value.Name + "` - " + interpolateText(description, doc, m.Symbols, "\t  ") + "\n")
			} else {
				f.WriteString("* `" + value.Name + "` - *Undocumented*\n")
			}
		}

		if ret, found := v.Tags["@return"]; found {
			f.WriteString("**Returns:** " + interpolateText(ret, doc, m.Symbols, "") + "\n\n")
			needs_newline = true
		}

		if needs_newline {
			f.WriteString("\n")
		}
		needs_newline = false
	}
	return
}
