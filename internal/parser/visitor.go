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
		} else {
			symbolName := doc.Blocks[0].Name + "#" + block.Name
			symbol.Location = symbolName
			v.Symbols[symbolName] = symbol
		}
	}

	return
}

type MarkdownVisitor struct {
	OutputDirectory string
	SkipPrivateDefs bool
	Symbols         map[string]Symbol
}

// Given a MixedText token list, return a string with all the parameters
// evaluated.
func (m *MarkdownVisitor) interpolateText(tokens MixedText, doc *Document) string {
	interpolationArray := make([]string, len(tokens))

	for i := 0; i < len(tokens); i++ {
		token := tokens[i]

		switch token.Type {
		case TOK_JDOC_NL:
			interpolationArray[i] = "\n"
		case TOK_JDOC_LINE:
			interpolationArray[i] = strings.TrimSpace(token.Lexeme) + " "
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

				symbol := m.Symbols[target]
				// TODO: The name of the link should be a proper definition
				str += "[" + symbol.Name + "](" + symbol.Location + ")"
				i++
			}

			interpolationArray[i] = str + " "
		}
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
		heading := "## "
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

			f.WriteString("```java\n" + v.Definition + "\n```\n\n")
		}

		// Before writing out content, write out any deprecated admonitions
		if ret, found := v.Tags["@deprecated"]; found {
			f.WriteString(":::caution Deprecated\n\n")
			f.WriteString(m.interpolateText(ret, doc) + "\n\n")
			f.WriteString(":::\n\n")
		}

		f.WriteString(m.interpolateText(v.Text, doc))
		f.WriteString("\n\n")

		if len(v.Params) > 0 {
			f.WriteString("* **Parameters:**" + "\n")
			needs_newline = true
		}

		// TODO: I'm not sure what the best way to ensure we get everything
		//       from Javadocs written out since it may contain more than we
		//       captured.
		for _, value := range v.Arguments {
			if description, found := v.Params[value.Name]; found {
				f.WriteString("\t* `" + value.Name + "` - " + m.interpolateText(description, doc) + "\n")
			} else {
				f.WriteString("\t* `" + value.Name + "` - *Undocumented*\n")
			}
		}

		if ret, found := v.Tags["@return"]; found {
			f.WriteString("* **Returns:** " + m.interpolateText(ret, doc) + "\n")
			needs_newline = true
		}

		if needs_newline {
			f.WriteString("\n")
		}
		needs_newline = false
	}
	return
}
