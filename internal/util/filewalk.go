/*
 * Copyright (c) 2022, Dana Burkart <dana.burkart@gmail.com>
 *
 * SPDX-License-Identifier: BSD-2-Clause
 */

package util

import (
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"
)

type SearchContext struct {
	Root string
	Files chan string
}

func (ctx *SearchContext)discover(directory string) {
	files, err := ioutil.ReadDir(directory)

	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		fullPath := filepath.Join(directory, file.Name())

		if file.IsDir() {
			ctx.discover(fullPath)
			continue
		}

		if strings.HasSuffix(file.Name(), ".java") {
			ctx.Files <- fullPath
		}
	}
}

func FileSearch(root string) *SearchContext {
	s := &SearchContext{
		Root: root,
		Files: make(chan string, 3),
	}

	go func() {
		s.discover(root)
		close(s.Files)
	}()

	return s
}
