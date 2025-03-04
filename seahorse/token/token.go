package token

const (
	TOKEN_UNK = iota
	
	TOKEN_ID      = iota // identifier
	TOKEN_STRING  = iota // "string"
	TOKEN_NUMERIC = iota // numeric
	TOKEN_EOF     = iota // end of file
	TOKEN_LF      = iota // newline character

	// Keyword
	TOKEN_VAR = iota
	TOKEN_IF  = iota

	// Delimiter
	TOKEN_OPENPAREN    = iota // (
	TOKEN_CLOSEPAREN   = iota // )
	TOKEN_OPENCURLY    = iota // {
	TOKEN_CLOSECURLY   = iota // }
	TOKEN_OPENBRACKET  = iota // [
	TOKEN_CLOSEBRACKET = iota // ]
	TOKEN_COMMA        = iota // ,

	// Operation
	TOKEN_PLUS   = iota // +
	TOKEN_MINUS  = iota // -
	TOKEN_MUL    = iota // *
	TOKEN_DIV    = iota // /
	TOKEN_MOD    = iota // %
	TOKEN_ASSIGN = iota // =
	TOKEN_EQ     = iota // ==
	TOKEN_NE     = iota // !=
	TOKEN_GT     = iota // >
	TOKEN_LT     = iota // <
	TOKEN_GE     = iota // >=
	TOKEN_LE     = iota // <=
	TOKEN_NOT    = iota // !
)

type Token struct {
	Kind  int
	Value string
	Filename string
	Line int
	Col int
}
