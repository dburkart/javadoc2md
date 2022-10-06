/*
 * Copyright (c) 2022, Dana Burkart <dana.burkart@gmail.com>
 *
 * SPDX-License-Identifier: BSD-2-Clause
 */

package parser

import "fmt"

// The Document struct represents a single "document" emitted by the transpiler.
type Document struct {
	Address string
	Package string
	Blocks  []Block
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

type ArgPair struct {
	Type string
	Name string
}

// A single Javadoc "block", whether for a class or a function
type Block struct {
	Doc        *Document
	Name       string
	Type       SymbolType
	Arguments  []ArgPair
	Text       Text
	Definition string
	Tags       map[string]Text
	Params     map[string]Text
	Attributes map[string]string
}

func (block *Block) Printdbg() {
	fmt.Println("Block: ", block.Name)
}

func MakeBlock() *Block {
	b := &Block{
		Name:       "",
		Text:       []Token{},
		Definition: "",
		Arguments:  []ArgPair{},
		Tags:       make(map[string]Text),
		Params:     make(map[string]Text),
		Attributes: make(map[string]string),
	}

	return b
}

func MakeDocument(address string) *Document {
	d := &Document{
		Address: address,
		Blocks:  []Block{},
	}

	return d
}
