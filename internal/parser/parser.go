/*
 * Copyright (c) 2022, Dana Burkart <dana.burkart@gmail.com>
 *
 * SPDX-License-Identifier: BSD-2-Clause
 */

package parser

import "strings"

func FormatDefinition(def string) string {
	def = strings.TrimSpace(def)
	def = strings.ReplaceAll(def, " ( ", "(")
	def = strings.ReplaceAll(def, " )", ")")
	def = strings.ReplaceAll(def, " , ", ", ")
	def = strings.ReplaceAll(def, " < ", "<")
	def = strings.ReplaceAll(def, " > ", "> ")

	return def
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

	t := <- scanner.Tokens

	for {
		t := ParseJavadoc(scanner, doc, t)

		if t.Type == TOK_EOF {
			break
		}
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
		t = <- scanner.Tokens

		if t.Type == TOK_JDOC_LINE || t.Type == TOK_JDOC_NL {
			block.Description = block.Description + t.Lexeme
		} else {
			break
		}
	}

	// Add tags to the tag map for the block, until we hit a non-tag
	for {
		if t.Type != TOK_JDOC_TAG {
			break
		}

		val := <- scanner.Tokens
		tagKey := t.Lexeme
		tagValue := ""

		if t.Lexeme == "@param" {
			fields := strings.Fields(val.Lexeme)
			tagKey = fields[0]
			val.Lexeme = strings.Join(fields[1:], " ")
		}

		// Tags can have multiple lines as their values, so we need to
		// capture all lines until the next tag / end
		for {
			if val.Type != TOK_JDOC_LINE {
				if t.Lexeme == "@param" {
					block.Params[tagKey] = tagValue
				} else {
					block.Tags[tagKey] = tagValue
				}
				t = val
				break
			}

			tagValue = tagValue + " " + val.Lexeme
			val = <- scanner.Tokens
			if val.Type == TOK_JDOC_NL {
				val = <- scanner.Tokens
			}
		}

		if t.Type != TOK_JDOC_TAG {
			break
		}
	}

	if t.Type == TOK_JDOC_END {
		t = <- scanner.Tokens
	}

	t = ParseJavaContext(scanner, block, t)
	block.Definition = FormatDefinition(block.Definition)

	document.Blocks = append(document.Blocks, *block)

	return t
}

func ParseJavaContext(scanner *Scanner, block *Block, head Token) Token {
	t := head
	lastID := ""
	for {
		if t.Type < TOK_JAVA_KEYWORD {
			return t
		}

		block.Definition += " " + t.Lexeme

		if t.Type == TOK_JAVA_KEYWORD || t.Type == TOK_JAVA_ANNOTATION {
			if t.Lexeme == "public" || t.Lexeme == "private" {
				block.Attributes["visibility"] = t.Lexeme
				goto next
			}

			if t.Lexeme == "class" || t.Lexeme == "interface" ||
			   t.Lexeme == "@class" || t.Lexeme == "@interface" {
				block.Doc.Type = t.Lexeme
				t = <- scanner.Tokens

				block.Definition += " " + t.Lexeme

				block.Name = t.Lexeme
				goto next
			}
		}

		if t.Type == TOK_JAVA_IDENTIFIER {
			lastID = t.Lexeme
			goto next
		}

		if t.Type == TOK_JAVA_PAREN_O && block.Name == "" {
			block.Name = lastID
		}

next: 	t = <- scanner.Tokens
	}
}