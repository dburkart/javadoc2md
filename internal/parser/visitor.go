/*
 * Copyright (c) 2022, Dana Burkart <dana.burkart@gmail.com>
 *
 * SPDX-License-Identifier: BSD-2-Clause
 */

package parser

import (
	"os"
	"path/filepath"
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
		symbol := Symbol{Type: block.Type, Package: doc.Package, Name: block.Name, QualifiedName: block.Name}
		if i == 0 {
			symbol.Location = block.Name
			v.Symbols[block.Name] = symbol
			v.Symbols[doc.Package+"."+block.Name] = symbol
		} else {
			// First, put together the symbol's qualified name
			qualifiedName := block.Name
			if symbol.Type == SYM_TYPE_METHOD {
				qualifiedName = block.Name + "("
				numArgs := len(block.Arguments)
				// For each argument, add to the symbol name
				for i, val := range block.Arguments {
					qualifiedName += val.Type
					if i < numArgs-1 {
						qualifiedName += ","
					}
				}
				qualifiedName += ")"
			}

			symbol.QualifiedName = qualifiedName
			doc.Blocks[i].QualifiedName = qualifiedName

			symbolName := doc.Blocks[0].Name + "#" + block.Name
			symbol.Location = doc.Blocks[0].Name + "#" + qualifiedName

			// If the vague, argument-less symbol already exists in the map
			// we want to only insert the exact symbol name below.
			if _, ok := v.Symbols[symbolName]; !ok {
				v.Symbols[symbolName] = symbol
				v.Symbols[doc.Package+"."+symbolName] = symbol
			}

			// Generate Qualified Name
			symbolName = doc.Blocks[0].Name + "#" + qualifiedName
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
		sectionName := "`" + v.Definition + "` {#" + v.QualifiedName + "}"

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
			f.WriteString(ret.Interpolate(doc, m.Symbols, "") + "\n\n")
			f.WriteString(":::\n\n")
		}

		f.WriteString(v.Text.Interpolate(doc, m.Symbols, ""))
		f.WriteString("\n\n")

		if len(v.Arguments) > 0 {
			f.WriteString("**Parameters:**" + "\n\n")
			needs_newline = true
		}

		var paramOrder []string
		resolvedParams := map[string]string{}
		// First iterate over any arguments, detecting undocumented fields in the process
		for _, value := range v.Arguments {
			paramOrder = append(paramOrder, value.Name)
			if description, found := v.Params[value.Name]; found {
				resolvedParams[value.Name] = description.Interpolate(doc, m.Symbols, "")
			} else {
				resolvedParams[value.Name] = "*Undocumented*"
			}
		}
		// Now iterate over anything "extra" in our params
		for k, v := range v.Params {
			if _, found := resolvedParams[k]; !found {
				paramOrder = append(paramOrder, k)
				resolvedParams[k] = v.Interpolate(doc, m.Symbols, "")
			}
		}
		// Finally, write out all params
		for _, p := range paramOrder {
			f.WriteString("* `" + p + "` - " + resolvedParams[p] + "\n")
		}

		if ret, found := v.Tags["@return"]; found {
			f.WriteString("\n**Returns:** " + ret.Interpolate(doc, m.Symbols, "") + "\n\n")
			needs_newline = true
		}

		if needs_newline {
			f.WriteString("\n")
		}
		needs_newline = false
	}
	return
}
