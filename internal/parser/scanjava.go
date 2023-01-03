/*
 * Copyright (c) 2022-2023, Dana Burkart <dana.burkart@gmail.com>
 *
 * SPDX-License-Identifier: BSD-2-Clause
 */

package parser

import (
	"strings"
	"unicode"
)

func ScanBegin(scanner *Scanner) ScanFn {
	scanner.SkipWhitespace()
	inComment := false

	for {
		if inComment {
			if strings.HasPrefix(scanner.InputToEnd(), "*/") {
				scanner.Pos += 2
				scanner.Start = scanner.Pos
				inComment = false
				continue
			}

			goto advance
		}

		// First, check if a JavaDoc is beginning
		if strings.HasPrefix(scanner.InputToEnd(), "/**") {
			return ScanJavadocStart
		}

		// Ignore anything in a comment
		if strings.HasPrefix(scanner.InputToEnd(), "/*") {
			inComment = true
			continue
		}

		// Check for a package name
		if strings.HasPrefix(scanner.InputToEnd(), "package") {
			return ScanPackageStatement
		}

	advance:
		ch := scanner.Next()
		scanner.Start = scanner.Pos

		if ch == EOF {
			scanner.Emit(TOK_EOF)
			return nil
		}
	}
}

func ScanPackageStatement(scanner *Scanner) ScanFn {
	scanner.Pos += len("package")
	scanner.Emit(TOK_JAVA_KEYWORD)
	return ScanPackageName
}

func ScanPackageName(scanner *Scanner) ScanFn {
	scanner.SkipWhitespace()
	for {
		ch := scanner.Peek()

		if (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') ||
			(ch >= '0' && ch <= '9') || ch == '.' || ch == '_' {
			scanner.Inc()
		} else {
			scanner.Emit(TOK_JAVA_IDENTIFIER)
			break
		}
	}
	return ScanBegin
}

func ScanJavadocStart(scanner *Scanner) ScanFn {
	scanner.Pos += len("/**")
	scanner.Emit(TOK_JDOC_START)
	return ScanJavadoc
}

func ScanJavadocEnd(scanner *Scanner) ScanFn {
	scanner.Pos += len("*/")
	scanner.Emit(TOK_JDOC_END)
	return ScanJavaLine
}

func ScanJavadoc(scanner *Scanner) ScanFn {
	scanner.SkipJavadocFiller()

	if strings.HasPrefix(scanner.InputToEnd(), "@") {
		return ScanJavadocTag
	}

	if strings.HasPrefix(scanner.InputToEnd(), "*/") {
		return ScanJavadocEnd
	}

	return ScanJavadocLine
}

func ScanJavadocLine(scanner *Scanner) ScanFn {

	for {
		ch := scanner.Peek()

		if ch == '*' {
			scanner.Inc()

			if scanner.Peek() == '/' {
				scanner.Dec()
				if scanner.Pos > scanner.Start {
					scanner.Emit(TOK_JDOC_LINE)
				}
				return ScanJavadocEnd
			}

			scanner.Dec()
		}

		if ch == '<' {
			if scanner.Pos > scanner.Start {
				scanner.Emit(TOK_JDOC_LINE)
			}
			return ScanJSXTag
		}

		if ch == '{' {
			if scanner.Pos > scanner.Start {
				scanner.Emit(TOK_JDOC_LINE)
			}
			return ScanJavadocParam
		}

		if ch == '\n' {
			if scanner.Pos > scanner.Start {
				scanner.Emit(TOK_JDOC_LINE)
			}
			scanner.Inc()
			scanner.Emit(TOK_JDOC_NL)
			return ScanJavadoc
		}

		scanner.Inc()
	}
}

func ScanJSXTag(scanner *Scanner) ScanFn {
	position := scanner.Pos
	// Eat '<'
	scanner.Inc()
	closeTag := scanner.Peek() == '/'

	// Iterate over runes until we find a matching '>'
	for {
		c := scanner.Next()

		if c == '\n' {
			scanner.Pos = position + 1
			return ScanJavadoc
		}

		if c == '>' {
			if closeTag {
				scanner.Emit(TOK_JSX_X)
			} else {
				scanner.Emit(TOK_JSX_O)
			}
			return ScanJavadocLine
		}
	}
}

func ScanJavadocTag(scanner *Scanner) ScanFn {
	for {
		ch := scanner.Next()

		if unicode.IsSpace(ch) {
			scanner.Rewind()
			lexeme := scanner.Input[scanner.Start:scanner.Pos]
			scanner.Emit(TOK_JDOC_TAG)

			if lexeme == "@param" {
				goto removeKey
			}

			return ScanJavadocLine
		}
	}

removeKey:
	scanner.Inc()
	scanner.Start += 1

	for {
		ch := scanner.Next()

		if unicode.IsSpace(ch) {
			scanner.Rewind()
			scanner.Emit(TOK_JDOC_LINE)
			return ScanJavadocLine
		}
	}
}

func ScanJavadocParam(scanner *Scanner) ScanFn {
	// Consume '{'
	scanner.Inc()

	// The next rune must be an '@' symbol for this to be a valid param
	if scanner.Peek() != '@' {
		return ScanJavadocLine
	}

	scanner.Start = scanner.Pos

	insideParam := false

	for {
		ch := scanner.Next()

		if ch == '}' {
			scanner.Rewind()
			if insideParam {
				scanner.Emit(TOK_JDOC_PARAM)
			} else {
				scanner.Emit(TOK_JDOC_LINE)
			}

			scanner.Inc()
			scanner.Emit(TOK_JDOC_PARAM_END)

			return ScanJavadocLine
		}

		if ch == '@' {
			insideParam = true
			continue
		}

		if ch == '\n' {
			scanner.Rewind()
			if insideParam {
				scanner.Emit(TOK_JDOC_PARAM)
				insideParam = false
			}

			// Emit a line if we have one
			if scanner.Start != scanner.Pos {
				scanner.Emit(TOK_JDOC_LINE)
			}

			scanner.Inc()
			scanner.Start += 1

			// Eat whitespace / "*"
			for {
				ch = scanner.Next()

				if unicode.IsSpace(ch) || ch == '*' {
					scanner.Start += 1
				} else {
					break
				}
			}

			scanner.Rewind()
		}

		if unicode.IsSpace(ch) && insideParam {
			scanner.Rewind()
			scanner.Emit(TOK_JDOC_PARAM)
			scanner.Inc()
			scanner.Start += 1
			insideParam = false
		}
	}
}

func ScanJavaLine(scanner *Scanner) ScanFn {
	for {
		scanner.SkipWhitespace()

		ch := scanner.Peek()

		if ch == ';' {
			return ScanBegin
		}

		if ch == '{' || ch == '}' {
			return ScanBegin
		}

		if strings.HasPrefix(scanner.InputToEnd(), "/**") {
			return ScanBegin
		}

		switch ch {
		case '.':
			scanner.Inc()
			scanner.Emit(TOK_JAVA_OPERATOR)
			continue
		case '?':
			scanner.Inc()
			scanner.Emit(TOK_JAVA_OTHER)
			continue
		case '>':
			if strings.HasPrefix(scanner.InputToEnd(), ">>") {
				scanner.Pos += 2
				scanner.Emit(TOK_JAVA_OPERATOR)
			}

			scanner.Inc()
			scanner.Emit(TOK_JAVA_OTHER)
			continue
		case '<':
			if strings.HasPrefix(scanner.InputToEnd(), "<<") {
				scanner.Pos += 2
				scanner.Emit(TOK_JAVA_OPERATOR)
			}
			scanner.Inc()
			scanner.Emit(TOK_JAVA_OTHER)
			continue
		case '~':
			scanner.Inc()
			scanner.Emit(TOK_JAVA_OPERATOR)
			continue
		case '&':
			scanner.Inc()
			scanner.Emit(TOK_JAVA_OPERATOR)
			continue
		case '(':
			scanner.Inc()
			scanner.Emit(TOK_JAVA_PAREN_O)
			continue
		case ')':
			scanner.Inc()
			scanner.Emit(TOK_JAVA_PAREN_X)
			continue
		case '[':
			scanner.Inc()
			scanner.Emit(TOK_JAVA_BRACKET_O)
			continue
		case ']':
			scanner.Inc()
			scanner.Emit(TOK_JAVA_BRACKET_X)
			continue
		case ',':
			scanner.Inc()
			scanner.Emit(TOK_JAVA_COMMA)
			continue
		case '=':
			scanner.Inc()
			scanner.Emit(TOK_JAVA_EQUAL)
			continue
		case '"':
			scanner.Inc()
			for {
				ch := scanner.Next()

				if ch == '"' {
					scanner.Emit(TOK_JAVA_STRING)
					break
				}
			}
			continue
		case '/':
			if strings.HasPrefix(scanner.InputToEnd(), "/*") {
				scanner.Pos += 2
				scanner.Emit(TOK_JAVA_COMMENT_O)
				continue
			}

			if strings.HasPrefix(scanner.InputToEnd(), "//") {
				for {
					if scanner.Peek() == '\n' {
						scanner.Inc()
						scanner.Start = scanner.Pos
						break
					}

					scanner.Inc()
				}
				continue
			}

			scanner.Pos += 1
			scanner.Emit(TOK_JAVA_OPERATOR)
			continue
		case '*':
			if strings.HasPrefix(scanner.InputToEnd(), "*/") {
				scanner.Pos += 2
				scanner.Emit(TOK_JAVA_COMMENT_O)
				continue
			}

			scanner.Pos += 1
			scanner.Emit(TOK_JAVA_OPERATOR)
			continue
		case '-':
			if strings.HasPrefix(scanner.InputToEnd(), "->") {
				scanner.Pos += 2
				scanner.Emit(TOK_JAVA_OPERATOR)
				continue
			}

			scanner.Pos += 1
			scanner.Emit(TOK_JAVA_OPERATOR)
			continue
		case '+':
			scanner.Pos += 1
			scanner.Emit(TOK_JAVA_OPERATOR)
			continue
		case 'c':
			if strings.HasPrefix(scanner.InputToEnd(), "class") {
				scanner.Pos += len("class")
				scanner.Emit(TOK_JAVA_KEYWORD)
				continue
			}
		case 'e':
			if strings.HasPrefix(scanner.InputToEnd(), "enum") {
				scanner.Pos += len("enum")
				scanner.Emit(TOK_JAVA_KEYWORD)
				continue
			}

			if strings.HasPrefix(scanner.InputToEnd(), "extends") {
				scanner.Pos += len("extends")
				scanner.Emit(TOK_JAVA_KEYWORD)
				continue
			}
		case 'i':
			if strings.HasPrefix(scanner.InputToEnd(), "interface") {
				scanner.Pos += len("interface")
				scanner.Emit(TOK_JAVA_KEYWORD)
				continue
			}
		case 'p':
			if strings.HasPrefix(scanner.InputToEnd(), "public") {
				scanner.Pos += len("public")
				scanner.Emit(TOK_JAVA_KEYWORD)
				continue
			}

			if strings.HasPrefix(scanner.InputToEnd(), "private") {
				scanner.Pos += len("private")
				scanner.Emit(TOK_JAVA_KEYWORD)
				continue
			}
		case 's':
			if strings.HasPrefix(scanner.InputToEnd(), "static") {
				scanner.Pos += len("static")
				scanner.Emit(TOK_JAVA_KEYWORD)
				continue
			}
		case '@':
			braces := 0
			for {
				if unicode.IsSpace(ch) && braces == 0 {
					scanner.Emit(TOK_JAVA_ANNOTATION)
					break
				}

				if ch == '(' || ch == '{' {
					braces += 1
				}

				if ch == ')' || ch == '}' {
					braces -= 1
				}

				scanner.Inc()
				ch = scanner.Peek()
			}
			continue
		}

		if (ch > '0' && ch < '9') || ch == '.' {
			for {
				if unicode.IsSpace(ch) {
					scanner.Emit(TOK_JAVA_NUMERIC)
					break
				}

				if (ch < '0' || ch > '9') && (ch < 'A' || ch > 'Z') && ch != '.' {
					scanner.Emit(TOK_JAVA_NUMERIC)
					break
				}

				scanner.Inc()
				ch = scanner.Peek()
			}
			continue
		}

		// Pull characters off until we have an identifier
		for {
			if unicode.IsSpace(ch) {
				scanner.Emit(TOK_JAVA_IDENTIFIER)
				break
			}

			if (ch < 'a' || ch > 'z') && (ch < 'A' || ch > 'Z') &&
				(ch < '0' || ch > '9') && ch != '_' && ch != '$' {
				scanner.Emit(TOK_JAVA_IDENTIFIER)
				break
			}

			scanner.Inc()
			ch = scanner.Peek()
		}
	}
}

func BeginScanningJavaCode(name, input string) *Scanner {
	s := &Scanner{
		Name:   name,
		Input:  input,
		State:  ScanBegin,
		Tokens: make(chan Token, 3),
	}

	return s
}
