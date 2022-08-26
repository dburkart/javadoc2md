/*
 * Copyright (c) 2022, Dana Burkart <dana.burkart@gmail.com>
 *
 * SPDX-License-Identifier: BSD-2-Clause
 */
 
package main
 
import "testing"
 
func SetupWithState(input string, state ScanFn) *Scanner {
	s := &Scanner{
		Name: "Test",
		Input: input,
		State: state,
		Tokens: make(chan Token, 3),
	}
	
	return s
}

func TestScanJavadocTag(t *testing.T) {
	s := SetupWithState("@tag and other stuff", ScanJavadocTag)
	
	s.State(s)
	token := <- s.Tokens
	
	if token.Type != TOK_JDOC_TAG {
		t.Errorf("got %q, wanted %q", token.Type, TOK_JDOC_TAG)
	}
	
	if token.Lexeme != "@tag" {
		t.Errorf("got %q, wanted @tag", token.Lexeme)
	}
}