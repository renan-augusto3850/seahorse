package transpiler

import (
	"fmt"
	"os"
	"seahorse/parser"
)

type Instance struct {
	input []*parser.Node
	ip    int // Instruction Pointer
}

func NewInstance(input []*parser.Node) *Instance {
	return &Instance {
		input: input,
		ip: 0,
	}
}

type Module struct {
	Filename string
	Text     string
}

func transpileFunction(funcall *parser.Node) string {
	function_name := funcall.Value
	switch function_name {
	case "printLog":
		if len(funcall.List) != 1 {
			fmt.Printf("printLog accepts only 1 argument but got %d\n", len(funcall.List))
			os.Exit(1)
		}
		arg := funcall.List[0]
		return fmt.Sprintf("print(%s)", transpileExpression(arg))
	case "input":
		if len(funcall.List) != 0 {
			fmt.Printf("input accepts no arguments but got %d\n", len(funcall.List))
			os.Exit(1)
		}
		return fmt.Sprintf("io.read()")
	default:
		fmt.Printf("Unknown function %s\n", function_name)
		os.Exit(1)
	}
	return "nil"
}

func transpileExpression(expr *parser.Node) string {
	switch expr.Kind {
	case parser.NODE_NUMERIC, parser.NODE_ID:
		return expr.Value
	case parser.NODE_STRING:
		return fmt.Sprintf("\"%s\"", expr.Value)
	case parser.NODE_BINOP:
		return fmt.Sprintf("(%s %s %s)", transpileExpression(expr.Left), expr.Value, transpileExpression(expr.Right))
	case parser.NODE_UNOP:
		return fmt.Sprintf("(%s%s)", expr.Value, transpileExpression(expr.Right))
	case parser.NODE_FUNCALL:
		return transpileFunction(expr)
	default:
		return ""
	}
}

func transpileVarStatement(expr *parser.Node) string {
	return fmt.Sprintf("local %s = %s", expr.Value, transpileExpression(expr.Right))
}

func transpileIfStatement(expr *parser.Node) string {
	return fmt.Sprintf("if %s then\n%send", transpileExpression(expr.Right), transpileBlock(expr.List))
}

func transpileBlock(block []*parser.Node) string {
	text := ""
	for i := range block {
		stmt := block[i]
		switch stmt.Kind {
		case parser.NODE_VAR:
			text += transpileVarStatement(stmt) + "\n"
		case parser.NODE_IF:
			text += transpileIfStatement(stmt) + "\n"
		default:
			text += transpileExpression(stmt) + "\n"
		}
	}
	return text
}

func (i *Instance) Transpile(filename string) Module {
	text := transpileBlock(i.input)
	return Module{
		Filename: filename,
		Text:     text,
	}
}
