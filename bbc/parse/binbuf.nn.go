package parse

import (
	"fmt"
	"strconv"

	"github.com/gemrs/gem/bbc/ast"
)
import (
	"bufio"
	"io"
	"strings"
)

type frame struct {
	i            int
	s            string
	line, column int
}
type Lexer struct {
	// The lexer runs in its own goroutine, and communicates via channel 'ch'.
	ch chan frame
	// We record the level of nesting because the action could return, and a
	// subsequent call expects to pick up where it left off. In other words,
	// we're simulating a coroutine.
	// TODO: Support a channel-based variant that compatible with Go's yacc.
	stack []frame
	stale bool

	// The 'l' and 'c' fields were added for
	// https://github.com/wagerlabs/docker/blob/65694e801a7b80930961d70c69cba9f2465459be/buildfile.nex
	// Since then, I introduced the built-in Line() and Column() functions.
	l, c int

	parseResult interface{}

	// The following line makes it easy for scripts to insert fields in the
	// generated code.
	// [NEX_END_OF_LEXER_STRUCT]
}

// NewLexerWithInit creates a new Lexer object, runs the given callback on it,
// then returns it.
func NewLexerWithInit(in io.Reader, initFun func(*Lexer)) *Lexer {
	type dfa struct {
		acc          []bool           // Accepting states.
		f            []func(rune) int // Transitions.
		startf, endf []int            // Transitions at start and end of input.
		nest         []dfa
	}
	yylex := new(Lexer)
	if initFun != nil {
		initFun(yylex)
	}
	yylex.ch = make(chan frame)
	var scan func(in *bufio.Reader, ch chan frame, family []dfa, line, column int)
	scan = func(in *bufio.Reader, ch chan frame, family []dfa, line, column int) {
		// Index of DFA and length of highest-precedence match so far.
		matchi, matchn := 0, -1
		var buf []rune
		n := 0
		checkAccept := func(i int, st int) bool {
			// Higher precedence match? DFAs are run in parallel, so matchn is at most len(buf), hence we may omit the length equality check.
			if family[i].acc[st] && (matchn < n || matchi > i) {
				matchi, matchn = i, n
				return true
			}
			return false
		}
		var state [][2]int
		for i := 0; i < len(family); i++ {
			mark := make([]bool, len(family[i].startf))
			// Every DFA starts at state 0.
			st := 0
			for {
				state = append(state, [2]int{i, st})
				mark[st] = true
				// As we're at the start of input, follow all ^ transitions and append to our list of start states.
				st = family[i].startf[st]
				if -1 == st || mark[st] {
					break
				}
				// We only check for a match after at least one transition.
				checkAccept(i, st)
			}
		}
		atEOF := false
		for {
			if n == len(buf) && !atEOF {
				r, _, err := in.ReadRune()
				switch err {
				case io.EOF:
					atEOF = true
				case nil:
					buf = append(buf, r)
				default:
					panic(err)
				}
			}
			if !atEOF {
				r := buf[n]
				n++
				var nextState [][2]int
				for _, x := range state {
					x[1] = family[x[0]].f[x[1]](r)
					if -1 == x[1] {
						continue
					}
					nextState = append(nextState, x)
					checkAccept(x[0], x[1])
				}
				state = nextState
			} else {
			dollar: // Handle $.
				for _, x := range state {
					mark := make([]bool, len(family[x[0]].endf))
					for {
						mark[x[1]] = true
						x[1] = family[x[0]].endf[x[1]]
						if -1 == x[1] || mark[x[1]] {
							break
						}
						if checkAccept(x[0], x[1]) {
							// Unlike before, we can break off the search. Now that we're at the end, there's no need to maintain the state of each DFA.
							break dollar
						}
					}
				}
				state = nil
			}

			if state == nil {
				lcUpdate := func(r rune) {
					if r == '\n' {
						line++
						column = 0
					} else {
						column++
					}
				}
				// All DFAs stuck. Return last match if it exists, otherwise advance by one rune and restart all DFAs.
				if matchn == -1 {
					if len(buf) == 0 { // This can only happen at the end of input.
						break
					}
					lcUpdate(buf[0])
					buf = buf[1:]
				} else {
					text := string(buf[:matchn])
					buf = buf[matchn:]
					matchn = -1
					ch <- frame{matchi, text, line, column}
					if len(family[matchi].nest) > 0 {
						scan(bufio.NewReader(strings.NewReader(text)), ch, family[matchi].nest, line, column)
					}
					if atEOF {
						break
					}
					for _, r := range text {
						lcUpdate(r)
					}
				}
				n = 0
				for i := 0; i < len(family); i++ {
					state = append(state, [2]int{i, 0})
				}
			}
		}
		ch <- frame{-1, "", line, column}
	}
	go scan(bufio.NewReader(in), yylex.ch, []dfa{
		// [0-9]+
		{[]bool{false, true}, []func(rune) int{ // Transitions
			func(r rune) int {
				switch {
				case 48 <= r && r <= 57:
					return 1
				}
				return -1
			},
			func(r rune) int {
				switch {
				case 48 <= r && r <= 57:
					return 1
				}
				return -1
			},
		}, []int{ /* Start-of-input transitions */ -1, -1}, []int{ /* End-of-input transitions */ -1, -1}, nil},

		// type
		{[]bool{false, false, false, false, true}, []func(rune) int{ // Transitions
			func(r rune) int {
				switch r {
				case 101:
					return -1
				case 112:
					return -1
				case 116:
					return 1
				case 121:
					return -1
				}
				return -1
			},
			func(r rune) int {
				switch r {
				case 101:
					return -1
				case 112:
					return -1
				case 116:
					return -1
				case 121:
					return 2
				}
				return -1
			},
			func(r rune) int {
				switch r {
				case 101:
					return -1
				case 112:
					return 3
				case 116:
					return -1
				case 121:
					return -1
				}
				return -1
			},
			func(r rune) int {
				switch r {
				case 101:
					return 4
				case 112:
					return -1
				case 116:
					return -1
				case 121:
					return -1
				}
				return -1
			},
			func(r rune) int {
				switch r {
				case 101:
					return -1
				case 112:
					return -1
				case 116:
					return -1
				case 121:
					return -1
				}
				return -1
			},
		}, []int{ /* Start-of-input transitions */ -1, -1, -1, -1, -1}, []int{ /* End-of-input transitions */ -1, -1, -1, -1, -1}, nil},

		// string
		{[]bool{false, false, false, false, false, false, true}, []func(rune) int{ // Transitions
			func(r rune) int {
				switch r {
				case 103:
					return -1
				case 105:
					return -1
				case 110:
					return -1
				case 114:
					return -1
				case 115:
					return 1
				case 116:
					return -1
				}
				return -1
			},
			func(r rune) int {
				switch r {
				case 103:
					return -1
				case 105:
					return -1
				case 110:
					return -1
				case 114:
					return -1
				case 115:
					return -1
				case 116:
					return 2
				}
				return -1
			},
			func(r rune) int {
				switch r {
				case 103:
					return -1
				case 105:
					return -1
				case 110:
					return -1
				case 114:
					return 3
				case 115:
					return -1
				case 116:
					return -1
				}
				return -1
			},
			func(r rune) int {
				switch r {
				case 103:
					return -1
				case 105:
					return 4
				case 110:
					return -1
				case 114:
					return -1
				case 115:
					return -1
				case 116:
					return -1
				}
				return -1
			},
			func(r rune) int {
				switch r {
				case 103:
					return -1
				case 105:
					return -1
				case 110:
					return 5
				case 114:
					return -1
				case 115:
					return -1
				case 116:
					return -1
				}
				return -1
			},
			func(r rune) int {
				switch r {
				case 103:
					return 6
				case 105:
					return -1
				case 110:
					return -1
				case 114:
					return -1
				case 115:
					return -1
				case 116:
					return -1
				}
				return -1
			},
			func(r rune) int {
				switch r {
				case 103:
					return -1
				case 105:
					return -1
				case 110:
					return -1
				case 114:
					return -1
				case 115:
					return -1
				case 116:
					return -1
				}
				return -1
			},
		}, []int{ /* Start-of-input transitions */ -1, -1, -1, -1, -1, -1, -1}, []int{ /* End-of-input transitions */ -1, -1, -1, -1, -1, -1, -1}, nil},

		// byte
		{[]bool{false, false, false, false, true}, []func(rune) int{ // Transitions
			func(r rune) int {
				switch r {
				case 98:
					return 1
				case 101:
					return -1
				case 116:
					return -1
				case 121:
					return -1
				}
				return -1
			},
			func(r rune) int {
				switch r {
				case 98:
					return -1
				case 101:
					return -1
				case 116:
					return -1
				case 121:
					return 2
				}
				return -1
			},
			func(r rune) int {
				switch r {
				case 98:
					return -1
				case 101:
					return -1
				case 116:
					return 3
				case 121:
					return -1
				}
				return -1
			},
			func(r rune) int {
				switch r {
				case 98:
					return -1
				case 101:
					return 4
				case 116:
					return -1
				case 121:
					return -1
				}
				return -1
			},
			func(r rune) int {
				switch r {
				case 98:
					return -1
				case 101:
					return -1
				case 116:
					return -1
				case 121:
					return -1
				}
				return -1
			},
		}, []int{ /* Start-of-input transitions */ -1, -1, -1, -1, -1}, []int{ /* End-of-input transitions */ -1, -1, -1, -1, -1}, nil},

		// struct
		{[]bool{false, false, false, false, false, false, true}, []func(rune) int{ // Transitions
			func(r rune) int {
				switch r {
				case 99:
					return -1
				case 114:
					return -1
				case 115:
					return 1
				case 116:
					return -1
				case 117:
					return -1
				}
				return -1
			},
			func(r rune) int {
				switch r {
				case 99:
					return -1
				case 114:
					return -1
				case 115:
					return -1
				case 116:
					return 2
				case 117:
					return -1
				}
				return -1
			},
			func(r rune) int {
				switch r {
				case 99:
					return -1
				case 114:
					return 3
				case 115:
					return -1
				case 116:
					return -1
				case 117:
					return -1
				}
				return -1
			},
			func(r rune) int {
				switch r {
				case 99:
					return -1
				case 114:
					return -1
				case 115:
					return -1
				case 116:
					return -1
				case 117:
					return 4
				}
				return -1
			},
			func(r rune) int {
				switch r {
				case 99:
					return 5
				case 114:
					return -1
				case 115:
					return -1
				case 116:
					return -1
				case 117:
					return -1
				}
				return -1
			},
			func(r rune) int {
				switch r {
				case 99:
					return -1
				case 114:
					return -1
				case 115:
					return -1
				case 116:
					return 6
				case 117:
					return -1
				}
				return -1
			},
			func(r rune) int {
				switch r {
				case 99:
					return -1
				case 114:
					return -1
				case 115:
					return -1
				case 116:
					return -1
				case 117:
					return -1
				}
				return -1
			},
		}, []int{ /* Start-of-input transitions */ -1, -1, -1, -1, -1, -1, -1}, []int{ /* End-of-input transitions */ -1, -1, -1, -1, -1, -1, -1}, nil},

		// bitstruct
		{[]bool{false, false, false, false, false, false, false, false, false, true}, []func(rune) int{ // Transitions
			func(r rune) int {
				switch r {
				case 98:
					return 1
				case 99:
					return -1
				case 105:
					return -1
				case 114:
					return -1
				case 115:
					return -1
				case 116:
					return -1
				case 117:
					return -1
				}
				return -1
			},
			func(r rune) int {
				switch r {
				case 98:
					return -1
				case 99:
					return -1
				case 105:
					return 2
				case 114:
					return -1
				case 115:
					return -1
				case 116:
					return -1
				case 117:
					return -1
				}
				return -1
			},
			func(r rune) int {
				switch r {
				case 98:
					return -1
				case 99:
					return -1
				case 105:
					return -1
				case 114:
					return -1
				case 115:
					return -1
				case 116:
					return 3
				case 117:
					return -1
				}
				return -1
			},
			func(r rune) int {
				switch r {
				case 98:
					return -1
				case 99:
					return -1
				case 105:
					return -1
				case 114:
					return -1
				case 115:
					return 4
				case 116:
					return -1
				case 117:
					return -1
				}
				return -1
			},
			func(r rune) int {
				switch r {
				case 98:
					return -1
				case 99:
					return -1
				case 105:
					return -1
				case 114:
					return -1
				case 115:
					return -1
				case 116:
					return 5
				case 117:
					return -1
				}
				return -1
			},
			func(r rune) int {
				switch r {
				case 98:
					return -1
				case 99:
					return -1
				case 105:
					return -1
				case 114:
					return 6
				case 115:
					return -1
				case 116:
					return -1
				case 117:
					return -1
				}
				return -1
			},
			func(r rune) int {
				switch r {
				case 98:
					return -1
				case 99:
					return -1
				case 105:
					return -1
				case 114:
					return -1
				case 115:
					return -1
				case 116:
					return -1
				case 117:
					return 7
				}
				return -1
			},
			func(r rune) int {
				switch r {
				case 98:
					return -1
				case 99:
					return 8
				case 105:
					return -1
				case 114:
					return -1
				case 115:
					return -1
				case 116:
					return -1
				case 117:
					return -1
				}
				return -1
			},
			func(r rune) int {
				switch r {
				case 98:
					return -1
				case 99:
					return -1
				case 105:
					return -1
				case 114:
					return -1
				case 115:
					return -1
				case 116:
					return 9
				case 117:
					return -1
				}
				return -1
			},
			func(r rune) int {
				switch r {
				case 98:
					return -1
				case 99:
					return -1
				case 105:
					return -1
				case 114:
					return -1
				case 115:
					return -1
				case 116:
					return -1
				case 117:
					return -1
				}
				return -1
			},
		}, []int{ /* Start-of-input transitions */ -1, -1, -1, -1, -1, -1, -1, -1, -1, -1}, []int{ /* End-of-input transitions */ -1, -1, -1, -1, -1, -1, -1, -1, -1, -1}, nil},

		// frame
		{[]bool{false, false, false, false, false, true}, []func(rune) int{ // Transitions
			func(r rune) int {
				switch r {
				case 97:
					return -1
				case 101:
					return -1
				case 102:
					return 1
				case 109:
					return -1
				case 114:
					return -1
				}
				return -1
			},
			func(r rune) int {
				switch r {
				case 97:
					return -1
				case 101:
					return -1
				case 102:
					return -1
				case 109:
					return -1
				case 114:
					return 2
				}
				return -1
			},
			func(r rune) int {
				switch r {
				case 97:
					return 3
				case 101:
					return -1
				case 102:
					return -1
				case 109:
					return -1
				case 114:
					return -1
				}
				return -1
			},
			func(r rune) int {
				switch r {
				case 97:
					return -1
				case 101:
					return -1
				case 102:
					return -1
				case 109:
					return 4
				case 114:
					return -1
				}
				return -1
			},
			func(r rune) int {
				switch r {
				case 97:
					return -1
				case 101:
					return 5
				case 102:
					return -1
				case 109:
					return -1
				case 114:
					return -1
				}
				return -1
			},
			func(r rune) int {
				switch r {
				case 97:
					return -1
				case 101:
					return -1
				case 102:
					return -1
				case 109:
					return -1
				case 114:
					return -1
				}
				return -1
			},
		}, []int{ /* Start-of-input transitions */ -1, -1, -1, -1, -1, -1}, []int{ /* End-of-input transitions */ -1, -1, -1, -1, -1, -1}, nil},

		// bit
		{[]bool{false, false, false, true}, []func(rune) int{ // Transitions
			func(r rune) int {
				switch r {
				case 98:
					return 1
				case 105:
					return -1
				case 116:
					return -1
				}
				return -1
			},
			func(r rune) int {
				switch r {
				case 98:
					return -1
				case 105:
					return 2
				case 116:
					return -1
				}
				return -1
			},
			func(r rune) int {
				switch r {
				case 98:
					return -1
				case 105:
					return -1
				case 116:
					return 3
				}
				return -1
			},
			func(r rune) int {
				switch r {
				case 98:
					return -1
				case 105:
					return -1
				case 116:
					return -1
				}
				return -1
			},
		}, []int{ /* Start-of-input transitions */ -1, -1, -1, -1}, []int{ /* End-of-input transitions */ -1, -1, -1, -1}, nil},

		// Fixed
		{[]bool{false, false, false, false, false, true}, []func(rune) int{ // Transitions
			func(r rune) int {
				switch r {
				case 70:
					return 1
				case 100:
					return -1
				case 101:
					return -1
				case 105:
					return -1
				case 120:
					return -1
				}
				return -1
			},
			func(r rune) int {
				switch r {
				case 70:
					return -1
				case 100:
					return -1
				case 101:
					return -1
				case 105:
					return 2
				case 120:
					return -1
				}
				return -1
			},
			func(r rune) int {
				switch r {
				case 70:
					return -1
				case 100:
					return -1
				case 101:
					return -1
				case 105:
					return -1
				case 120:
					return 3
				}
				return -1
			},
			func(r rune) int {
				switch r {
				case 70:
					return -1
				case 100:
					return -1
				case 101:
					return 4
				case 105:
					return -1
				case 120:
					return -1
				}
				return -1
			},
			func(r rune) int {
				switch r {
				case 70:
					return -1
				case 100:
					return 5
				case 101:
					return -1
				case 105:
					return -1
				case 120:
					return -1
				}
				return -1
			},
			func(r rune) int {
				switch r {
				case 70:
					return -1
				case 100:
					return -1
				case 101:
					return -1
				case 105:
					return -1
				case 120:
					return -1
				}
				return -1
			},
		}, []int{ /* Start-of-input transitions */ -1, -1, -1, -1, -1, -1}, []int{ /* End-of-input transitions */ -1, -1, -1, -1, -1, -1}, nil},

		// Var8
		{[]bool{false, false, false, false, true}, []func(rune) int{ // Transitions
			func(r rune) int {
				switch r {
				case 56:
					return -1
				case 86:
					return 1
				case 97:
					return -1
				case 114:
					return -1
				}
				return -1
			},
			func(r rune) int {
				switch r {
				case 56:
					return -1
				case 86:
					return -1
				case 97:
					return 2
				case 114:
					return -1
				}
				return -1
			},
			func(r rune) int {
				switch r {
				case 56:
					return -1
				case 86:
					return -1
				case 97:
					return -1
				case 114:
					return 3
				}
				return -1
			},
			func(r rune) int {
				switch r {
				case 56:
					return 4
				case 86:
					return -1
				case 97:
					return -1
				case 114:
					return -1
				}
				return -1
			},
			func(r rune) int {
				switch r {
				case 56:
					return -1
				case 86:
					return -1
				case 97:
					return -1
				case 114:
					return -1
				}
				return -1
			},
		}, []int{ /* Start-of-input transitions */ -1, -1, -1, -1, -1}, []int{ /* End-of-input transitions */ -1, -1, -1, -1, -1}, nil},

		// Var16
		{[]bool{false, false, false, false, false, true}, []func(rune) int{ // Transitions
			func(r rune) int {
				switch r {
				case 49:
					return -1
				case 54:
					return -1
				case 86:
					return 1
				case 97:
					return -1
				case 114:
					return -1
				}
				return -1
			},
			func(r rune) int {
				switch r {
				case 49:
					return -1
				case 54:
					return -1
				case 86:
					return -1
				case 97:
					return 2
				case 114:
					return -1
				}
				return -1
			},
			func(r rune) int {
				switch r {
				case 49:
					return -1
				case 54:
					return -1
				case 86:
					return -1
				case 97:
					return -1
				case 114:
					return 3
				}
				return -1
			},
			func(r rune) int {
				switch r {
				case 49:
					return 4
				case 54:
					return -1
				case 86:
					return -1
				case 97:
					return -1
				case 114:
					return -1
				}
				return -1
			},
			func(r rune) int {
				switch r {
				case 49:
					return -1
				case 54:
					return 5
				case 86:
					return -1
				case 97:
					return -1
				case 114:
					return -1
				}
				return -1
			},
			func(r rune) int {
				switch r {
				case 49:
					return -1
				case 54:
					return -1
				case 86:
					return -1
				case 97:
					return -1
				case 114:
					return -1
				}
				return -1
			},
		}, []int{ /* Start-of-input transitions */ -1, -1, -1, -1, -1, -1}, []int{ /* End-of-input transitions */ -1, -1, -1, -1, -1, -1}, nil},

		// u?int(8|16|24|32|64)
		{[]bool{false, false, false, false, false, false, false, true, false, false, true, true, true, true}, []func(rune) int{ // Transitions
			func(r rune) int {
				switch r {
				case 49:
					return -1
				case 50:
					return -1
				case 51:
					return -1
				case 52:
					return -1
				case 54:
					return -1
				case 56:
					return -1
				case 105:
					return 2
				case 110:
					return -1
				case 116:
					return -1
				case 117:
					return 1
				}
				return -1
			},
			func(r rune) int {
				switch r {
				case 49:
					return -1
				case 50:
					return -1
				case 51:
					return -1
				case 52:
					return -1
				case 54:
					return -1
				case 56:
					return -1
				case 105:
					return 2
				case 110:
					return -1
				case 116:
					return -1
				case 117:
					return -1
				}
				return -1
			},
			func(r rune) int {
				switch r {
				case 49:
					return -1
				case 50:
					return -1
				case 51:
					return -1
				case 52:
					return -1
				case 54:
					return -1
				case 56:
					return -1
				case 105:
					return -1
				case 110:
					return 3
				case 116:
					return -1
				case 117:
					return -1
				}
				return -1
			},
			func(r rune) int {
				switch r {
				case 49:
					return -1
				case 50:
					return -1
				case 51:
					return -1
				case 52:
					return -1
				case 54:
					return -1
				case 56:
					return -1
				case 105:
					return -1
				case 110:
					return -1
				case 116:
					return 4
				case 117:
					return -1
				}
				return -1
			},
			func(r rune) int {
				switch r {
				case 49:
					return 8
				case 50:
					return 5
				case 51:
					return 6
				case 52:
					return -1
				case 54:
					return 9
				case 56:
					return 7
				case 105:
					return -1
				case 110:
					return -1
				case 116:
					return -1
				case 117:
					return -1
				}
				return -1
			},
			func(r rune) int {
				switch r {
				case 49:
					return -1
				case 50:
					return -1
				case 51:
					return -1
				case 52:
					return 13
				case 54:
					return -1
				case 56:
					return -1
				case 105:
					return -1
				case 110:
					return -1
				case 116:
					return -1
				case 117:
					return -1
				}
				return -1
			},
			func(r rune) int {
				switch r {
				case 49:
					return -1
				case 50:
					return 12
				case 51:
					return -1
				case 52:
					return -1
				case 54:
					return -1
				case 56:
					return -1
				case 105:
					return -1
				case 110:
					return -1
				case 116:
					return -1
				case 117:
					return -1
				}
				return -1
			},
			func(r rune) int {
				switch r {
				case 49:
					return -1
				case 50:
					return -1
				case 51:
					return -1
				case 52:
					return -1
				case 54:
					return -1
				case 56:
					return -1
				case 105:
					return -1
				case 110:
					return -1
				case 116:
					return -1
				case 117:
					return -1
				}
				return -1
			},
			func(r rune) int {
				switch r {
				case 49:
					return -1
				case 50:
					return -1
				case 51:
					return -1
				case 52:
					return -1
				case 54:
					return 11
				case 56:
					return -1
				case 105:
					return -1
				case 110:
					return -1
				case 116:
					return -1
				case 117:
					return -1
				}
				return -1
			},
			func(r rune) int {
				switch r {
				case 49:
					return -1
				case 50:
					return -1
				case 51:
					return -1
				case 52:
					return 10
				case 54:
					return -1
				case 56:
					return -1
				case 105:
					return -1
				case 110:
					return -1
				case 116:
					return -1
				case 117:
					return -1
				}
				return -1
			},
			func(r rune) int {
				switch r {
				case 49:
					return -1
				case 50:
					return -1
				case 51:
					return -1
				case 52:
					return -1
				case 54:
					return -1
				case 56:
					return -1
				case 105:
					return -1
				case 110:
					return -1
				case 116:
					return -1
				case 117:
					return -1
				}
				return -1
			},
			func(r rune) int {
				switch r {
				case 49:
					return -1
				case 50:
					return -1
				case 51:
					return -1
				case 52:
					return -1
				case 54:
					return -1
				case 56:
					return -1
				case 105:
					return -1
				case 110:
					return -1
				case 116:
					return -1
				case 117:
					return -1
				}
				return -1
			},
			func(r rune) int {
				switch r {
				case 49:
					return -1
				case 50:
					return -1
				case 51:
					return -1
				case 52:
					return -1
				case 54:
					return -1
				case 56:
					return -1
				case 105:
					return -1
				case 110:
					return -1
				case 116:
					return -1
				case 117:
					return -1
				}
				return -1
			},
			func(r rune) int {
				switch r {
				case 49:
					return -1
				case 50:
					return -1
				case 51:
					return -1
				case 52:
					return -1
				case 54:
					return -1
				case 56:
					return -1
				case 105:
					return -1
				case 110:
					return -1
				case 116:
					return -1
				case 117:
					return -1
				}
				return -1
			},
		}, []int{ /* Start-of-input transitions */ -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1}, []int{ /* End-of-input transitions */ -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1}, nil},

		// Int[a-zA-Z0-9]+
		{[]bool{false, false, false, false, true}, []func(rune) int{ // Transitions
			func(r rune) int {
				switch r {
				case 73:
					return 1
				case 110:
					return -1
				case 116:
					return -1
				}
				switch {
				case 97 <= r && r <= 122:
					return -1
				case 48 <= r && r <= 57:
					return -1
				case 65 <= r && r <= 90:
					return -1
				}
				return -1
			},
			func(r rune) int {
				switch r {
				case 73:
					return -1
				case 110:
					return 2
				case 116:
					return -1
				}
				switch {
				case 97 <= r && r <= 122:
					return -1
				case 48 <= r && r <= 57:
					return -1
				case 65 <= r && r <= 90:
					return -1
				}
				return -1
			},
			func(r rune) int {
				switch r {
				case 73:
					return -1
				case 110:
					return -1
				case 116:
					return 3
				}
				switch {
				case 97 <= r && r <= 122:
					return -1
				case 48 <= r && r <= 57:
					return -1
				case 65 <= r && r <= 90:
					return -1
				}
				return -1
			},
			func(r rune) int {
				switch r {
				case 73:
					return 4
				case 110:
					return 4
				case 116:
					return 4
				}
				switch {
				case 97 <= r && r <= 122:
					return 4
				case 48 <= r && r <= 57:
					return 4
				case 65 <= r && r <= 90:
					return 4
				}
				return -1
			},
			func(r rune) int {
				switch r {
				case 73:
					return 4
				case 110:
					return 4
				case 116:
					return 4
				}
				switch {
				case 97 <= r && r <= 122:
					return 4
				case 48 <= r && r <= 57:
					return 4
				case 65 <= r && r <= 90:
					return 4
				}
				return -1
			},
		}, []int{ /* Start-of-input transitions */ -1, -1, -1, -1, -1}, []int{ /* End-of-input transitions */ -1, -1, -1, -1, -1}, nil},

		// [{}\[\]<>,]
		{[]bool{false, true}, []func(rune) int{ // Transitions
			func(r rune) int {
				switch r {
				case 44:
					return 1
				case 60:
					return 1
				case 62:
					return 1
				case 91:
					return 1
				case 93:
					return 1
				case 123:
					return 1
				case 125:
					return 1
				}
				return -1
			},
			func(r rune) int {
				switch r {
				case 44:
					return -1
				case 60:
					return -1
				case 62:
					return -1
				case 91:
					return -1
				case 93:
					return -1
				case 123:
					return -1
				case 125:
					return -1
				}
				return -1
			},
		}, []int{ /* Start-of-input transitions */ -1, -1}, []int{ /* End-of-input transitions */ -1, -1}, nil},

		// [a-zA-Z_]+([0-9a-zA-Z_]+)?
		{[]bool{false, true, true, true}, []func(rune) int{ // Transitions
			func(r rune) int {
				switch r {
				case 95:
					return 1
				}
				switch {
				case 48 <= r && r <= 57:
					return -1
				case 65 <= r && r <= 90:
					return 1
				case 97 <= r && r <= 122:
					return 1
				}
				return -1
			},
			func(r rune) int {
				switch r {
				case 95:
					return 2
				}
				switch {
				case 48 <= r && r <= 57:
					return 3
				case 65 <= r && r <= 90:
					return 2
				case 97 <= r && r <= 122:
					return 2
				}
				return -1
			},
			func(r rune) int {
				switch r {
				case 95:
					return 2
				}
				switch {
				case 48 <= r && r <= 57:
					return 3
				case 65 <= r && r <= 90:
					return 2
				case 97 <= r && r <= 122:
					return 2
				}
				return -1
			},
			func(r rune) int {
				switch r {
				case 95:
					return 3
				}
				switch {
				case 48 <= r && r <= 57:
					return 3
				case 65 <= r && r <= 90:
					return 3
				case 97 <= r && r <= 122:
					return 3
				}
				return -1
			},
		}, []int{ /* Start-of-input transitions */ -1, -1, -1, -1}, []int{ /* End-of-input transitions */ -1, -1, -1, -1}, nil},

		// [ \t\n]+
		{[]bool{false, true}, []func(rune) int{ // Transitions
			func(r rune) int {
				switch r {
				case 9:
					return 1
				case 10:
					return 1
				case 32:
					return 1
				}
				return -1
			},
			func(r rune) int {
				switch r {
				case 9:
					return 1
				case 10:
					return 1
				case 32:
					return 1
				}
				return -1
			},
		}, []int{ /* Start-of-input transitions */ -1, -1}, []int{ /* End-of-input transitions */ -1, -1}, nil},

		// \/\/[^\n]*
		{[]bool{false, false, true, true}, []func(rune) int{ // Transitions
			func(r rune) int {
				switch r {
				case 10:
					return -1
				case 47:
					return 1
				}
				return -1
			},
			func(r rune) int {
				switch r {
				case 10:
					return -1
				case 47:
					return 2
				}
				return -1
			},
			func(r rune) int {
				switch r {
				case 10:
					return -1
				case 47:
					return 3
				}
				return 3
			},
			func(r rune) int {
				switch r {
				case 10:
					return -1
				case 47:
					return 3
				}
				return 3
			},
		}, []int{ /* Start-of-input transitions */ -1, -1, -1, -1}, []int{ /* End-of-input transitions */ -1, -1, -1, -1}, nil},

		// \/\*([^*]|[\r\n]|(\*+([^*\/]|[\r\n])))*\*\/
		{[]bool{false, false, false, false, false, false, false, true, false, false}, []func(rune) int{ // Transitions
			func(r rune) int {
				switch r {
				case 10:
					return -1
				case 13:
					return -1
				case 42:
					return -1
				case 47:
					return 1
				}
				return -1
			},
			func(r rune) int {
				switch r {
				case 10:
					return -1
				case 13:
					return -1
				case 42:
					return 2
				case 47:
					return -1
				}
				return -1
			},
			func(r rune) int {
				switch r {
				case 10:
					return 3
				case 13:
					return 3
				case 42:
					return 5
				case 47:
					return 4
				}
				return 4
			},
			func(r rune) int {
				switch r {
				case 10:
					return 3
				case 13:
					return 3
				case 42:
					return 5
				case 47:
					return 4
				}
				return 4
			},
			func(r rune) int {
				switch r {
				case 10:
					return 3
				case 13:
					return 3
				case 42:
					return 5
				case 47:
					return 4
				}
				return 4
			},
			func(r rune) int {
				switch r {
				case 10:
					return 6
				case 13:
					return 6
				case 42:
					return 8
				case 47:
					return 7
				}
				return 9
			},
			func(r rune) int {
				switch r {
				case 10:
					return 3
				case 13:
					return 3
				case 42:
					return 5
				case 47:
					return 4
				}
				return 4
			},
			func(r rune) int {
				switch r {
				case 10:
					return -1
				case 13:
					return -1
				case 42:
					return -1
				case 47:
					return -1
				}
				return -1
			},
			func(r rune) int {
				switch r {
				case 10:
					return 6
				case 13:
					return 6
				case 42:
					return 8
				case 47:
					return -1
				}
				return 9
			},
			func(r rune) int {
				switch r {
				case 10:
					return 3
				case 13:
					return 3
				case 42:
					return 5
				case 47:
					return 4
				}
				return 4
			},
		}, []int{ /* Start-of-input transitions */ -1, -1, -1, -1, -1, -1, -1, -1, -1, -1}, []int{ /* End-of-input transitions */ -1, -1, -1, -1, -1, -1, -1, -1, -1, -1}, nil},

		// .
		{[]bool{false, true}, []func(rune) int{ // Transitions
			func(r rune) int {
				return 1
			},
			func(r rune) int {
				return -1
			},
		}, []int{ /* Start-of-input transitions */ -1, -1}, []int{ /* End-of-input transitions */ -1, -1}, nil},
	}, 0, 0)
	return yylex
}

func NewLexer(in io.Reader) *Lexer {
	return NewLexerWithInit(in, nil)
}

// Text returns the matched text.
func (yylex *Lexer) Text() string {
	return yylex.stack[len(yylex.stack)-1].s
}

// Line returns the current line number.
// The first line is 0.
func (yylex *Lexer) Line() int {
	if len(yylex.stack) == 0 {
		return 0
	}
	return yylex.stack[len(yylex.stack)-1].line
}

// Column returns the current column number.
// The first column is 0.
func (yylex *Lexer) Column() int {
	if len(yylex.stack) == 0 {
		return 0
	}
	return yylex.stack[len(yylex.stack)-1].column
}

func (yylex *Lexer) next(lvl int) int {
	if lvl == len(yylex.stack) {
		l, c := 0, 0
		if lvl > 0 {
			l, c = yylex.stack[lvl-1].line, yylex.stack[lvl-1].column
		}
		yylex.stack = append(yylex.stack, frame{0, "", l, c})
	}
	if lvl == len(yylex.stack)-1 {
		p := &yylex.stack[lvl]
		*p = <-yylex.ch
		yylex.stale = false
	} else {
		yylex.stale = true
	}
	return yylex.stack[lvl].i
}
func (yylex *Lexer) pop() {
	yylex.stack = yylex.stack[:len(yylex.stack)-1]
}

// Lex runs the lexer. Always returns 0.
// When the -s option is given, this function is not generated;
// instead, the NN_FUN macro runs the lexer.
func (yylex *Lexer) Lex(lval *yySymType) int {
OUTER0:
	for {
		switch yylex.next(0) {
		case 0:
			{
				var err error
				lval.ival, err = strconv.Atoi(yylex.Text())
				if err != nil {
					yylex.Error(err.Error())
				}
				return tNumber
			}
		case 1:
			{
				return tType
			}
		case 2:
			{
				return tStringType
			}
		case 3:
			{
				return tByteType
			}
		case 4:
			{
				return tStruct
			}
		case 5:
			{
				return tBitStruct
			}
		case 6:
			{
				return tFrame
			}
		case 7:
			{
				return tBitsType
			}
		case 8:
			{
				return tFrameFixed
			}
		case 9:
			{
				return tFrameVar8
			}
		case 10:
			{
				return tFrameVar16
			}
		case 11:
			{
				var err error
				lval.n, err = ast.ParseIntegerType(yylex.Text())
				if err != nil {
					yylex.Error(err.Error())
				}
				return tIntegerType
			}
		case 12:
			{
				lval.sval = yylex.Text()
				return tIntegerFlag
			}
		case 13:
			{
				return int(yylex.Text()[0])
			}
		case 14:
			{
				lval.sval = yylex.Text()
				return tIdentifier
			}
		case 15:
			{ /* eat up whitespace */
			}
		case 16:
			{ /* eat up one-line comments */
			}
		case 17:
			{ /* eat up multi-line comments. ugly but functional regex */
			}
		case 18:
			{
				yylex.Error(fmt.Sprintf("unrecognized character: %v", yylex.Text()))
			}
		default:
			break OUTER0
		}
		continue
	}
	yylex.pop()

	return 0
}
