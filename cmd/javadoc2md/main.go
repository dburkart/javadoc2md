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

	"github.com/dburkart/javadoc2md/internal/parser"
)

func discover(searchDirectory string, documents *[]*parser.Document) {
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

			s := parser.BeginScanning(fullPath, string(content[:]))
			d := parser.ParseDocument(s, fullPath)

			*documents = append(*documents, d)
		}
	}
}

func main() {
	var outputDirectory string
	var inputDirectory string

	var documents []*parser.Document

	flag.StringVar(&inputDirectory, "input", ".", "Input directory to transpile")
	flag.StringVar(&outputDirectory, "output", ".", "Output directory to receive markdown files")

	flag.Parse()

	discover(inputDirectory, &documents)

	options := parser.VisitorConfigOptions{
		OutputDirectory: outputDirectory,
	}

	parser.VisitDocuments(&options, &documents)
}
