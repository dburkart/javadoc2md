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

    "github.com/dburkart/javadoc2md/internal/parser"
    "github.com/dburkart/javadoc2md/internal/util"
)

func main() {
    var outputDirectory string
    var inputDirectory string

    flag.StringVar(&inputDirectory, "input", ".", "Input directory to transpile")
    flag.StringVar(&outputDirectory, "output", ".", "Output directory to receive markdown files")

    flag.Parse()

    ctx := util.FileSearch(inputDirectory)
    documents := make(chan *parser.Document)
    var wg sync.WaitGroup

    for {
        fileToParse, ok := <- ctx.Files

        if !ok {
            break
        }

        content, err := ioutil.ReadFile(fileToParse)
        if err != nil {
            fmt.Println(err)
        }

        go func() {
            s := parser.BeginScanning(fileToParse, string(content[:]))
            d := parser.ParseDocument(s, fileToParse)

            documents <- d
            fmt.Println("Parsed", fileToParse)
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
    }

    parser.VisitDocuments(&options, documents)
}
