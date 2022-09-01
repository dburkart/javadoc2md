/*
 * Copyright (c) 2022, Dana Burkart <dana.burkart@gmail.com>
 *
 * SPDX-License-Identifier: BSD-2-Clause
 */

package parser

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
	TOK_JAVA_EQUAL
	TOK_JAVA_STRING
	TOK_JAVA_OPERATOR
	TOK_JAVA_BRACKET_O
	TOK_JAVA_BRACKET_X
	TOK_JAVA_IDENTIFIER
	TOK_JAVA_NUMERIC
	TOK_JAVA_ANNOTATION
	TOK_JAVA_COMMENT_O
	TOK_JAVA_COMMENT_X
	TOK_JAVA_OTHER
)

type Token struct {
	Type TokenType
	Lexeme string
}
