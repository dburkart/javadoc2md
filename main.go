/*
 * Copyright (c) 2022, Dana Burkart <dana.burkart@gmail.com>
 *
 * SPDX-License-Identifier: BSD-2-Clause
 */

package main

import (
	"io/ioutil"
	"flag"
	"fmt"
	"log"
	"path/filepath"
	"strings"
)

func discover(searchDirectory string, documents *[]*Document) {
	files, err := ioutil.ReadDir(searchDirectory)

	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		fullPath := filepath.Join(searchDirectory, file.Name())

		if file.IsDir() {
			discover(fullPath, documents)
		}

		if strings.HasSuffix(file.Name(), ".java") {
			fmt.Println("Parsing", fullPath)
			content, err := ioutil.ReadFile(fullPath)
			if err != nil {
				fmt.Println(err)
			}

			s := BeginScanning(fullPath, string(content[:]))
			d := ParseDocument(s, fullPath)

			*documents = append(*documents, d)
		}
	}
}

func main() {
	var outputDirectory string
	var inputDirectory string

	var documents []*Document

	flag.StringVar(&inputDirectory, "input", ".", "Input directory to transpile")
	flag.StringVar(&outputDirectory, "output", ".", "Output directory to receive markdown files")

	flag.Parse()

	discover(inputDirectory, &documents)

	visitors := []Visitor{
		&SymbolVisitor{Symbols: make(map[string]string)},
		&MarkdownVisitor{OutputDirectory: outputDirectory},
	}

	for _, v := range visitors {
		for _, d := range documents {
			v.visit(d)
		}
	}
}
