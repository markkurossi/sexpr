//
// Copyright (c) 2023 Markku Rossi
//
// All rights reserved.
//

package scheme

import (
	"fmt"
	"io"
	"os"
	"path"

	"github.com/markkurossi/scheme/types"
)

var loadBuiltins = []Builtin{
	{
		Name:   "scheme::load",
		Args:   []string{"caller<string>", "filename<string>"},
		Return: types.Any,
		Native: func(scm *Scheme, args []Value) (Value, error) {
			caller, ok := args[0].(String)
			if !ok {
				return nil, fmt.Errorf("invalid caller: %v", args[0])
			}
			f, ok := args[1].(String)
			if !ok {
				return nil, fmt.Errorf("invalid filename: %v", args[1])
			}
			file := string(f)
			if !path.IsAbs(file) {
				file = path.Join(path.Dir(string(caller)), file)
			}
			if scm.Params.Verbose {
				fmt.Printf("load: %v\n", file)
			}
			return scm.LoadFile(file)
		},
	},
	{
		Name: "scheme::stack-trace",
		Return: &types.Type{
			Enum: types.EnumList,
			Element: &types.Type{
				Enum: types.EnumPair,
				Car:  types.String,
				Cdr:  types.InexactInteger,
			},
		},
		Native: func(scm *Scheme, args []Value) (Value, error) {
			stack := scm.StackTrace()

			var result, tail Pair

			for _, frame := range stack {
				p := NewPair(
					NewPair(String(frame.Source),
						NewNumber(0, frame.Line)),
					nil)
				if tail == nil {
					result = p
				} else {
					tail.SetCdr(p)
				}
				tail = p
			}

			return result, nil
		},
	},
}

// LoadFile loads and compiles the file.
func (scm *Scheme) LoadFile(file string) (Value, error) {
	in, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer in.Close()
	return scm.Load(file, in)
}

// Load loads and compiles the input.
func (scm *Scheme) Load(source string, in io.Reader) (Value, error) {
	c := NewCompiler(scm)

	library, err := c.Compile(source, in)
	if err != nil {
		return nil, err
	}
	if false {
		fmt.Printf("Code:\n")
		for _, c := range library.Init {
			fmt.Printf("%s\n", c)
		}
	}

	return NewPair(&Identifier{Name: "library"},
		NewPair(library.Name,
			NewPair(library.Exports,
				NewPair(library.Imports,
					NewPair(
						&Lambda{
							Impl: &LambdaImpl{
								Return:   types.Any,
								Source:   library.Source,
								Code:     library.Init,
								PCMap:    library.PCMap,
								Captures: true,
							},
						},
						nil))))), nil
}
