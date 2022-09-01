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
}

func VisitDocuments(options *VisitorConfigOptions, docs *[]*Document) {
	visitors := []Visitor{
		&SymbolVisitor{Symbols: make(map[string]string)},
		&MarkdownVisitor{OutputDirectory: options.OutputDirectory},
	}

	for _, v := range visitors {
		for _, d := range *docs {
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
	Symbols map[string]string
}

func (v *SymbolVisitor) visit(doc *Document) (err bool, description string) {
	err = false
	description = ""

	for i, block := range doc.Blocks {
		if i == 0 {
			v.Symbols[block.Name] = block.Name
		} else {
			symbolName := doc.Blocks[0].Name + "#" + block.Name
			v.Symbols[symbolName] = symbolName
		}
	}

	return
}

type MarkdownVisitor struct {
	OutputDirectory string
}

func (v *MarkdownVisitor) visit(doc *Document) (err bool, description string) {
	err = false
	description = ""
	needs_newline := false

	f, createErr := os.Create(filepath.Join(v.OutputDirectory, doc.Blocks[0].Name + ".md"))
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

		// f.WriteString(heading + v.Name + "\n\n")
		f.WriteString(heading + sectionName + "\n\n")

		// Write out the definition separately if this is the first block
		if i == 0 {
			f.WriteString("```java\n" + v.Definition + "\n```\n\n")
		}

		f.WriteString(strings.TrimSpace(v.Description) + "\n\n")

		if len(v.Params) > 0 {
			f.WriteString("* **Parameters:**" + "\n")
			needs_newline = true
		}

		for key, value := range v.Params {
			f.WriteString("\t* `" + key + "` -" + value + "\n")
		}

		if ret, found := v.Tags["@return"]; found {
			f.WriteString("* **Returns:** " + ret + "\n")
			needs_newline = true
		}

		if needs_newline {
			f.WriteString("\n")
		}
		needs_newline = false
	}
	return
}