/*
 * Copyright (c) 2022, Dana Burkart <dana.burkart@gmail.com>
 *
 * SPDX-License-Identifier: BSD-2-Clause
 */

package main

import (
	"flag"
	"fmt"
	"io/ioutil"

	"github.com/dburkart/javadoc2md/internal/parser"
	"github.com/dburkart/javadoc2md/internal/util"
)

func main() {
	var outputDirectory string
	var inputDirectory string

	var documents []*parser.Document

	flag.StringVar(&inputDirectory, "input", ".", "Input directory to transpile")
	flag.StringVar(&outputDirectory, "output", ".", "Output directory to receive markdown files")

	flag.Parse()

	ctx := util.FileSearch(inputDirectory)

	for {
		fileToParse, ok := <- ctx.Files

		if !ok {
			break
		}

		content, err := ioutil.ReadFile(fileToParse)
		if err != nil {
			fmt.Println(err)
		}

		fmt.Println("Parsing", fileToParse)
		s := parser.BeginScanning(fileToParse, string(content[:]))
		d := parser.ParseDocument(s, fileToParse)

		documents = append(documents, d)
	}

	options := parser.VisitorConfigOptions{
		OutputDirectory: outputDirectory,
	}

	parser.VisitDocuments(&options, &documents)
}
