//
// Copyright (c) 2022-2023 Markku Rossi
//
// All rights reserved.
//

package scheme

import (
	"fmt"
	"strings"

	"github.com/markkurossi/scheme/types"
)

var (
	_ Value = &BigInt{}
	_ Value = &Bytevector{}
	_ Value = &Frame{}
	_ Value = &Identifier{}
	_ Value = &Lambda{}
	_ Value = &PlainPair{}
	_ Value = &Port{}
	_ Value = &Vector{}
	_ Value = Boolean(true)
	_ Value = Character('@')
	_ Value = Int(0)
	_ Value = Keyword(0)
	_ Value = String("string")
)

// Value implements a Scheme value.
type Value interface {
	Scheme() string
	Eq(o Value) bool
	Equal(o Value) bool
	Type() *types.Type
}

// ToString returns a display representation of the value.
func ToString(v Value) string {
	if v == nil {
		return "'()"
	}
	return fmt.Sprintf("%v", v)
}

// ToScheme returns a Scheme representation of the value.
func ToScheme(v Value) string {
	if v == nil {
		return "'()"
	}
	return v.Scheme()
}

// Flags define symbol flags.
type Flags int

// Symbol flags.
const (
	FlagDefined Flags = 1 << iota
	FlagConst
)

func (f Flags) String() string {
	var result string
	if f&FlagDefined != 0 {
		result += " defined"
	}
	if f&FlagConst != 0 {
		result += " const"
	}
	return strings.TrimSpace(result)
}

// Identifier implements identifier values.
type Identifier struct {
	Name       string
	Point      Point
	GlobalType *types.Type
	Global     Value
	Flags      Flags
}

// Scheme returns the value as a Scheme string.
func (v *Identifier) Scheme() string {
	return v.String()
}

// Eq tests if the argument value is eq? to this value.
func (v *Identifier) Eq(o Value) bool {
	return v.Equal(o)
}

// Equal tests if the argument value is equal to this value.
func (v *Identifier) Equal(o Value) bool {
	ov, ok := o.(*Identifier)
	return ok && v.Name == ov.Name
}

// Type implements Value.Type.
func (v *Identifier) Type() *types.Type {
	return types.Symbol
}

func (v *Identifier) String() string {
	return v.Name
}
