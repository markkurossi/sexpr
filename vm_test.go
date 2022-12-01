//
// Copyright (c) 2022 Markku Rossi
//
// All rights reserved.
//

package scheme

import (
	"fmt"
	"strings"
	"testing"
)

var vmTests = []struct {
	i string
	v Value
	o string
}{
	{
		i: `(display "The length of \"Hello, world!\" is ")
(display (string-length "Hello, world!"))
(display ".")
(newline)`,
		o: `The length of "Hello, world!" is 13.
`,
	},
	{
		i: `(define (print msg) (display msg) (newline))
(print "Hello, lambda!")
(print "Hello, world!")`,
		o: `Hello, lambda!
Hello, world!
`,
	},
	{
		i: `(define (print msg) (display msg) (newline))
(define msg "Hello, msg!")
(set! msg "Hello, set!")
(print msg)
`,
		o: `Hello, set!
`,
	},
	{
		i: `
(define (say-maker header msg trailer)
  (lambda (pre post)
    (display header)
    (display pre)
    (display msg)
    (display post)
    (display trailer)
    (newline)))

(define a (say-maker "<html>" "Hello, a!" "</html>"))
(define b (say-maker "<div>" "Hello, b!" "</div>"))

(a "(" ")")
(b "{" "}")
`,
		o: `<html>(Hello, a!)</html>
<div>{Hello, b!}</div>
`,
	},
	{
		i: `(+ 1 2 3)`,
		v: NewNumber(0, 6),
		o: ``,
	},
	{
		i: `(+)`,
		v: NewNumber(0, 0),
		o: ``,
	},
	{
		i: `(* 1 2 3)`,
		v: NewNumber(0, 6),
		o: ``,
	},
	{
		i: `(*)`,
		v: NewNumber(0, 1),
		o: ``,
	},
	{
		i: `(begin 1 2 3 4)`,
		v: NewNumber(0, 4),
		o: ``,
	},
	{
		i: `(if #t 1 2)`,
		v: NewNumber(0, 1),
		o: ``,
	},

	{
		i: `(if #f 1 2)`,
		v: NewNumber(0, 2),
		o: ``,
	},
	{
		i: `(pair? (cons 1 2))`,
		v: Boolean(true),
	},
	{
		i: `(pair? 3)`,
		v: Boolean(false),
	},
	{
		i: `(cons 1 2)`,
		v: NewPair(NewNumber(0, 1), NewNumber(0, 2)),
	},
	{
		i: `(car (cons 1 2))`,
		v: NewNumber(0, 1),
	},
	{
		i: `(cdr (cons 1 2))`,
		v: NewNumber(0, 2),
	},
	{
		i: `(define v (cons 1 2)) (set-car! v 42) v`,
		v: NewPair(NewNumber(0, 42), NewNumber(0, 2)),
	},
	{
		i: `(define v (cons 1 2)) (set-cdr! v 42) v`,
		v: NewPair(NewNumber(0, 1), NewNumber(0, 42)),
	},
	{
		i: `(null? (list))`,
		v: Boolean(true),
	},
	{
		i: `(null? (list 1))`,
		v: Boolean(false),
	},
	{
		i: `(list? (list))`,
		v: Boolean(true),
	},
	{
		i: `(list? (list 1))`,
		v: Boolean(true),
	},
	{
		i: `(list? 1)`,
		v: Boolean(false),
	},
	{
		i: `(length (list))`,
		v: NewNumber(0, 0),
	},
	{
		i: `(length (list 1 2 3))`,
		v: NewNumber(0, 3),
	},
	{
		i: `(list-tail (list 1 2 3) 2)`,
		v: NewPair(NewNumber(0, 3), nil),
	},
	{
		i: `(list-tail (list 1 2 3) 3)`,
		v: nil,
	},
	{
		i: `(list-ref (list 1 2 3) 0)`,
		v: NewNumber(0, 1),
	},
	{
		i: `(list-ref (list 1 2 3) 1)`,
		v: NewNumber(0, 2),
	},
	{
		i: `(list-ref (list 1 2 3) 2)`,
		v: NewNumber(0, 3),
	},
}

func TestVM(t *testing.T) {
	for idx, test := range vmTests {
		scm, err := New()
		if err != nil {
			t.Fatalf("failed to create virtual machine: %v", err)
		}
		stdout := &strings.Builder{}
		scm.Stdout = stdout

		v, err := scm.Eval(fmt.Sprintf("test-%d", idx),
			strings.NewReader(test.i))
		if err != nil {
			t.Fatalf("Test %d: Eval failed: %v", idx, err)
		}
		if !Equal(v, test.v) {
			t.Errorf("Eval failed: got %v, expected %v", v, test.v)
		}
		output := stdout.String()
		if output != test.o {
			t.Errorf("unexpected output: got '%v', expected '%v'",
				output, test.o)
		}
	}
}
