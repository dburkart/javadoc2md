/*
 * Copyright (c) 2022, Dana Burkart <dana.burkart@gmail.com>
 *
 * SPDX-License-Identifier: BSD-2-Clause
 */

package parser

import (
	"testing"
)

func TestTagTypeSimple(t *testing.T) {
	tag := JSXTag{Index: 0, Tag: "<b>"}
	if tag.Type() != "b" {
		t.Errorf("got %q, wanted 'b'", tag.Type())
	}
}

func TestTagTypeAHref(t *testing.T) {
	tag := JSXTag{Index: 0, Tag: "<a href='#'>"}
	if tag.Type() != "a" {
		t.Errorf("got %q, wanted 'a'", tag.Type())
	}
}

func TestSelfClosedTag(t *testing.T) {
	tag := JSXTag{Index: 0, Tag: "<p/>"}
	if tag.Type() != "p" {
		t.Errorf("got %q, wanted 'p'", tag.Type())
	}
}

func TestClosingTag(t *testing.T) {
	tag := JSXTag{Index: 0, Tag: "</b>"}
	if tag.Type() != "b" {
		t.Errorf("got %q, wanted 'b'", tag.Type())
	}
}

func TestAutoCloseTags(t *testing.T) {
	tag := JSXTag{Index: 0, Tag: "<p>"}
	if tag.Close() != "<p/>" {
		t.Errorf("got %q, wanted '<p/>'", tag.Tag)
	}

	// Test tag which is already closed
	if tag.Close() != "<p/>" {
		t.Errorf("got %q, wanted '<p/>'", tag.Tag)
	}
}

func TestSimpleAttribute(t *testing.T) {
	tag := JSXTag{Index: 0, Tag: "<a href='#'>"}
	attributes := tag.Attributes()

	if attributes["href"] != "#" {
		t.Errorf("got attribute '%q', wanted '#'", attributes["href"])
	}
}

func TestMultipleAttributes(t *testing.T) {
	tag := JSXTag{Index: 0, Tag: "<person name='bob' gender='male'>"}
	attributes := tag.Attributes()

	if attributes["name"] != "bob" {
		t.Errorf("got attribute '%q' for name, wanted 'bob'", attributes["name"])
	}

	if attributes["gender"] != "male" {
		t.Errorf("got attribute '%q' for gender, wanted 'male'", attributes["gender"])
	}
}
