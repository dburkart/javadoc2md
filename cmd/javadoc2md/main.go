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
	"sync"

	"github.com/dburkart/javadoc2md/internal/logger"
	"github.com/dburkart/javadoc2md/internal/parser"
	"github.com/dburkart/javadoc2md/internal/util"
)

func main() {
	var outputDirectory string
	var inputDirectory string
	var skipPrivateDefs bool

	flag.StringVar(&inputDirectory, "input", ".", "Input directory to transpile")
	flag.StringVar(&outputDirectory, "output", ".", "Output directory to receive markdown files")
	skipPrivateDefs = *flag.Bool("skip-private", false, "Skip private definitions")

	flag.Parse()

	logger.Initialize()

	ctx := util.FileSearch(inputDirectory)
	documents := make(chan *parser.Document)
	var wg sync.WaitGroup

	for {
		fileToParse, ok := <-ctx.Files

		if !ok {
			break
		}

		content, err := ioutil.ReadFile(fileToParse)
		if err != nil {
			fmt.Println(err)
		}

		go func() {
			s := parser.BeginScanningJavaCode(fileToParse, string(content[:]))
			d := parser.ParseDocument(s, fileToParse)

			documents <- d
			wg.Done()
		}()
		wg.Add(1)
	}

	go func() {
		wg.Wait()
		close(documents)
	}()

	options := parser.VisitorConfigOptions{
		OutputDirectory: outputDirectory,
		SkipPrivateDefs: skipPrivateDefs,
	}

	parser.VisitDocuments(&options, documents)
}
