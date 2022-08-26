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

func TestScanJavadocStart(t *testing.T) {
	input := `/**
	* Test`
	
	s := SetupWithState(input, ScanJavadocStart)
	
	s.State(s)
	token := <- s.Tokens
	
	if token.Type != TOK_JDOC_START {
		t.Errorf("got %q, wanted %q", token.Type, TOK_JDOC_START)
	}
}

func TestScanJavadocEnd(t *testing.T) {
	input := "*/"
	
	s := SetupWithState(input, ScanJavadocEnd)
	
	s.State(s)
	token := <- s.Tokens
	
	if token.Type != TOK_JDOC_END {
		t.Errorf("got %q, wanted %q", token.Type, TOK_JDOC_END)
	}
}

func TestScanJavadocLine(t *testing.T) {
	input := `This is a line assumed to be in a Javadoc
	* This is the next line`
	
	s := SetupWithState(input, ScanJavadocLine)
	
	s.State(s)
	token := <- s.Tokens
	
	if token.Type != TOK_JDOC_LINE {
		t.Errorf("got %q, wanted %q", token.Type, TOK_JDOC_LINE)
	}
	
	if token.Lexeme != "This is a line assumed to be in a Javadoc" {
		t.Errorf("got %q, wanted 'This is a line assumed to be in a Javadoc'", token.Lexeme)
	}
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