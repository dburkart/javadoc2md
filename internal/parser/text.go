/*
 * Copyright (c) 2022, Dana Burkart <dana.burkart@gmail.com>
 *
 * SPDX-License-Identifier: BSD-2-Clause
 */

package parser

// Really, we should be building an AST since Javadoc can have parameters
// virtually anywhere, but storing token lists in Blocks is simpler for now.
type Text []Token
