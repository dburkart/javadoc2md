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
    symbolVisitor := SymbolVisitor{Symbols: make(map[string]string)}
    for {
        doc, ok := <- docs
        if !ok {
            break
        }

        symbolVisitor.visit(doc)
        documents = append(documents, doc)
    }

    visitors := []Visitor{
        &MarkdownVisitor{OutputDirectory: options.OutputDirectory, SkipPrivateDefs: options.SkipPrivateDefs},
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
    SkipPrivateDefs bool
}

func (m *MarkdownVisitor) interpolateText(tokens []Token) string {
    interpolatedText := ""

    for i := 0; i < len(tokens); i++ {
        token := tokens[i]

        switch token.Type {
        case TOK_JDOC_NL:
        case TOK_JDOC_LINE:
            interpolatedText += token.Lexeme
        case TOK_JDOC_PARAM:
            if token.Lexeme == "@code" {
                interpolatedText += "`"
                interpolatedText += tokens[i+1].Lexeme
                interpolatedText += "` "
                i++
            }
        }
    }

    return interpolatedText
}

func (m *MarkdownVisitor) visit(doc *Document) (err bool, description string) {
    err = false
    description = ""
    needs_newline := false

    if m.SkipPrivateDefs && doc.Blocks[0].Attributes["visibility"] == "private" {
        return
    }

    f, createErr := os.Create(filepath.Join(m.OutputDirectory, doc.Blocks[0].Name + ".md"))
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

        // Before writing out content, write out any deprecated admonitions
        if ret, found := v.Tags["@deprecated"]; found {
            f.WriteString(":::caution Deprecated\n\n")
            f.WriteString(strings.TrimSpace(ret) + "\n\n")
            f.WriteString(":::\n\n")
        }

        f.WriteString(m.interpolateText(v.Text))
        f.WriteString("\n\n")

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
