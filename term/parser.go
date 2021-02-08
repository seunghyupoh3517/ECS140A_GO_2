package term

import (
	"errors"
	// "fmt"
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

var mixedArray = [][][]interface{} {{nil, nil, nil, nil, {Term_NT, tokenEOF}, {Term_NT, tokenEOF}, {Term_NT, tokenEOF}}, {nil, nil , nil, nil, {tokenAtom, NT1}, {tokenNumber}, {tokenVariable}}, {{}, {tokenLpar, Args_NT, tokenRpar}, {}, {},  nil, nil, nil},  {nil, nil, nil, nil, {Term_NT, NT2}, {Term_NT, NT2}, {Term_NT, NT2}}, {nil, nil, {}, {tokenComma, Args_NT}, nil, nil, nil}}



// implements the Parse method with Grammar struct
func (g Grammar) Parse(str string) (*Term, error) {
	// TODO: matrix in the global
	// The parseTable can be the type of
	// []interface{} which is the list in a single cell
	var finalTerm = &Term{}
	parseTable := mixedArray
	var stk1 = []*Term{} 					// functor
	var stk2 = [][]*Term{}   				// argument List
	var stk_ptr = 0
	var termMap = map[string]*Term{}  	// term.toString() -> *term


	// Test if the given is empty
	// emplex := newLexer(str)
	// emp, _:= emplex.next()
	// if emp.typ == tokenEOF{
	// 	fmt.Println("gun - line 82")
	// 	return nil, nil
	// }

	// Tokennize the input string
	lex := newLexer(str)
	var	tokenList []*Token
	for {
		token, err := lex.next()

		if err == ErrLexer {
			// validating the given string, return error if can't parse to token
			// TODO: double check with error object return type
			// fmt.Println("***** Error from line 84 ")
			// fmt.Println("ni gan ha? - line 95")
			return nil, ErrParser
			} else {
				if token.typ == tokenEOF {
					tokenList = append(tokenList, token)
					break
				}
				tokenList = append(tokenList, token)
			}
	}

	// Handle corner case when the input is an empty string ""
	// if len(tokenList) == 1 && tokenList[0].typ == tokenEOF{
	// 	return nil, nil
	// }


	// pointer point to the current token in the list
	var tokenInd = 0
	// initialize the stack
	// stack needs to accept two data types, nonTerminal & tokenType
	var stack []interface{}
	stack = append(stack, Start_NT)

	if len(tokenList) != 1 && tokenList[0].typ != tokenEOF {
	 	for len(stack) != 0 {

	 		ind := len(stack) - 1		// index of top element in the stack
	 		//fmt.Println("************* Index checking Line 111")
	 		topOfStack := stack[ind]


	 		//fmt.Println("************* Index error Line 113")

	//		fmt.Println("==================================")
	//		fmt.Println("Token type:", tokenList[tokenInd].typ, " Token literal:", tokenList[tokenInd].literal)
	//		fmt.Println("The topOfStack type is:", topOfStack)
	//		fmt.Println("==================================")

	 		switch typ := topOfStack.(type) { // tokenType or nonTerminal
	 		case tokenType:
	// 			fmt.Println("==================================")
	//			fmt.Println("********** Line 120: bottle of stack:", stack[0])
	//			fmt.Println("********** Line 123: token size:", len(tokenList))
	//			fmt.Println("********** Line 123: token Ind:", tokenInd)
	//			fmt.Println("==================================")

				// fmt.Println("==================================")
				// fmt.Println("Token type:", tokenList[tokenInd].typ, " Token literal:", tokenList[tokenInd].literal)
				// fmt.Println("==================================")

	 			if tokenList[tokenInd].typ == topOfStack {
					if topOfStack == tokenAtom && tokenList[tokenInd + 1].typ == tokenLpar {
						// indicator for create the functor term and push items to two stacks
						temp := &Term{Typ: relationMap[tokenList[tokenInd].typ], Literal: tokenList[tokenInd].literal}
						str := temp.String()
						if val, ok := termMap[str]; ok {
							stk1 = append(stk1, val) 		// - CHECK SYNTAX
						} else {
							termMap[str] = temp
							stk1 = append(stk1, temp)
						}

						var tempList = []*Term{}			// arguments list
						stk2 = append(stk2, tempList)		// push empty term[] to the stk2
						stk_ptr++
					} else if topOfStack == tokenRpar {
						// indicator for creating the compound Term
						if stk_ptr > 0 {
							// create the compound term
							temp := &Term{Typ: TermCompound, Functor: stk1[stk_ptr - 1] , Args: stk2[stk_ptr - 1]}

							// pop out the top of two stacks
							stk_ptr--
							stk1 = stk1[:stk_ptr]
							stk2 = stk2[:stk_ptr]

							// check if exits in the termMap avoid duplicate compound
							str := temp.String()
							if val, ok := termMap[str]; ok {
								temp = val 		// use the old *term if exits in the map
							} else {
								// put the new compound in the termMap
								termMap[str] = temp
							}

							// append the new created compound into the next level
							if stk_ptr > 0 {
								stk2[stk_ptr - 1] = append(stk2[stk_ptr - 1], temp)
							} else {
								// we create the last final compound
								//fmt.Println("************ Check Line 174:", temp.String())
								//return temp, nil
								finalTerm = temp
							}
						}

					} else if (topOfStack == tokenAtom || topOfStack == tokenNumber || topOfStack == tokenVariable)  {
						// general case
						temp := &Term{Typ: relationMap[tokenList[tokenInd].typ], Literal: tokenList[tokenInd].literal} // 1. Create Term struct
						str := temp.String()

	//					fmt.Println("*****************************")
	//					fmt.Println(str)
	//					fmt.Println(temp.Typ)
	//					fmt.Println(temp.Args)
	//					fmt.Println("*****************************")

						if val, ok := termMap[str]; ok {
							temp = val
						} else { 										// 3. if not, put new Term into termMap - if no duiplicate, use new Term to append to stk2
							termMap[str] = temp

						}
						finalTerm = temp


						if stk_ptr > 0 {
							//fmt.Println("************* Index checking Line 188")
							stk2[stk_ptr - 1] = append(stk2[stk_ptr - 1], temp)
							//fmt.Println("************* Index error Line 190")
						}
					} else if (topOfStack == tokenEOF) {
						//fmt.Println("******************** I reach the end of input token *******")
							// only a single term left, return it
						 	if len(termMap) > 0 {
							 	return finalTerm, nil
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
	 				//fmt.Println("******* Line 231 stack size:", len(stack))
	 				tokenInd += 1;
	 			} else {
	 				// terminal is not match
	 				// TODO: double check here
	 				// fmt.Println("***** Error from line 184 ")
	 				return nil, ErrParser
	 			}

	 		case nonTermial:
	 			//fmt.Println("############### cao !!")
	 			//fmt.Println("row:",topOfStack, "coloum:", tokenList[tokenInd].typ)
	 			// when the top is non terminal
	 			// check the value in the parsing table with given token
	 			if parseTable[typ][tokenList[tokenInd].typ] != nil {
	 				// value inside the cell, find the transition to other state
	 				var transList = parseTable[typ][tokenList[tokenInd].typ]
	 				var listIndex = len(transList) -1

	 				//fmt.Println("$$$$$$$$$$$$$$$$$$ gan!!!")

	 				stack = stack[:ind]		// pop out the top non terminal before push
	 				// push T -> X1 X2 X3 to the stack in reverse order
	 				for listIndex >= 0 {
		 				stack = append(stack, transList[listIndex])
						listIndex -= 1
		 			}
	 			} else {
	 				// no value at given cell, invalid input string
	 				//fmt.Println("@@@@@@@@@@@@@@@ cao!!!!")
	 				// fmt.Println("***** Error from line 208 ")
	 				return nil, ErrParser
	 			}

	 		// default:
	 			// invalid type in the stack
	 			//fmt.Println(typ)
	 			// fmt.Println("***** Error from line 214 ")
	 			// return nil, ErrParser
				 	// return nil, nil
	 		}
	 	}
	}
	
	return nil, nil
}
