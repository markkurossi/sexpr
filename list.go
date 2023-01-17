//
// Copyright (c) 2022-2023 Markku Rossi
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
	Eq(o Value) bool
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

// From returns pair's start location.
func (pair *PlainPair) From() Point {
	return Point{}
}

// To returns pair's end location.
func (pair *PlainPair) To() Point {
	return Point{}
}

// SetTo sets pair's end location.
func (pair *PlainPair) SetTo(p Point) {
}

// Errorf returns an error with the pair's location.
func (pair *PlainPair) Errorf(format string, a ...interface{}) error {
	return fmt.Errorf(format, a...)
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

// Eq tests if the argument value is eq? to this value.
func (pair *PlainPair) Eq(o Value) bool {
	ov, ok := o.(*PlainPair)
	return ok && pair == ov
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
	from Point
	to   Point
	PlainPair
}

// NewLocationPair creates a new pair with the car and cdr values and
// location information.
func NewLocationPair(from, to Point, car, cdr Value) Pair {
	return &LocationPair{
		from: from,
		to:   to,
		PlainPair: PlainPair{
			car: car,
			cdr: cdr,
		},
	}
}

// From returns pair's start location.
func (pair *LocationPair) From() Point {
	return pair.from
}

// To returns pair's end location.
func (pair *LocationPair) To() Point {
	return pair.to
}

// SetTo sets pair's end location.
func (pair *LocationPair) SetTo(p Point) {
	pair.to = p
}

// Eq tests if the argument value is eq? to this value.
func (pair *LocationPair) Eq(o Value) bool {
	ov, ok := o.(*LocationPair)
	return ok && pair == ov
}

// Errorf returns an error with the pair's location.
func (pair *LocationPair) Errorf(format string, a ...interface{}) error {
	msg := fmt.Sprintf(format, a...)
	return fmt.Errorf("%s: %s", pair.from, msg)
}

func (pair *LocationPair) String() string {
	return pair.PlainPair.String()
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

// ListPairs returns the list pairs as []Pair.
func ListPairs(list Value) ([]Pair, bool) {
	var result []Pair

	err := MapPairs(func(idx int, p Pair) error {
		result = append(result, p)
		return nil
	}, list)
	if err != nil {
		return nil, false
	}
	return result, true
}

// ListValues returns the list values as []Value.
func ListValues(list Value) ([]Value, bool) {
	var result []Value

	err := Map(func(idx int, v Value) error {
		result = append(result, v)
		return nil
	}, list)
	if err != nil {
		return nil, false
	}
	return result, true
}

// Map maps function for each element of the list. The function
// returns nil if the argument list is a list and map functions
// returns nil for each of its element.
func Map(f func(idx int, v Value) error, list Value) error {
	return MapPairs(func(idx int, p Pair) error {
		return f(idx, p.Car())
	}, list)
}

// MapPairs maps function for each pair of the list. The function
// returns nil if the argument list is a list and map functions
// returns nil for each of its element.
func MapPairs(f func(idx int, p Pair) error, list Value) error {
	if list == nil {
		return nil
	}
	pair, ok := list.(Pair)
	if !ok {
		return ErrorInvalidList
	}

	for idx := 0; pair != nil; idx++ {
		if err := f(idx, pair); err != nil {
			point := pair.From()
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
		Native: func(scm *Scheme, l *Lambda, args []Value) (Value, error) {
			_, ok := args[0].(Pair)
			return Boolean(ok), nil
		},
	},
	{
		Name: "cons",
		Args: []string{"obj1", "obj2"},
		Native: func(scm *Scheme, l *Lambda, args []Value) (Value, error) {
			return NewPair(args[0], args[1]), nil
		},
	},
	{
		Name: "car",
		Args: []string{"pair"},
		Native: func(scm *Scheme, l *Lambda, args []Value) (Value, error) {
			pair, ok := args[0].(Pair)
			if !ok {
				return nil, l.Errorf("not a pair: %v", args[0])
			}
			return pair.Car(), nil
		},
	},
	{
		Name: "cdr",
		Args: []string{"pair"},
		Native: func(scm *Scheme, l *Lambda, args []Value) (Value, error) {
			pair, ok := args[0].(Pair)
			if !ok {
				return nil, l.Errorf("not a pair: %v", args[0])
			}
			return pair.Cdr(), nil
		},
	},
	{
		Name: "set-car!",
		Args: []string{"pair", "obj"},
		Native: func(scm *Scheme, l *Lambda, args []Value) (Value, error) {
			pair, ok := args[0].(Pair)
			if !ok {
				return nil, l.Errorf("not a pair: %v", args[0])
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
		Native: func(scm *Scheme, l *Lambda, args []Value) (Value, error) {
			pair, ok := args[0].(Pair)
			if !ok {
				return nil, l.Errorf("not a pair: %v", args[0])
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
		Native: func(scm *Scheme, l *Lambda, args []Value) (Value, error) {
			return Boolean(args[0] == nil), nil
		},
	},
}
