package transpiler

import (
	"fmt"
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

func transpileExpression(expr *parser.Node) string {
	switch expr.Kind {
	case parser.NODE_NUMERIC, parser.NODE_STRING, parser.NODE_ID:
		return expr.Value
	case parser.NODE_BINOP:
		return fmt.Sprintf("(%s %s %s)", transpileExpression(expr.Left), expr.Value, transpileExpression(expr.Right))
	case parser.NODE_UNOP:
		return fmt.Sprintf("(%s%s)", expr.Value, transpileExpression(expr.Right))
	case parser.NODE_FUNCALL:
		args := ""
		for i, arg := range expr.List {
			if i > 0 {
				args += ","
			}
			args += transpileExpression(arg)
		}
		return fmt.Sprintf("%s(%s)", expr.Value, args)
	default:
		return ""
	}
}

func transpileVarStatement(expr *parser.Node) string {
	return fmt.Sprintf("local %s = %s", expr.Value, transpileExpression(expr.Right))
}

func (i *Instance) Transpile(filename string) Module {
	text := ""
	for i.ip < len(i.input) {
		stmt := i.input[i.ip]
		switch stmt.Kind {
		case parser.NODE_VAR:
			text += transpileVarStatement(stmt) + "\n"
		default:
			text += transpileExpression(stmt) + "\n"
		}
		i.ip++
	}
	return Module{
		Filename: filename,
		Text:     text,
	}
}
