/*
 * Copyright (c) 2022, Dana Burkart <dana.burkart@gmail.com>
 *
 * SPDX-License-Identifier: BSD-2-Clause
 */

package main

import (
	"io/ioutil"
	"fmt"
	"os"
)

func main() {
	f := os.Args[1]

	content, err := ioutil.ReadFile(f)
	if err != nil {
		fmt.Println(err)
	}

	s := BeginScanning(f, string(content[:]))

	d := ParseDocument(s, f)

	symbolVisitor := SymbolVisitor{Symbols: make(map[string]string)}
	symbolVisitor.visit(d)

	visitor := MarkdownVisitor{}
	visitor.visit(d)
}
