/*
 * Copyright (c) 2022, Dana Burkart <dana.burkart@gmail.com>
 *
 * SPDX-License-Identifier: BSD-2-Clause
 */

package parser

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"

	"github.com/dburkart/javadoc2md/internal/logger"
)

// This is Uuuugly. We can do better.
func FormatDefinition(def string) string {
	def = strings.TrimSpace(def)
	def = strings.ReplaceAll(def, " ( ", "(")
	def = strings.ReplaceAll(def, " )", ")")
	def = strings.ReplaceAll(def, " , ", ", ")
	def = strings.ReplaceAll(def, " ,", ",")
	def = strings.ReplaceAll(def, " < ", "<")
	def = strings.ReplaceAll(def, " > ", "> ")

	m := regexp.MustCompile(`(.*)(,)$`)

	return m.Copy().ReplaceAllString(def, "$1")
}

// Given a line, splitKey pulls off the first word, and returns it
// along with the unmodified remainder of the line
func splitKey(line string) (head string, remainder string) {
	first, last, inHead := 0, 0, false
	for index, character := range line {
		if !inHead {
			if unicode.IsSpace(character) {
				continue
			}

			first = index
			inHead = true
		}

		if unicode.IsSpace(character) {
			last = index
			break
		}
	}

	if last == 0 {
		head = line[first:]
		remainder = ""
	} else {
		head = line[first:last]
		remainder = line[last+1:]
	}
	return
}

func ParseDocument(scanner *Scanner, path string) *Document {
	// First, set up our scan loop
	go func() {
		for {
			scanner.State = scanner.State(scanner)
			if scanner.State == nil {
				break
			}
		}
	}()

	doc := MakeDocument(path)

	t := <-scanner.Tokens

	// If we see a package name, save that aside before processing any Javadocs
	if t.Type == TOK_JAVA_KEYWORD && t.Lexeme == "package" {
		t = <-scanner.Tokens
		// TODO: What if it's not an identifier?
		doc.Package = t.Lexeme
		t = <-scanner.Tokens
	}

	for {
		if t.Type == TOK_EOF {
			break
		}

		// Skip anything which isn't the start of javadoc comment
		if t.Type != TOK_JDOC_START {
			t = <-scanner.Tokens
			continue
		}

		t = ParseJavadoc(scanner, doc, t)
	}

	return doc
}

func ParseJavadoc(scanner *Scanner, document *Document, t Token) Token {
	if t.Type != TOK_JDOC_START {
		return t
	}

	// Make our Javadoc block
	block := MakeBlock()
	block.Doc = document

	// Pull off lines until we hit the first tag
	for {
		t = <-scanner.Tokens

		if t.Type == TOK_JDOC_LINE || t.Type == TOK_JDOC_NL || t.Type == TOK_JDOC_PARAM ||
			t.Type == TOK_JSX_O || t.Type == TOK_JSX_X {
			block.Text = append(block.Text, t)
		} else {
			break
		}
	}

	// Add tags to the tag map for the block, until we hit a non-tag
	for {
		if t.Type != TOK_JDOC_TAG {
			break
		}

		val := <-scanner.Tokens
		tagKey := t.Lexeme

		if t.Lexeme == "@param" {
			tagKey, val.Lexeme = splitKey(val.Lexeme)
		}

		// Tags can have multiple lines as their values, so we need to
		// capture all lines until the next tag / end
		for {
			if val.Type != TOK_JDOC_LINE && val.Type != TOK_JDOC_PARAM && val.Type != TOK_JDOC_NL &&
				val.Type != TOK_JSX_O && val.Type != TOK_JSX_X {
				t = val
				break
			}

			if t.Lexeme == "@param" {
				block.Params[tagKey] = append(block.Params[tagKey], val)
			} else {
				block.Tags[tagKey] = append(block.Tags[tagKey], val)
			}
			val = <-scanner.Tokens
		}

		if t.Type != TOK_JDOC_TAG {
			break
		}
	}

	if t.Type == TOK_JDOC_END {
		t = <-scanner.Tokens
	}

	t = ParseJavaContext(scanner, block, t)
	block.Definition = FormatDefinition(block.Definition)

	document.Blocks = append(document.Blocks, *block)

	if block.Name == "" {
		logger.Debug("Could not introspect name from block " + fmt.Sprint(len(document.Blocks)) + " in document " + document.Address)
	}

	return t
}

func ParseJavaContext(scanner *Scanner, block *Block, head Token) Token {
	var lastToken Token

	t := head
	lastID, lastLastID := "", ""
	inArgumentList := false
	for {
		if t.Type < TOK_JAVA_KEYWORD {
			if block.Name == "" {
				block.Name = lastID
			}

			return t
		}

		block.Definition += " " + t.Lexeme

		if t.Type == TOK_JAVA_KEYWORD || t.Type == TOK_JAVA_ANNOTATION {
			if t.Lexeme == "public" || t.Lexeme == "private" {
				block.Attributes["visibility"] = t.Lexeme
				goto next
			}

			switch t.Lexeme {
			case "class", "@class":
				block.Type = SYM_TYPE_CLASS
			case "interface", "@interface":
				block.Type = SYM_TYPE_INTERFACE
			case "enum":
				block.Type = SYM_TYPE_ENUM
			}

			if t.Lexeme == "class" || t.Lexeme == "interface" ||
				t.Lexeme == "@class" || t.Lexeme == "@interface" ||
				t.Lexeme == "enum" {
				t = <-scanner.Tokens

				block.Definition += " " + t.Lexeme

				block.Name = t.Lexeme
				goto next
			}
		}

		if t.Type == TOK_JAVA_IDENTIFIER {
			lastLastID = lastID
			lastID = t.Lexeme
			goto next
		}

		if t.Type == TOK_JAVA_PAREN_O && block.Name == "" {
			block.Name = lastID
			inArgumentList = true
			block.Type = SYM_TYPE_METHOD
		}

		if t.Type == TOK_JAVA_EQUAL && block.Name == "" {
			block.Name = lastID
			block.Type = SYM_TYPE_FIELD
		}

		// Record our arguments
		if inArgumentList {
			pair := ArgPair{Type: lastLastID, Name: lastID}

			if t.Type == TOK_JAVA_PAREN_X {
				inArgumentList = false
				// There were no arguments in this case
				if lastToken.Type == TOK_JAVA_PAREN_O {
					goto next
				}
				block.Arguments = append(block.Arguments, pair)
			}

			if t.Type == TOK_JAVA_COMMA {
				block.Arguments = append(block.Arguments, pair)
			}
		}

	next:
		lastToken = t
		t = <-scanner.Tokens
	}
}
