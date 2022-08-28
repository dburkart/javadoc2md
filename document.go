/*
 * Copyright (c) 2022, Dana Burkart <dana.burkart@gmail.com>
 *
 * SPDX-License-Identifier: BSD-2-Clause
 */

package main

// The Document struct represents a single "document" emitted by the transpiler.
type Document struct {
	Address string
	Blocks []Block
}

func (document *Document) AddBlock(block Block) {
	document.Blocks = append(document.Blocks, block)
}

// A single Javadoc "block", whether for a class or a function
type Block struct {
	Name string
	Description string
	Tags map[string]string
}

func MakeBlock(name string) *Block {
	b := &Block{
		Name: name,
		Description: "",
		Tags: make(map[string]string),
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
