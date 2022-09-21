/*
 * Copyright (c) 2022, Dana Burkart <dana.burkart@gmail.com>
 *
 * SPDX-License-Identifier: BSD-2-Clause
 */

package parser

type SymbolType int

const (
	SYM_TYPE_INVALID SymbolType = iota
	SYM_TYPE_CLASS
	SYM_TYPE_INTERFACE
	SYM_TYPE_ENUM
	SYM_TYPE_METHOD
	SYM_TYPE_FIELD
)

type Symbol struct {
	Type     SymbolType
	Name     string // Short name
	Package  string
	Parent   string // Fields, methods, inner classes
	Location string
}
