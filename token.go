/*
 * Copyright (c) 2022, Dana Burkart <dana.burkart@gmail.com>
 *
 * SPDX-License-Identifier: BSD-2-Clause
 */

package main

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
	TOK_JAVA_COMMA
	TOK_JAVA_IDENTIFIER
)

type Token struct {
	Type TokenType
	Lexeme string
}
