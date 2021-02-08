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

var relationMap = map[tokenType]TermType { // tokenType -> TermType
	tokenAtom: TermAtom,
	tokenNumber: TermNumber,
	tokenVariable: TermVariable,
}

// NewParser creates a struct of a type that satisfies the Parser interface.
func NewParser() Parser {
	var parseGrammar Parser = Grammar{make(map[*Term][]*Term)}
	return parseGrammar
}

// equivalent grammar
// <term> ::= ATOM <new> | NUM | VAR
// <new> ::= nil | ( <args> )
// <args> ::= <term> <new2>
// <new2> ::= nil | , <args>

// nonTerminal enumerates all types to non terminal
type nonTermial int

const (
	Start_NT nonTermial = iota
	Term_NT
	NT1
	Args_NT
	NT2
)

// implements the Parse method with Grammar struct
func (g Grammar) Parse(str string) (*Term, error) {
	// TODO: matrix in the global
	// The parseTable can be the type of
	// []interface{} which is the list in a single cell
	var parseTable [][][]interface{}
	var stk1 []*Term{} // functor
	var stk2 [][]*Term{} // argument List
	var stk_ptr = 0
	var termMap = map[string]*Term{} // term.toString() -> *term
	lex := newLexer(str)

	// convert the input string to tokens
	var	tokenList []*Token
	for {
		token, err := lex.next()
		if err == ErrLexer {
			// validating the given string, return error if can't parse to token
			// TODO: double check with error object return type
			return nil, err
			} else {
				if token.typ == tokenEOF {
					break
				}
				tokenList = append(tokenList, token)
			}
	}

	// [ "bar", "(", "x", ")" ]
	// pointer point to the current token in the list
	var tokenInd = 0
	// initialize the stack
	// stack needs to accept two data types, nonTerminal & tokenType
	var stack []interface{}
	stack = append(stack, Start_NT)

 	for len(stack) != 0 {
 		ind := len(stack) - 1		// index of top element in the stack
 		topOfStack := stack[ind]	
 		switch typ := topOfStack.(type) { // tokenType or nonTerminal
 		case tokenType:
 			if tokenList[tokenInd].typ == topOfStack {
				if tokenList[tokenInd+1].typ == tokenLpar && topOfStack == tokenAtom {
					temp := Term{Typ: relationMap[tokenList[tokenInd].typ], Literal: tokenList[tokenInd].literal}
					str := temp.String()
					if val, ok := termMap[str]; ok {
						stk1 = append(stk1, val) // - CHECK SYNTAX
					} else { 
						termMap[str] = &temp
						stk1 = append(stk1, temp)
					}
					stk2[stk_ptr] = append(stk2[stk_ptr], nil) // either way push empty into stk2, just to have equal level
					str_ptr++
				}  
				else if topOfStack == tokenRpar {
					if stk_ptr > 0 {
						temp_stk1 := stk1[stk_ptr] // functor term
						stk1 = stk1[:stk_ptr] 
						temp_stk2[stk_ptr] := stk2[stk_ptr] // list of args terms
						// TODO: pop stk2
					}
					temp := Term{Typ: TermCompound, Functor: temp_stk1 , Args: temp_stk2[]}
					str := temp.String() 
					if val, ok := termMap[str]; ok {
						// TODO: Pass val to final outcome, Grammar
					} else {
						// TODO: Pass temp to final outcome, Grammar
					}

					stk_ptr-- // pop pop
					if stk_ptr > 0 {
						for len(stk2[stk_ptr]) != 0 {
							stk2[stk_ptr] = appned(stk2[stk_ptr], temp_stk2[])
						}
					}
				} 

				else if tokenList[tokenInd+1].typ != tokenLpar && (topOfStack == tokenAtom || topOfStack == tokenNumber || topOfStack == tokenVariable)  { // general case
					temp := Term{Typ: relationMap[tokenList[tokenInd].typ], Literal: tokenList[tokenInd].literal} // 1. Create Term struct
					str := temp.String()
					if val, ok := termMap[str]; ok { // 2. Check duplicate from the termMap - if duplicate, just use the val and append to stk2
						stk2[stk_ptr] = append(stk2[stk_ptr], val) // - CHECK SYNTAX
					} else { // 3. if not, put new Term into termMap - if no duiplicate, use new Term to append to stk2
						termMap[str] = &temp
						stk2[stk_ptr] = append(stk2[stk_ptr], termMap[str]) // - CHECK SYNTAX
					}
				} 

				// when the top is terminal
 				// pop out the terminal and advance the tokenList index
 				// when the tokenType(terminal) are the same
 				// TODO: indicator for create compound 
 				// check if the topOfStack == tokenRpar 
 				// then do compound creation
 				// call helper function  createrCompund()
 				// TODO: indicator for pushing to stacks
 				// check if the topOfStack == tokenAtom && tokenList[tokenInd + 1] == tokenLpar
 				// Term atom being detect
 				// create the term for this atom. do pushing items to two stacks (stk1 - termAtom, skt2 - with empty_list( []*Term ) )
 				// create the term, push into top of stk2's list
 				// but remember check the map if exits already
 				stack = stack[:ind]		// pop out the top element
 				tokenInd += 1;
 			} else {
 				// terminal is not match
 				// TODO: double check here
 				return nil, ErrParser
 			}

 		case nonTermial:
 			// when the top is non terminal
 			// check the value in the parsing table with given token
 			if parseTable[topOfStack][tokenList[tokenInd]] != nil {
 				// value inside the cell, find the transition to other state
 				var transList = parseTable[topOfStack][tokenList[tokenInd]]
 				var listIndex = len(transList) -1

 				stack = stack[:ind]		// pop out the top non terminal before push
 				// push T -> X1 X2 X3 to the stack in reverse order
 				for listIndex >= 0 {
	 				stack = append(stack, transList[listIndex])
					listIndex -= 1
	 			}
 			} else {
 				// no value at given cell, invalid input string
 				return nil, ErrParser
 			}

 		default:
 			// invalid type in the stack
 			return nil, ErrParser
 		}
 	}

	return nil, nil
}

