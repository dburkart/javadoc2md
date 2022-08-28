package main

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

	for {
		t := ParseJavadoc(scanner, doc)

		if t.Type == TOK_EOF {
			break
		}
	}

	return doc
}

func ParseJavadoc(scanner *Scanner, document *Document) Token {
	t := <- scanner.Tokens

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
		block.Tags[t.Lexeme] = val.Lexeme

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

	document.Blocks = append(document.Blocks, *block)

	return t
}
