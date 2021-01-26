package branch

import (
	"go/ast"
	"go/parser"
	"go/token"
)
//func Inspect(node Node, f func(Node) bool)
//Inspect traverses an AST in depth-first order: It starts by calling f(node); node must not be nil. 
//If f returns true, 
//Inspect invokes f recursively for each of the non-nil children of node, followed by a call of f(nil).
//https://golang.org/pkg/go/ast/#Inspect
func branchCount(fn *ast.FuncDecl) uint {
	// TODO: Write the branchCount function,
	// count the number of branching statements in function fn
	//https://play.golang.org/p/cq2OI6CA6v_n
	var Count uint = 0
	
	ast.Inspect(fn, func (node ast.Node) bool{
		//10-15lines
		//detect if for switch range, type switch, goto, continue, break, fallthrough
		// If we return true, we keep recursing under this AST node.
        // If we return false, we won't visit anything under this AST node.
		switch node.(type) {
		case *ast.IfStmt: //find if
			Count++
		case *ast.ForStmt: //find for
			Count++
		case *ast.SwitchStmt: //find switch
			Count++
		case *ast.RangeStmt: //find Range
			Count++
		case *ast.TypeSwitchStmt: //find type
			Count++
		case *ast.BranchStmt: //include goto, continue, break, fallthrough
			Count++
		// default:
		// 	return false
		}
	
	return true
})

	return Count
}

// ComputeBranchFactors returns a map from the name of the function in the given
// Go code to the number of branching statements it contains.
func ComputeBranchFactors(src string) map[string]uint {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "src.go", src, 0)
	if err != nil {
		panic(err)
	}

	m := make(map[string]uint)
	for _, decl := range f.Decls {
		switch fn := decl.(type) {
		case *ast.FuncDecl:
			m[fn.Name.Name] = branchCount(fn)
		}
	}

	return m
}
