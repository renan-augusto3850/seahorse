package lexer

import (
	"seahorse/token"
	"unicode"
)

func Lexer(input string, filename string) []token.Token {
	var tree []token.Token
	code := input
	line := 1;
	col := 1;
	for i := 0; i < len(code); {
		c := rune(code[i])
		cStr := string(c)
		switch {
		case c == '\n':
			tree = append(tree, token.Token{Kind: token.TOKEN_LF, Value: "\n", Filename: filename, Line: line, Col: col})
			line += 1;
			col = 1;
		case unicode.IsSpace(c):
			i++
			col++;
			continue
		case unicode.IsDigit(c): // Números
			start := i
			for i < len(code) && unicode.IsDigit(rune(code[i])) {
				i++
				col++;
			}
			tree = append(tree, token.Token{Kind: token.TOKEN_NUMERIC, Value: code[start:i], Filename: filename, Line: line, Col: col})
			continue
		case c == '"':
			var value string
			i++ // Pula a primeira aspas
			for i < len(input) && input[i] != '"' {
				value += string(input[i])
				i++
				col++;
			}
			tree = append(tree, token.Token{Kind: token.TOKEN_STRING, Value: value, Filename: filename, Line: line, Col: col})
		case cStr == "+": // Operadores
			tree = append(tree, token.Token{Kind: token.TOKEN_PLUS, Value: cStr, Filename: filename, Line: line, Col: col})
		case cStr == "-":
			tree = append(tree, token.Token{Kind: token.TOKEN_MINUS, Value: cStr, Filename: filename, Line: line, Col: col})
		case cStr == "*":
			tree = append(tree, token.Token{Kind: token.TOKEN_MUL, Value: cStr, Filename: filename, Line: line, Col: col})
		case cStr == "/":
			tree = append(tree, token.Token{Kind: token.TOKEN_DIV, Value: cStr, Filename: filename, Line: line, Col: col})
		case cStr == "=":
			if i+1 < len(input) && input[i+1] == '=' {
				tree = append(tree, token.Token{Kind: token.TOKEN_EQ, Value: "==", Filename: filename, Line: line, Col: col})
				i++
				col++
			} else {
				tree = append(tree, token.Token{Kind: token.TOKEN_ASSIGN, Value: cStr, Filename: filename, Line: line, Col: col})
			}
		case cStr == "!":
			if i+1 < len(input) && input[i+1] == '=' {
				tree = append(tree, token.Token{Kind: token.TOKEN_NE, Value: "!=", Filename: filename, Line: line, Col: col})
				i++
				col++
			} else {
				tree = append(tree, token.Token{Kind: token.TOKEN_NOT, Value: cStr, Filename: filename, Line: line, Col: col})
			}
		case cStr == "<":
			if i+1 < len(input) && input[i+1] == '=' {
				tree = append(tree, token.Token{Kind: token.TOKEN_LE, Value: "<=", Filename: filename, Line: line, Col: col})
				i++
				col++
			} else {
				tree = append(tree, token.Token{Kind: token.TOKEN_LT, Value: cStr, Filename: filename, Line: line, Col: col})
			}
		case cStr == ">":
			if i+1 < len(input) && input[i+1] == '=' {
				tree = append(tree, token.Token{Kind: token.TOKEN_GE, Value: ">=", Filename: filename, Line: line, Col: col})
				i++
				col++
			} else {
				tree = append(tree, token.Token{Kind: token.TOKEN_GT, Value: cStr, Filename: filename, Line: line, Col: col})
			}
		case cStr == ",":
			tree = append(tree, token.Token{Kind: token.TOKEN_COMMA, Value: cStr, Filename: filename, Line: line, Col: col})
		case cStr == "(":
			tree = append(tree, token.Token{Kind: token.TOKEN_OPENPAREN, Value: cStr, Filename: filename, Line: line, Col: col})
		case cStr == ")":
			tree = append(tree, token.Token{Kind: token.TOKEN_CLOSEPAREN, Value: cStr, Filename: filename, Line: line, Col: col})
		case unicode.IsLetter(c):
			start := i
			for i < len(code) && (unicode.IsLetter(rune(code[i])) || unicode.IsDigit(rune(code[i]))) {
				i++
				col++
			}
			value := code[start:i]
			switch value {
			case "var":
				tree = append(tree, token.Token{Kind: token.TOKEN_VAR, Value: value, Filename: filename, Line: line, Col: col})
			case "if":
				tree = append(tree, token.Token{Kind: token.TOKEN_IF, Value: value, Filename: filename, Line: line, Col: col})
			default:
				tree = append(tree, token.Token{Kind: token.TOKEN_ID, Value: value, Filename: filename, Line: line, Col: col})
				continue
			}
		default:
			tree = append(tree, token.Token{Kind: token.TOKEN_UNK, Value: cStr, Filename: filename, Line: line, Col: col})
		}
		i++
		col++
	}
	tree = append(tree, token.Token{Kind: token.TOKEN_EOF, Value: "", Filename: filename, Line: line, Col: col})
	return tree
}
