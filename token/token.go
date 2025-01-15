// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package token

type Pos int

// IsValid reports whether the position is valid.
func (p Pos) IsValid() bool {
	return p != NoPos
}

// The zero value for Pos is NoPos; there is no file and line information
// associated with it, and NoPos.IsValid() is false. NoPos is always
// smaller than any other Pos value. The corresponding Position value
// for NoPos is the zero value for Position.
const NoPos Pos = 0

type Token string

const (
	ADD  = "+"
	SUB  = "-"
	MUL  = "*"
	DIV  = "/"
	IDIV = "\\"
	MOD  = "Mod"
	EXP  = "^"

	IS = "Is"

	BITAND = "&"

	NOT = "Not"
	AND = "And"
	OR  = "Or"
	XOR = "Xor"
	EQV = "Eqv"
	IMP = "Imp"

	EQ  = "="
	NEQ = "<>"

	LT = "<"
	GT = ">"

	LT_ASSIGN = "<="
	GT_ASSIGN = ">="

	APOSTROPHE = "'"

	FALSE   = "False"
	TRUE    = "True"
	NOTHING = "Nothing"

	BYVAL = "ByVal"
	BYREF = "ByRef"

	SET   = "Set"
	CONST = "Const"

	// Data Types

	EMPTY    = "Empty"
	NULL     = "Null"
	BOOLEAN  = "Boolean"
	BYTE     = "Byte"
	INTEGER  = "Integer"
	CURRENCY = "Currency"
	LONG     = "Long"
	SINGLE   = "Single"
	DOUBLE   = "Double"
	DATE     = "Date"
	STRING   = "String"
	OBJECT   = "Object"
	ERROR    = "Error"
)
