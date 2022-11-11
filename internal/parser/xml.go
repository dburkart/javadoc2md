/*
 * Copyright (c) 2022, Dana Burkart <dana.burkart@gmail.com>
 *
 * SPDX-License-Identifier: BSD-2-Clause
 */

package parser

import "unicode"

type XMLTag struct {
	Index int
	Tag   string
}

func (j *XMLTag) Type() string {
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

func (j *XMLTag) Attributes() map[string]string {
	attributes := make(map[string]string)

	tagType := j.Type()
	key := ""
	value := ""
	inValue := false
	quoteChar := '"'

	for i := 0; i < len(j.Tag); i++ {
		ch := rune(j.Tag[i])
		if (unicode.IsSpace(ch) || ch == '/' || ch == '>') && !inValue {
			if key != tagType {
				attributes[key] = value
			}
			key = ""
			value = ""
			continue
		}

		if inValue {
			if ch == quoteChar {
				inValue = false
				continue
			}
			value = value + string(ch)
		} else {
			if ch == '=' {
				quoteChar = rune(j.Tag[i+1])
				i++
				inValue = true
				continue
			}

			key = key + string(ch)
		}
	}

	return attributes
}

func (j *XMLTag) Close() string {
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
