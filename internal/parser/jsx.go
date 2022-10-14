/*
 * Copyright (c) 2022, Dana Burkart <dana.burkart@gmail.com>
 *
 * SPDX-License-Identifier: BSD-2-Clause
 */

package parser

import "unicode"

type JSXTag struct {
	Index int
	Tag   string
}

func (j *JSXTag) Type() string {
	start, end := 1, 0

	for i, value := range j.Tag {
		if unicode.IsSpace(value) {
			end = i
			break
		}

		if value == '/' && i == len(j.Tag)-2 {
			end = i
			break
		} else if value == '/' {
			start = i + 1
		}

		if value == '>' {
			end = i
			break
		}
	}

	return j.Tag[start:end]
}

func (j *JSXTag) Close() string {
	isClosed := false
	for i, value := range j.Tag {
		if value == '>' && j.Tag[i-1] == '/' {
			isClosed = true
		}
	}

	if !isClosed {
		j.Tag = j.Tag[:len(j.Tag)-1] + "/>"
	}

	return j.Tag
}
