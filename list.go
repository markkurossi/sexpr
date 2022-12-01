//
// Copyright (c) 2022 Markku Rossi
//
// All rights reserved.
//

package scheme

import (
	"errors"
	"fmt"
	"strings"
)

var (
	_ Pair  = &PlainPair{}
	_ Pair  = &LocationPair{}
	_ Value = &PlainPair{}
)

// Pair implements a Scheme pair.
type Pair interface {
	Locator
	Car() Value
	Cdr() Value
	SetCar(v Value) error
	SetCdr(v Value) error
	Scheme() string
	Equal(o Value) bool
}

// PlainPair implements a Scheme pair with car and cdr values.
type PlainPair struct {
	car Value
	cdr Value
}

// NewPair creates a new pair with the car and cdr values.
func NewPair(car, cdr Value) Pair {
	return &PlainPair{
		car: car,
		cdr: cdr,
	}
}

// Location implements the Locator interface.
func (pair *PlainPair) Location() Point {
	return Point{}
}

// Car returns the pair's car value.
func (pair *PlainPair) Car() Value {
	return pair.car
}

// Cdr returns the pair's cdr value.
func (pair *PlainPair) Cdr() Value {
	return pair.cdr
}

// SetCar sets the pair's car value.
func (pair *PlainPair) SetCar(v Value) error {
	pair.car = v
	return nil
}

// SetCdr sets the pair's cdr value.
func (pair *PlainPair) SetCdr(v Value) error {
	pair.cdr = v
	return nil
}

// Scheme returns the value as a Scheme string.
func (pair *PlainPair) Scheme() string {
	return pair.String()
}

// Equal tests if the argument value is equal to this value.
func (pair *PlainPair) Equal(o Value) bool {
	ov, ok := o.(Pair)
	return ok && Equal(pair.car, ov.Car()) && Equal(pair.cdr, ov.Cdr())
}

func (pair *PlainPair) String() string {
	var str strings.Builder
	str.WriteRune('(')

	i := Pair(pair)
	first := true
loop:
	for {
		if first {
			first = false
		} else {
			str.WriteRune(' ')
		}
		if i.Car() == nil {
			str.WriteString("nil")
		} else {
			str.WriteString(i.Car().Scheme())
		}
		switch cdr := i.Cdr().(type) {
		case Pair:
			i = cdr

		case nil:
			break loop

		default:
			str.WriteString(" . ")
			str.WriteString(fmt.Sprintf("%v", cdr))
			break loop
		}
	}
	str.WriteRune(')')

	return str.String()
}

// LocationPair implements a Scheme pair with location information.
type LocationPair struct {
	Point
	PlainPair
}

// Location implements the Locator interface.
func (pair *LocationPair) Location() Point {
	return pair.Point
}

func (pair *LocationPair) String() string {
	return pair.PlainPair.String()
}

// NewLocationPair creates a new pair with the car and cdr values and
// location information.
func NewLocationPair(point Point, car, cdr Value) Pair {
	return &LocationPair{
		Point: point,
		PlainPair: PlainPair{
			car: car,
			cdr: cdr,
		},
	}
}

// ErrorInvalidList is used to indicate when a malformed or otherwise
// invalid list is passed to list functions.
var ErrorInvalidList = errors.New("invalid list")

// ListLength check if the argument value is a valid Scheme list.
func ListLength(list Value) (int, bool) {
	var count int

	err := Map(func(idx int, v Value) error { count++; return nil }, list)
	if err != nil {
		return 0, false
	}
	return count, true
}

// Map maps function for each element of the list. The function
// returns nil if the argument list is a list and map functions
// returns nil for each of its element.
func Map(f func(idx int, v Value) error, list Value) error {
	if list == nil {
		return nil
	}
	pair, ok := list.(Pair)
	if !ok {
		return ErrorInvalidList
	}

	for idx := 0; pair != nil; idx++ {
		if err := f(idx, pair.Car()); err != nil {
			point := pair.Location()
			if point.Undefined() {
				return err
			}
			return fmt.Errorf("%s: %v", point, err)
		}
		switch cdr := pair.Cdr().(type) {
		case Pair:
			pair = cdr

		case nil:
			pair = nil

		default:
			return ErrorInvalidList
		}
	}
	return nil
}

// Car returns the car element of the pair.
func Car(pair Value, ok bool) (Value, bool) {
	if !ok {
		return pair, false
	}
	p, ok := pair.(Pair)
	if !ok {
		return pair, false
	}
	return p.Car(), true
}

// Cdr returns the cdr element of the cons cell.
func Cdr(pair Value, ok bool) (Value, bool) {
	if !ok {
		return pair, false
	}
	p, ok := pair.(Pair)
	if !ok {
		return pair, false
	}
	return p.Cdr(), true
}

var listBuiltins = []Builtin{
	{
		Name: "pair?",
		Args: []string{"obj"},
		Native: func(scm *Scheme, args []Value) (Value, error) {
			_, ok := args[0].(Pair)
			return Boolean(ok), nil
		},
	},
	{
		Name: "cons",
		Args: []string{"obj1", "obj2"},
		Native: func(scm *Scheme, args []Value) (Value, error) {
			return NewPair(args[0], args[1]), nil
		},
	},
	{
		Name: "car",
		Args: []string{"pair"},
		Native: func(scm *Scheme, args []Value) (Value, error) {
			pair, ok := args[0].(Pair)
			if !ok {
				return nil, fmt.Errorf("car: not a pair: %v", args[0])
			}
			return pair.Car(), nil
		},
	},
	{
		Name: "cdr",
		Args: []string{"pair"},
		Native: func(scm *Scheme, args []Value) (Value, error) {
			pair, ok := args[0].(Pair)
			if !ok {
				return nil, fmt.Errorf("cdr: not a pair: %v", args[0])
			}
			return pair.Cdr(), nil
		},
	},
	{
		Name: "set-car!",
		Args: []string{"pair", "obj"},
		Native: func(scm *Scheme, args []Value) (Value, error) {
			pair, ok := args[0].(Pair)
			if !ok {
				return nil, fmt.Errorf("set-car!: not a pair: %v", args[0])
			}
			err := pair.SetCar(args[1])
			if err != nil {
				return nil, err
			}
			return nil, nil
		},
	},
	{
		Name: "set-cdr!",
		Args: []string{"pair", "obj"},
		Native: func(scm *Scheme, args []Value) (Value, error) {
			pair, ok := args[0].(Pair)
			if !ok {
				return nil, fmt.Errorf("set-cdr!: not a pair: %v", args[0])
			}
			err := pair.SetCdr(args[1])
			if err != nil {
				return nil, err
			}
			return nil, nil
		},
	},
	{
		Name: "null?",
		Args: []string{"obj"},
		Native: func(scm *Scheme, args []Value) (Value, error) {
			return Boolean(args[0] == nil), nil
		},
	},
	{
		Name: "list?",
		Args: []string{"obj"},
		Native: func(scm *Scheme, args []Value) (Value, error) {
			_, ok := ListLength(args[0])
			return Boolean(ok), nil
		},
	},
	{
		Name: "list",
		Args: []string{"[obj]..."},
		Native: func(scm *Scheme, args []Value) (Value, error) {
			var result Pair
			for i := len(args) - 1; i >= 0; i-- {
				result = NewPair(args[i], result)
			}
			return result, nil
		},
	},
	{
		Name: "length",
		Args: []string{"obj"},
		Native: func(scm *Scheme, args []Value) (Value, error) {
			length, ok := ListLength(args[0])
			if !ok {
				return nil, fmt.Errorf("length: not a list: %v", args[0])
			}
			return NewNumber(0, length), nil
		},
	},
	// XXX append
	{
		Name: "reverse",
		Args: []string{"list"},
		Native: func(scm *Scheme, args []Value) (Value, error) {
			var result Pair

			err := Map(func(idx int, v Value) error {
				result = NewPair(v, result)
				return nil
			}, args[0])
			if err != nil {
				return nil, fmt.Errorf("reverse: %v", err)
			}
			return result, nil
		},
	},
	{
		Name: "list-tail",
		Args: []string{"list", "k"},
		Native: func(scm *Scheme, args []Value) (Value, error) {
			kn, ok := args[1].(Number)
			if !ok {
				return nil, fmt.Errorf("list-tail: invalid index: %v", args[1])
			}
			k := int(kn.Int64())
			pair := args[0]
			var i int

			for i = 0; i < k && pair != nil; i++ {
				pair, ok = Cdr(pair, true)
				if !ok {
					return nil, fmt.Errorf("list-tail: invalid list: %v", pair)
				}
			}
			if i < k {
				return nil,
					fmt.Errorf("list-tail: index %v out of range for list %v",
						k, args[0])
			}

			return pair, nil
		},
	},
	{
		Name: "list-ref",
		Args: []string{"list", "k"},
		Native: func(scm *Scheme, args []Value) (Value, error) {
			kn, ok := args[1].(Number)
			if !ok {
				return nil, fmt.Errorf("list-ref: invalid index: %v", args[1])
			}
			k := int(kn.Int64())
			pair := args[0]
			var i int

			for i = 0; i < k && pair != nil; i++ {
				pair, ok = Cdr(pair, true)
				if !ok {
					return nil, fmt.Errorf("list-ref: invalid list: %v", pair)
				}
			}
			if pair == nil {
				return nil,
					fmt.Errorf("list-ref: index %v out of range for list %v",
						k, args[0])
			}

			result, ok := Car(pair, true)
			if !ok {
				return nil, fmt.Errorf("list-ref: invalid list: %v", pair)
			}

			return result, nil
		},
	},
}
