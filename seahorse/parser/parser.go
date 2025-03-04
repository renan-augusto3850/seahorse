package parser

import (
	"fmt"
	"seahorse/token"
)

const (
	NODE_BINOP   = iota // a + b, b / (a + x)
	NODE_UNOP    = iota // !a, -b
	NODE_NUMERIC = iota // any numeric value
	NODE_STRING  = iota // "string"
	NODE_ID      = iota // foo, bar, baz
	NODE_FUNCALL = iota // foo(), bar(), baz()
	NODE_VAR     = iota // var foo = bar
	NODE_IF      = iota // if expr { block }
)

// Error function to log syntax errors.
func (p *Parser) ParseError(message string) {
	t := p.currentToken()
	value := t.Value
	if t.Kind == token.TOKEN_EOF { value = "End of file" }
	if t.Kind == token.TOKEN_LF { value = "Newline (\\n)" }
	msg := fmt.Sprintf("%s", message)
	p.Err = ParseError {
		Message: msg,
		Near: value,
		Filename: t.Filename,
		Line: t.Line,
		Col: t.Col,
	}
}

// Node is the universal construction that can hold any
// type of structure parsed by Seahorse. Depending on
// the field Kind, the other fields are used in different ways
type Node struct {
	Kind  int
	Value string
	Left  *Node
	Right *Node
	List []*Node
}

type Parser struct {
	tokens []token.Token
	pos    int
	Err    ParseError
}

type ParseError struct {
	Message  string
	Near     string
	Filename string
	Line     int
	Col      int
}

// Creates new Parser instance.
func NewParser(tokens []token.Token) *Parser {
	return &Parser{
		tokens: tokens,
		pos:    0,
	}
}

// Returns the token at the pointer of the parser.
func (p *Parser) currentToken() token.Token {
	if p.pos >= len(p.tokens) {
		return token.Token{Kind: token.TOKEN_EOF}
	}
	return p.tokens[p.pos]
}

// Advances the pointer of the parser.
func (p *Parser) consume() {
	p.pos++
}

var binaryOperation = []int{
	token.TOKEN_PLUS,
	token.TOKEN_MINUS,
	token.TOKEN_MUL,
	token.TOKEN_DIV,
	token.TOKEN_EQ,
	token.TOKEN_NE,
	token.TOKEN_LT,
	token.TOKEN_GT,
	token.TOKEN_GE,
	token.TOKEN_LE,
}
// var unaryOperation = []int{ token.TOKEN_NOT, token.TOKEN_MINUS }

// Returns precedence according to proper binary operation Token Kind.
func getBinaryOperationPrecedence(kind int) int {
	switch kind {
	case token.TOKEN_EQ: return 2;
	case token.TOKEN_NE: return 2;
	case token.TOKEN_GE: return 2;
	case token.TOKEN_LT: return 2;
	case token.TOKEN_GT: return 2;
	case token.TOKEN_LE: return 2;
	case token.TOKEN_PLUS: return 3;
	case token.TOKEN_MINUS: return 3;
	case token.TOKEN_MUL: return 4;
	case token.TOKEN_DIV: return 4;
	}
	return -1;
}

// Returns precedence according to proper unary operation Token Kind.
func getUnaryOperationPrecedence(kind int) int {
	switch kind {
	case token.TOKEN_NOT: return 5;
	case token.TOKEN_MINUS: return 5;
	}
	return -1;
}

// Parse expressions separated by ,
func (p *Parser) parseArgumentList() []*Node {
	list := []*Node{}
	first := p.parseExpression(0)
	if first != nil {
		list = append(list, first)
	}
	for p.expectTokens([]int{ token.TOKEN_COMMA }) {
		if p.currentToken().Kind == token.TOKEN_CLOSEPAREN {
			break
		}
		p.consume()
		x := p.parseExpression(0)
		if x == nil {
			p.ParseError("Expected expression after ','")
			return nil
		}
		list = append(list, x)
	}
	return list
}

// Parse arbitrary expression
func (p *Parser) parseExpression(precedence int) *Node {
	node := p.parseTerm()
	if node == nil { return nil }
	if p.currentToken().Kind == token.TOKEN_OPENPAREN && node.Kind == NODE_ID {
		p.consume()
		list := p.parseArgumentList()
		if list == nil { return nil }
		if !p.expectToken(token.TOKEN_CLOSEPAREN) {
			p.ParseError("Expected ')' after argument list.")
			return nil
		}
		p.consume()
		node = &Node{
			Kind:  NODE_FUNCALL,
			Value: node.Value,
			List:  list,
		}
	} else {
		for p.expectTokens(binaryOperation) {
			next_prec := getBinaryOperationPrecedence(p.currentToken().Kind)
			if next_prec < precedence {
				break
			}
			operator := p.currentToken()
			p.consume()
			rightNode := p.parseExpression(next_prec)
			if rightNode == nil { return nil }
			node = &Node{
				Kind:  NODE_BINOP,
				Value: operator.Value,
				Left:  node,
				Right: rightNode,
			}
		}
	}
	return node
}

// Returns if a specific Token Kind is the current token's kind
func (p *Parser) expectToken(kind int) bool {
	return p.currentToken().Kind == kind
}

// Returns if one of the provided token kinds is the current token's kind
func (p *Parser) expectTokens(kinds []int) bool {
	var expects []bool;
	for i := 0; i < len(kinds); i++ {
		result := p.currentToken().Kind == kinds[i]
		expects = append(expects, result);
	}
	for j := 0; j < len(expects); j++ {
		if expects[j] {
			return true;
		}
	}
	return false;
}

// Parse basic piece of an expression.
func (p *Parser) parseTerm() *Node {
	t := p.currentToken()
	switch t.Kind {
	case token.TOKEN_NUMERIC:
		p.consume()
		return &Node{
			Kind:  NODE_NUMERIC,
			Value: t.Value,
		}
	case token.TOKEN_NOT:
		p.consume()
		return &Node{
			Kind:  NODE_UNOP,
			Value: t.Value,
			Right: p.parseExpression(getUnaryOperationPrecedence(token.TOKEN_NOT)),
		}
	case token.TOKEN_MINUS:
		p.consume()
		return &Node{
			Kind:  NODE_UNOP,
			Value: t.Value,
			Right: p.parseExpression(getUnaryOperationPrecedence(token.TOKEN_MINUS)),
		}
	case token.TOKEN_OPENPAREN:
		p.consume()
		node := p.parseExpression(0)
		if !p.expectToken(token.TOKEN_CLOSEPAREN) {
			p.ParseError("Expected ')' after expression.")
			return nil
		}
		return node
	case token.TOKEN_ID:
		p.consume()
		return &Node{
			Kind:  NODE_ID,
			Value: t.Value,
		}
	case token.TOKEN_STRING:
		p.consume()
		return &Node{
			Kind:  NODE_STRING,
			Value: t.Value,
		}
	case token.TOKEN_LF:
		p.consume()
		return p.parseTerm()
	}
	p.ParseError(fmt.Sprintf("Unexpected token: %d", t.Kind))
	return nil
}

// Parse var statement
func (p *Parser) parseVarStatement() *Node {
	p.expectToken(token.TOKEN_VAR)
	p.consume()
	if !p.expectToken(token.TOKEN_ID) {
		p.ParseError("Expected identifier")
		return nil
	}
	name := p.currentToken().Value
	p.consume()
	if !p.expectToken(token.TOKEN_ASSIGN) {
		p.ParseError("Expected =")
		return nil
	}
	p.consume()
	expr := p.parseExpression(0)
	if expr == nil {
		p.ParseError("Expected expression")
		return nil
	}
	if !p.expectEnd() {
		p.ParseError("Expected end of statement")
		return nil
	}
	p.consume()
	return &Node {
		Kind: NODE_VAR,
		Value: name,
		Right: expr,
	}
}

// Parse if statement
func (p *Parser) parseIfStatement() *Node {
	p.expectToken(token.TOKEN_IF)
	p.consume()
	expr := p.parseExpression(0)
	p.expectToken(token.TOKEN_OPENCURLY)
	p.consume()
	block := p.Parse(true)
	if block == nil { return nil }
	p.expectToken(token.TOKEN_CLOSECURLY)
	p.consume()
	return &Node {
		Kind: NODE_IF,
		Right: expr,
		List: block,
	}
}

func (p *Parser) expectEnd() bool {
	return p.expectTokens([]int{token.TOKEN_EOF, token.TOKEN_LF})
}

// Parse
func (p *Parser) Parse(inside_block bool) []*Node {
	nodes := []*Node{}
	for p.currentToken().Kind != token.TOKEN_EOF {
		switch p.currentToken().Kind {
		case token.TOKEN_VAR:
			x := p.parseVarStatement()
			if x == nil { return nil }
			nodes = append(nodes, x)
		case token.TOKEN_CLOSECURLY:
			if !inside_block {
				p.ParseError("Unexpected }")
				return nil
			}
			p.consume()
		case token.TOKEN_IF:
			x := p.parseIfStatement()
			if x == nil { return nil }
			nodes = append(nodes, x)
		case token.TOKEN_LF:
			p.consume()
			continue
		default:
			x := p.parseExpression(0)
			if x == nil { return nil }
			if !p.expectEnd() {
				p.ParseError("Expected end of statement")
				return nil
			}
			p.consume()
			nodes = append(nodes, x)
		}
	}
	return nodes
}
