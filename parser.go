package main

import "strings"

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
		if t.Lexeme == "@param" {
			fields := strings.Fields(val.Lexeme)
			block.Params[fields[0]] = strings.Join(fields[1:], " ")
		} else {
			block.Tags[t.Lexeme] = val.Lexeme
		}

		t = <- scanner.Tokens
		if t.Type == TOK_JDOC_NL {
			t = <- scanner.Tokens
			continue
		}
		if t.Type != TOK_JDOC_TAG {
			break
		}
	}

	if t.Type == TOK_JDOC_END {
		t = <- scanner.Tokens
	}

	t = ParseJavaContext(scanner, block, t)

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

		if t.Type == TOK_JAVA_KEYWORD {
			if t.Lexeme == "public" || t.Lexeme == "private" {
				block.Attributes["visibility"] = t.Lexeme
				goto next
			}

			if t.Lexeme == "class" {
				t = <- scanner.Tokens

				block.Name = t.Lexeme
				goto next
			}
		}

		if t.Type == TOK_JAVA_IDENTIFIER {
			lastID = t.Lexeme
			goto next
		}

		if t.Type == TOK_JAVA_PAREN_O && block.Name == ""{
			block.Name = lastID
		}

next: 	t = <- scanner.Tokens
	}
}
