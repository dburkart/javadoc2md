/*
 * Copyright (c) 2022, Dana Burkart <dana.burkart@gmail.com>
 *
 * SPDX-License-Identifier: BSD-2-Clause
 */

package parser

import "unicode"

type jsxTag struct {
	index int
	tag   string
}

func (j *jsxTag) tagType() string {
	start, end := 1, 0

	for i, value := range j.tag {
		if unicode.IsSpace(value) {
			end = i - 1
			break
		}

		if value == '/' && i == len(j.tag)-1 {
			end = i - 1
			break
		} else if value == '/' {
			start = i + 1
		}

		if value == '>' {
			end = i
			break
		}
	}

	return j.tag[start:end]
}

func (j *jsxTag) close() string {
	isClosed := false
	for i, value := range j.tag {
		if value == '>' && j.tag[i-1] == '/' {
			isClosed = true
		}
	}

	if !isClosed {
		j.tag = j.tag[:len(j.tag)-1] + "/>"
	}

	return j.tag
}
