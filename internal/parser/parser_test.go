/*
 * Copyright (c) 2022, Dana Burkart <dana.burkart@gmail.com>
 *
 * SPDX-License-Identifier: BSD-2-Clause
 */

package parser

import (
	"testing"
)

func TestSimpleClass(t *testing.T) {
	input := `
/**
 * This is a Simple Class
 */
public class SimpleClass {

	/**
     * Adds two numbers together
	 */
	public int add(int a, int b);
}`
	s := BeginScanningJavaCode("Test Simple Class", input)
	d := ParseDocument(s, "foo/bar/baz")

	if len(d.Blocks) != 2 {
		t.Errorf("expected 2 block")
	}

	if d.Blocks[0].Name != "SimpleClass" {
		t.Errorf("got class name of %s, wanted SimpleClass", d.Blocks[0].Name)
	}

	if len(d.Blocks[1].Arguments) != 2 {
		t.Errorf("got %d arguments, wanted 2", len(d.Blocks[1].Arguments))
	}
}
