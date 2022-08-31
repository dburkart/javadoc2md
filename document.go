/*
 * Copyright (c) 2022, Dana Burkart <dana.burkart@gmail.com>
 *
 * SPDX-License-Identifier: BSD-2-Clause
 */

package main

import "fmt"

// The Document struct represents a single "document" emitted by the transpiler.
type Document struct {
	Address string
	Type string
	Blocks []Block
}

func (document *Document) AddBlock(block Block) {
	document.Blocks = append(document.Blocks, block)
}

func (document *Document) Printdbg() {
	fmt.Println("Document: ", document.Address)
	for _, v := range document.Blocks {
		v.Printdbg()
	}
}

// A single Javadoc "block", whether for a class or a function
type Block struct {
	Doc *Document
	Name string
	Description string
	Definition string
	Tags map[string]string
	Params map[string]string

	Attributes map[string]string
}

func (block *Block) Printdbg() {
	fmt.Println("Block: ", block.Name)
}

func MakeBlock() *Block {
	b := &Block{
		Name: "",
		Description: "",
		Definition: "",
		Tags: make(map[string]string),
		Params: make(map[string]string),
		Attributes: make(map[string]string),
	}

	return b
}

func MakeDocument(address string) *Document {
	d := &Document{
		Address: address,
		Blocks: []Block{},
	}

	return d
}
