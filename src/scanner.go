package src

import (
	"fmt"

	"github.com/caelondev/mutex/src/helpers"
	"github.com/caelondev/mutex/src/lexer"
)

type Scanner struct {
	SourceCode []rune
	Tokens []*lexer.Token

	Start int
	Current int
	Line int
}

func NewScanner(sourceCode string) *Scanner {
	return &Scanner{
		SourceCode: []rune(sourceCode),
		Tokens: []*lexer.Token{},
		Start: 0,
		Current: 0,
		Line: 1,
	}
}

func (s *Scanner) ScanTokens() []*lexer.Token {
	for !s.isEOF() {
		s.Start = s.Current
		s.ScanToken()
	}

	s.Tokens = append(s.Tokens, &lexer.Token{
		TokenType: lexer.EOF,
		Lexeme: "",
		Literal: nil,
		Line: s.Line,
	})

	return s.Tokens
}

func (s *Scanner) isEOF() bool {
	return s.Current >= len(s.SourceCode)
}

func (s *Scanner) ScanToken() {
	c := s.advance()

	switch c {
	case '(':
		s.addToken(lexer.LEFT_PARENTHESIS)
	case ')':
		s.addToken(lexer.RIGHT_PARENTHESIS)
	case '{':
		s.addToken(lexer.LEFT_BRACE)
	case '}':
		s.addToken(lexer.LEFT_BRACE)
	case ':':
		s.addToken(lexer.COLON)
	case ';':
		s.addToken(lexer.SEMICOLON)
	case ',':
		s.addToken(lexer.COMMA)
	case '.':
		s.addToken(lexer.DOT)
	case '-':
		s.addToken(lexer.MINUS)
	case '*':
		s.addToken(lexer.STAR)
	case '+':
		s.addToken(lexer.PLUS)
	case '<':
		s.addToken(helpers.Ternary(s.match('='), lexer.LESS_EQUAL, lexer.LESS).(lexer.TokenType))
	case '>':
		s.addToken(helpers.Ternary(s.match('='), lexer.GREATER_EQUAL, lexer.GREATER).(lexer.TokenType))
	case '=':
		s.addToken(helpers.Ternary(s.match('='), lexer.EQUAL_TO, lexer.ASSIGNMENT).(lexer.TokenType))
	case '!':
		s.addToken(helpers.Ternary(s.match('='), lexer.NOT_EQUAL, lexer.NOT).(lexer.TokenType))
	case '/':
		s.handleSlash()
	case '"':
		s.handleString()

	case ' ', '\r', '\t':
		// Ignore whitespace ---
		break
	case '\n':
		s.Line++

	default:
		mutex.reportError(s.Line, fmt.Sprintf("Unexpected token found: %c", c))
	}
}

func (s *Scanner) handleString() {
	for s.peek() != '"' && !s.isEOF() {
      if (s.peek() == '\n') {
				s.Line++
			}
      s.advance();
  }

	if (s.isEOF()) {
      mutex.reportError(s.Line, "Unterminated string.");
      return;
    }

    s.advance(); // Eat closing ".

    // Trim the surrounding quotes.
		value := s.SourceCode[s.Start + 1 : s.Current - 1]
    s.addTokenWithLiteral(lexer.STRING, value);
}

func (s *Scanner) handleSlash() {
	if s.match('/') { // Check another slash
		for s.peek() != '\n' && !s.isEOF() {
			s.advance() // Eat tokens until EOF or newline ---
		}
	} else {
		s.addToken(lexer.SLASH)
	}
}

func (s *Scanner) advance() rune {
	result := s.SourceCode[s.Current]
	s.Current++
	return result
}

func (s *Scanner) addToken(tokenType lexer.TokenType) {
	s.addTokenWithLiteral(tokenType, nil)
}

func (s *Scanner) addTokenWithLiteral(tokenType lexer.TokenType, literal any) {
	text := string(s.SourceCode[s.Start:s.Current])
	s.Tokens = append(s.Tokens, &lexer.Token{
		TokenType: tokenType,
		Lexeme: text,
		Literal: literal,
		Line: s.Line,
	})
}

func (s *Scanner) match(expected rune) bool {
	if s.isEOF() {
		return false
	}
	if s.SourceCode[s.Current] != expected {
		return false
	}

	s.Current++
	return true
}

func (s *Scanner) peek() rune {
	if s.isEOF() {
		return 0
	}
	return s.SourceCode[s.Current]
}

func (s *Scanner) peekNext() rune {
	return s.SourceCode[s.Current+1]
}

func (s *Scanner) isNumber(c rune) bool {
	return c >= '0' && c  <= '9'
}
