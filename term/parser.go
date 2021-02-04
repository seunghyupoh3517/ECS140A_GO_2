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

// New unambiguous grammar with left  factorization 
// <term> ::= ATOM <new> | NUM | VAR
// <new> ::= nil | ( <args> )
// <args> ::= <term> <new2>
// <new2> ::= nil | , <args>

// Parsing  table
//      |	$			atom				num 			var					(				)			,
//		-----------------------------------------------------------------------------------------------------------------
// S	|   X			S::=term $			S::=term $		S::=term $			X				X			X  
// term |	X			term::=atom new 	term:=num		term:=var			X				X			X
// new	|	new::=nil	X					X				X					new::=( args )	new::=nil	new::=nil
// args	|	X			args::=term new2	args:=term new2 args:=term  new2	X				X			X
// new2	|	X			X					X				X					X				new2:=nil	new2:=, args

//	testing bar(1, a, foo(X))
//	Matched				Stack				Input (token)		Action 
//						term$				bar(1, a, foo(X))$
//						atom new$			bar(1, a, foo(X))$	output term -> atom new	
//	bar					new$				(1, a, foo(X))$		match bar - atom
//	bar					(args)$				(1, a, foo(X))$		output new -> (args)
//	bar(				args)$				1, a, foo(X))$		match (
//	bar(				term new2)$ 		1, a, foo(X))$		output args -> term new2
// 	bar(				num new2)$			1, a, foo(X))$		output term -> num
//	bar(1				new2)$				, a, foo(X))$		match 1 - num
//	bar(1				, args)$			, a, foo(X))$		output new2 -> , args
//	bar(1,				args)$				a, foo(X))$			match ,
//	bar(1,				term new2)$ 		a, foo(X))$			output args -> term new2
// 	bar(1,				atom new new2)$		a, foo(X))$			output term -> atom new
//	bar(1, a			new new2)$			, foo(X))$			match a - atom
//  bar(1, a			nil new2)$			, foo(X))$			output new -> nil
//  bar(1, a			, args)$			, foo(X))$			output new2 -> , args
//  bar(1, a,   		args)$				foo(X))$			match ,
//	bar(1, a,  	 		term new2)$			foo(X))$			output args -> term new2
//	bar(1, a, 			atom new new2)$		foo(X))$			output term -> atom new
//	bar(1, a, foo		new new2)$			(X))$				match foo - atom
// 	bar(1, a, foo		(args)new2)$		(X))$				output new -> (args)
//  bar(1, a, foo(  	args)new2)$			X))$				match (
//  bar(1, a, foo(  	term new2)new2)$	X))$				output args -> term new2
//  bar(1, a, foo(  	var new2)new2)$		X))$				output term -> var
//  bar(1, a, foo(X 	new2)new2)$			))$					match X - var
//  bar(1, a, foo(X	  	nill)new2)$			))$					output new2 -> nil
//  bar(1, a, foo(X)  	new2)$				)$					match )
//  bar(1, a, foo(X)  	nil)$				)$					output new2 -> nil
//  bar(1, a, foo(X))  	$					$					match )
//	COMPLETED


// I was going to use two different types - non terminal and terminal, and differentiate the numbers with the type however
// I cannot assign two different types to an array, thus, I am using one type containing both nonterminal and terminal and am going to
// differentiate them by if condition nonterminal 1 - 5, terminal 6 - 13 - deduct 5 to use it as index.
// Pasing table ERROR if the first entry is 0  
type parsingEntry int
const ( 		
	S		parsingEntry = iota + 1 // 1
	Term	// 2
	New 	// 3
	Args	// 4
	Neww	// 5
	Dollar	// When entry > 5, deduct 5 in order to use it with the index 
	Atom	// 7
	Num		// 8
	Var		// 9
	Lpar	// 10
	Rpar	// 11
	Comma	// 12
	Nil		// 13
)

type rule [3]parsingEntry
S1, S2, S3 := rule{Term, Dollar}, rule{Term, Dollar}, rule{Term, Dollar}
Term1, Term2, Term3 := rule{Atom, New}, rule{Num}, rule{Var}
New1, New2, New3, New4 := rule{Nil}, rule{Lpar, Args, Rpar}, rule{Nil}, rule{Nil}
Args1, Args2, Arg3 := rule{Term, Neww}, rule{Term, Neww}, rule{Term, Neww}
Neww1, Neww2 := rule{Nil}, rule{Comma, Args}
	
var parseTable [5][7]rule
parseTable[0][1], parseTable[0][2], parseTable[0][3] = S1, S2, S3
parseTable[1][1], parseTable[1][2], parseTable[1][3] = Term1, Term2, Term3
parseTable[2][0], parseTable[2][4], parseTable[2][5], parseTable[2][6] = New1, New2, New3, New4
parseTable[3][1], parseTable[3][2], parseTable[3][3] = Args1, Args2, Arg3
parseTable[4][5], parseTable[4][6] = Neww1, Neww2

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
