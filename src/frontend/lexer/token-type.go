package lexer

type TokenType int

const (

	// Single-char token ---
	LEFT_PARENTHESIS TokenType = iota
	RIGHT_PARENTHESIS
	LEFT_BRACE
	RIGHT_BRACE
	COMMA
	DOT
	MINUS
	PLUS
	COLON
	SEMICOLON
	SLASH
	STAR
	MODULO

	// One or two char tokens ---
	NOT
	NOT_EQUAL
	EQUAL_TO
	GREATER
	GREATER_EQUAL
	LESS
	LESS_EQUAL
	ASSIGNMENT

	// Literals ---
	IDENTIFIER
	STRING
	NUMBER

	// Keywords
	AND
	CLASS
	ELSE
	FALSE
	TRUE
	FUNCTION
	IF
	NIL
	OR
	MUTABLE
	IMMUTABLE
	PRINT
	RETURN
	SUPER
	THIS
	VAR
	WHILE

	EOF
)

var RESERVED_KEYWORDS map[string]TokenType = map[string]TokenType{
	"and": AND,
	"class": CLASS,
	"else": ELSE,
	"false": FALSE,
	"true": TRUE,
	"fn": FUNCTION,
	"if": IF,
	"nil": NIL,
	"or": OR,
	"print": PRINT,
	"return": RETURN,
	"super": SUPER,
	"this": THIS,
	"mut": MUTABLE,
	"imm": IMMUTABLE,
	"var": VAR,
	"while": WHILE,
}

func TokenTypeString(t TokenType) string {
	switch t {
	case LEFT_PARENTHESIS:
		return "LEFT_PARENTHESIS"
	case RIGHT_PARENTHESIS:
		return "RIGHT_PARENTHESIS"
	case LEFT_BRACE:
		return "LEFT_BRACE"
	case RIGHT_BRACE:
		return "RIGHT_BRACE"
	case COMMA:
		return "COMMA"
	case DOT:
		return "DOT"
	case MINUS:
		return "MINUS"
	case PLUS:
		return "PLUS"
	case SEMICOLON:
		return "SEMICOLON"
	case SLASH:
		return "SLASH"
	case STAR:
		return "STAR"

	// one/two char
	case NOT:
		return "NOT"
	case NOT_EQUAL:
		return "NOT_EQUAL"
	case EQUAL_TO:
		return "EQUAL_TO"
	case GREATER:
		return "GREATER"
	case GREATER_EQUAL:
		return "GREATER_EQUAL"
	case LESS:
		return "LESS"
	case LESS_EQUAL:
		return "LESS_EQUAL"
	case ASSIGNMENT:
		return "ASSIGNMENT"

	// literals
	case IDENTIFIER:
		return "IDENTIFIER"
	case STRING:
		return "STRING"
	case NUMBER:
		return "NUMBER"

	// keywords
	case AND:
		return "AND"
	case CLASS:
		return "CLASS"
	case ELSE:
		return "ELSE"
	case FALSE:
		return "FALSE"
	case TRUE:
		return "TRUE"
	case FUNCTION:
		return "FUNCTION"
	case IF:
		return "IF"
	case NIL:
		return "NIL"
	case OR:
		return "OR"
	case PRINT:
		return "PRINT"
	case RETURN:
		return "RETURN"
	case SUPER:
		return "SUPER"
	case THIS:
		return "THIS"
	case MUTABLE:
		return "MUTABLE"
	case IMMUTABLE:
		return "IMMUTABLE"
	case VAR:
		return "VAR"
	case WHILE:
		return "WHILE"

	case EOF:
		return "EOF"

	default:
		return "UNKNOWN_TOKEN"
	}
}
