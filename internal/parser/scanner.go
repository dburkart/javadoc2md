/*
 * Copyright (c) 2022, Dana Burkart <dana.burkart@gmail.com>
 *
 * SPDX-License-Identifier: BSD-2-Clause
 */

package parser

import (
	"unicode"
	"unicode/utf8"
)

type Scanner struct {
	Name   string
	Input  string
	Tokens chan Token
	State  ScanFn

	Start     int
	Pos       int
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

// SkipJavadocFiller skips the "filler" characters at the beginning of
// each line, which usually constitutes an asterisk and some whitespace.
func (this *Scanner) SkipJavadocFiller() {
	for {
		ch := this.Next()

		// Eat the asterisk, one more space, and bail out
		if ch == '*' && this.Peek() != '/' {
			ch = this.Peek()
			if unicode.IsSpace(ch) && ch != '\n' {
				this.Inc()
			}
			break
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
