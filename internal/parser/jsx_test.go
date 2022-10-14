/*
 * Copyright (c) 2022, Dana Burkart <dana.burkart@gmail.com>
 *
 * SPDX-License-Identifier: BSD-2-Clause
 */

package parser

import "testing"

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
