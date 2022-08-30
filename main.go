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

func main() {
	var outputDirectory string
	var inputDirectory string

	var documents []*Document

	flag.StringVar(&inputDirectory, "input", ".", "Input directory to transpile")
	flag.StringVar(&outputDirectory, "output", ".", "Output directory to receive markdown files")

	flag.Parse()

	files, err := ioutil.ReadDir(inputDirectory)

	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		fullPath := filepath.Join(inputDirectory, file.Name())

		if file.IsDir() {
			fmt.Println("Skipping directory", fullPath)
			continue
		}

		if strings.HasSuffix(file.Name(), ".java") {
			fmt.Println("Parsing", fullPath)
			content, err := ioutil.ReadFile(fullPath)
			if err != nil {
				fmt.Println(err)
			}

			s := BeginScanning(fullPath, string(content[:]))
			d := ParseDocument(s, fullPath)

			documents = append(documents, d)
		}
	}

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
