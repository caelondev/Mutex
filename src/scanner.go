package src

import (
	"fmt"
	"strconv"

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
		Tokens: make([]*lexer.Token, 0, len(sourceCode) + 1), // Over-allocate for faster tokenization process
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
		s.addToken(lexer.RIGHT_BRACE)
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
		if isNumber(c) {
			s.handleNumber()
		} else if isAlphabet(c) {	
			s.handleIdentifier()
		} else {
			mutex.reportError(s.Line, fmt.Sprintf("Unexpected token found: %c", c))
		}
	}
}

func (s *Scanner) handleIdentifier() {
	for isAlphanumeric(s.peek()) {
		s.advance()
	}


	keyword := string(s.SourceCode[s.Start : s.Current])
	if tokenType, ok := lexer.RESERVED_KEYWORDS[keyword]; ok {
		s.addToken(tokenType)
		return
	}


	s.addToken(lexer.IDENTIFIER)
}

func (s *Scanner) handleNumber() {
    // Eat digits before the dot
    for isNumber(s.peek()) {
        s.advance()
    }

    // Handle decimal part
    if s.peek() == '.' {
        if !isNumber(s.peekNext()) { // no number after dot
            mutex.reportError(s.Line,
                fmt.Sprintf("Expected number after '.' but got '%c'", s.peekNext()))
            return
        }

        s.advance() // Eat the dot

        for isNumber(s.peek()) {
            s.advance() // Eat digits after dot
        }
    }

    // Parse the number string
    value := string(s.SourceCode[s.Start:s.Current])
    parsedNumber, err := strconv.ParseFloat(value, 64)
    if err != nil {
        mutex.reportError(s.Line,
            fmt.Sprintf("Failed to parse number '%s': %s", value, err))
        return
    }

    s.addTokenWithLiteral(lexer.NUMBER, parsedNumber)
}

func (s *Scanner) handleString() {
	for s.peek() != '"' && !s.isEOF() {
      if (s.peek() == '\n') {
				currentValue := string(
					s.SourceCode[s.Start + 1 : s.Current],
				)

				mutex.reportError(s.Line, fmt.Sprintf("Missing closing string ('\"') after string value \"%s\"", currentValue))
				return
			}
      s.advance();
  }

	if (s.isEOF()) {
		currentValue := string(
			s.SourceCode[s.Start + 1 : s.Current],
		)
		mutex.reportError(s.Line, fmt.Sprintf("Missing closing string ('\"') after string value \"%s\"", currentValue))
		return;
	}

	s.match('"'); // Eat closing ".

	// Trim the surrounding quotes.
	value := string(
		s.SourceCode[s.Start + 1 : s.Current - 1],
		)

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
    if s.Current+1 >= len(s.SourceCode) {
        return 0 
    }
    return s.SourceCode[s.Current+1]
}

func isNumber(c rune) bool {
	return c >= '0' && c <= '9'
}

func isAlphabet(c rune) bool {
	return 	(c >= 'a' && c <= 'z') ||
					(c >= 'A' && c <= 'Z') ||
					c == '_'
}

func isAlphanumeric(c rune) bool {
	return isAlphabet(c) || isNumber(c)
}
