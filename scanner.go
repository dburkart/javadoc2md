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

const EOF rune = 0

type TokenType int

const (
	TOK_INVALID TokenType = iota
	TOK_EOF
	
	// This content exists only in JavaDocs
	TOK_JDOC_START					// /**
	TOK_JDOC_END					//  */
	TOK_JDOC_TAG					// @...
	TOK_JDOC_NL						// Newlines are significant inside JavaDocs
	TOK_JDOC_LINE   				// everything else
	
	// This content is java-related
	TOK_JAVA_KEYWORD
	TOK_JAVA_PAREN_O
	TOK_JAVA_PAREN_X
	TOK_JAVA_CURLY_O
	TOK_JAVA_CURLO_X
	TOK_JAVA_CLASS
	TOK_JAVA_TYPE
	TOK_JAVA_IDENTIFIER
	TOK_JAVA_UNKNOWN
)

type Token struct {
	Type TokenType
	Lexeme string
}

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

func (this *Scanner) Peek() rune {
	ch := this.Next()
	this.Pos -= this.RuneWidth
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

func ScanBegin(scanner *Scanner) ScanFn {
	scanner.SkipWhitespace()
	
	// First, check if a JavaDoc is beginning
	if strings.HasPrefix(scanner.InputToEnd(), "/**") {
		return ScanJavadocStart
	}

	scanner.Emit(TOK_EOF)
	return nil
}

func ScanJavadocStart(scanner *Scanner) ScanFn {
	scanner.Pos += len("/**")
	scanner.Emit(TOK_JDOC_START)
	return ScanJavadoc
}

func ScanJavadocEnd(scanner *Scanner) ScanFn {
	scanner.Pos += len("*/")
	scanner.Emit(TOK_JDOC_END)
	return ScanBegin
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
			scanner.Emit(TOK_JDOC_TAG)
			return ScanJavadocLine
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