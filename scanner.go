/*
 * Copyright (c) 2022, Dana Burkart <dana.burkart@gmail.com>
 *
 * SPDX-License-Identifier: BSD-2-Clause
 */

package main

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

type Scanner struct {
	Name string
	Input string
	Tokens chan Token
	State ScanFn

	Start int
	Pos int
	RuneWidth int
}

type ScanFn func(*Scanner) ScanFn

func (this *Scanner) Emit(tokenType TokenType) {
	this.Tokens <- Token{Type: tokenType, Lexeme: this.Input[this.Start:this.Pos]}
	this.Start = this.Pos
}

func (this *Scanner) Inc() {
	this.Pos++
	if this.Pos >= utf8.RuneCountInString(this.Input) {
		this.Emit(TOK_EOF)
	}
}

func (this *Scanner) Dec() {
	this.Pos--
}

func (this *Scanner) Next() rune {
	if this.Pos >= utf8.RuneCountInString(this.Input) {
		return EOF
	}

	result, width := utf8.DecodeRuneInString(this.Input[this.Pos:])

	this.Pos += width
	this.RuneWidth = width
	return result
}

func (this *Scanner) Rewind() {
	this.Pos -= this.RuneWidth
}

func (this *Scanner) Peek() rune {
	ch := this.Next()
	this.Rewind()
	return ch
}

func (this *Scanner) InputToEnd() string {
	return this.Input[this.Pos:]
}

func (this *Scanner) SkipWhitespace() {
	for {
		ch := this.Next()

		if !unicode.IsSpace(ch) {
			this.Dec()
			break
		}

		if ch == EOF {
			this.Emit(TOK_EOF)
			break
		}
	}

	this.Start = this.Pos
}

// Skips whitespace, except for newlines. This function is useful within
// Javadoc comments since newlines are sometimes significant.
func (this *Scanner) SkipLinearWhitespace() {
	for {
		ch := this.Next()

		// Consider '*' whitespace in this function, as long as it's not followed by /
		if ch == '*' && this.Peek() != '/' {
			continue
		}

		if ch == '\n' || !unicode.IsSpace(ch) {
			this.Dec()
			break
		}

		if ch == EOF {
			this.Emit(TOK_EOF)
			break
		}
	}

	this.Start = this.Pos
}

// Scanning functions

func ScanBegin(scanner *Scanner) ScanFn {
	scanner.SkipWhitespace()

	for {
		// First, check if a JavaDoc is beginning
		if strings.HasPrefix(scanner.InputToEnd(), "/**") {
			return ScanJavadocStart
		}

		ch := scanner.Next()
		scanner.Start = scanner.Pos

		if ch == EOF {
			scanner.Emit(TOK_EOF)
		}
	}
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
	scanner.SkipLinearWhitespace()

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

	return nil
}

func ScanJavadocTag(scanner *Scanner) ScanFn {
	for {
		ch := scanner.Next()

		if unicode.IsSpace(ch) {
			scanner.Rewind()
			scanner.Emit(TOK_JDOC_TAG)
			return ScanJavadocLine
		}
	}
}

func ScanJavaLine(scanner *Scanner) ScanFn {
	for {
		scanner.SkipWhitespace()

		ch := scanner.Peek()

		if ch == '{' {
			return ScanBegin
		}

		if strings.HasPrefix(scanner.InputToEnd(), "/**") {
			return ScanBegin
		}

		switch (ch) {
			case '(':
				scanner.Inc()
				scanner.Emit(TOK_JAVA_PAREN_O)
				continue
			case ')':
				scanner.Inc()
				scanner.Emit(TOK_JAVA_PAREN_X)
				continue
			case ',':
				scanner.Inc()
				scanner.Emit(TOK_JAVA_COMMA)
				continue
			case 'c':
				if strings.HasPrefix(scanner.InputToEnd(), "class") {
					scanner.Pos += len("class")
					scanner.Emit(TOK_JAVA_KEYWORD)
					continue
				}
			case 'e':
				if strings.HasPrefix(scanner.InputToEnd(), "extends") {
					scanner.Pos += len("extends")
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
		}

		// Pull characters off until we have an identifier
		for {
			if unicode.IsSpace(ch) {
				scanner.Emit(TOK_JAVA_IDENTIFIER)
				break
			}

			if (ch < 'a' || ch > 'z') && (ch < 'A' || ch > 'Z') && (ch < '0' || ch > '9'){
				scanner.Emit(TOK_JAVA_IDENTIFIER)
				break
			}

			scanner.Inc()
			ch = scanner.Peek()
		}
	}
}

func BeginScanning(name, input string) *Scanner {
	s := &Scanner{
		Name: name,
		Input: input,
		State: ScanBegin,
		Tokens: make(chan Token, 3),
	}

	return s
}
