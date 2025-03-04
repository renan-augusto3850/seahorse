package main

import (
	"fmt"
    "io/ioutil"
	"seahorse/lexer"
	"seahorse/parser"
	"seahorse/transpiler"
	"os"
)

func main() {
	args := os.Args
	if len(args) < 2 {
		fmt.Printf("Usage: %s FILENAME\n", args[0])
		fmt.Println("Positional argument FILENAME required.")
		return
	}

	input_file, err := ioutil.ReadFile(args[1])
    if err != nil {
        fmt.Println("Could not read input file.")
   		return
    }

    input := string(input_file)
	tokens := lexer.Lexer(input, args[1])
	p := parser.NewParser(tokens)
	ast := p.Parse()
	if ast == nil {
		err := p.Err;
		fmt.Printf("Syntax error at %s:%d:%d near %s\n",
					err.Filename, err.Line, err.Col, err.Near);
		fmt.Printf("    %s\n", err.Message)
		return
	}
	i := transpiler.NewInstance(ast)

	x := i.Transpile(fmt.Sprintf("%s.lua", args[1]));

	// Write transpiled result to file
	file, err := os.Create(x.Filename)
	if err != nil {
		fmt.Println("Could not write output file.")
		return
	}
	file.Write([]byte(x.Text))
}

// func printAST(node *parser.Node, indent int) {
// 	if node == nil { return }
// 	for i := 0; i < indent; i++ {
// 		fmt.Print("  ")
// 	}
// 	switch node.Kind {
// 	case parser.NODE_BINOP: fmt.Printf("Binary Operation: %s\n", node.Value)
// 	case parser.NODE_UNOP: fmt.Printf("Unary Operation: %s\n", node.Value)
// 	case parser.NODE_ID: fmt.Printf("Identifier: %s\n", node.Value)
// 	case parser.NODE_FUNCALL: fmt.Printf("Function Call: %s\n", node.Value)
// 	case parser.NODE_NUMERIC: fmt.Printf("Numeric: %s\n", node.Value)
// 	case parser.NODE_STRING: fmt.Printf("String: \"%s\"\n", node.Value)
// 	case parser.NODE_VAR: fmt.Printf("Var: \"%s\"\n", node.Value)
// 	}

// 	printAST(node.Left, indent+1)
// 	printAST(node.Right, indent+1)
// 	printASTList(node.List, indent+1)
// }

// func printASTList(list []*parser.Node, ident int) {
// 	for i := 0; i < len(list); i++ {
// 		fmt.Printf("\t--- %d ---\n", i);
// 		printAST(list[i], ident);
// 	}
// }
