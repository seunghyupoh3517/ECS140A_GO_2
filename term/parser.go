package term

import (
	"errors"
)

// ErrParser is the error value returned by the Parser if the string is not a
// valid term.
// See also https://golang.org/pkg/errors/#New
// and // https://golang.org/pkg/builtin/#error
var ErrParser = errors.New("parser error")

//
// <term>     ::= ATOM | NUM | VAR | <compound>
// <compound> ::= <functor> LPAR <args> RPAR
// <functor>  ::= ATOM
// <args>     ::= <term> | <term> COMMA <args>
//

// Parser is the interface for the term parser.
// Do not change the definition of this interface.
type Parser interface {
	Parse(string) (*Term, error)
}

// Grammar has map from token to the term
type Grammar struct {
	grammar map[*Term][]*Term    // term -> term[]
}

// NewParser creates a struct of a type that satisfies the Parser interface.
func NewParser() Parser {
	var parseGrammar Parser := Grammar{make(map[*Token]*Term)}
	return parseGrammar
}


// what dose Parse need to do?
// it receives a string, we need to convert it to the term type, Otherwise, return a error for invalid string given to us
// we can apply the lexer to separate the string to tokens
// For converting to the term, we need to figure out several attributes and assign to the Term struct
	// term type
	// content
	// functor, the name of the function
	// args
	// use term_test as reference, it shows how to initialize a Term object

type nonTermial int

const (
	Start nonTermial = iota
	Term 
	New 
	Compound 
	Funct 
	Argus 
	New'
	// arr[S][atom]
)

func (g Grammar) Parse(str string) (*Term, error) {
	// TODO: matrix in the global
	var parseTable [][]int

	lex := newLexer(str)

	// convert the input string to tokens
	var	tokenList []*Token
	for token, err := lex.next(); err != ErrLexer {
		tokenList.append(token)
	} else {
		// validating the given string, return error if can't parse to token
		// TODO: double check with error object return type
		return nil, fmt.Errorf(lex.ErrLexer)
	}

	// [ "bar", "(", "x", ")" ]
	// pointer point to the current token
	var tokenInd = 0
	// initialize the stack
	var stack []nonTermial
	stack.append(Start)
 	while (len(stack) != 0 ) {
 		ind := len(stack) - 1		// ind 
 		topOfStack := stack[ind]

 		switch typ := topOfStack.(type) {
 		case tokenType:
 			// when the top is terminal
 			// pop out the terminal and advance the tokenList index 
 			// when the tokenType(terminal) are the same
 			if tokenList[tokenInd].typ == topOfStack.typ {
 				stack = stack[:ind]		// pop out the top element
 				tokenInd++;	
 			} else {
 				// terminal is not match
 				return nil, fmt.Errorf("Invalid input string")
 			}

 		case nonTermial: 
 			// when the top is non terminal
 			if parseTable[topOfStack][]

 		default:
 			// invalid type in the stack
 			return nil, fmt.Errorf("Invalid type in the stack")
 		}
 	}


}