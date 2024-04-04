package lexer

import (
	"bufio"
	"io"
)

type Position struct {
	Line   int
	Column int
}

func NewPosition(line, column int) Position {
	return Position{Line: line, Column: column}
}

type Scanner struct {
	Reader    *bufio.Reader
	Current   rune
	LineCount int
	CharCount int
}

func NewScanner(reader io.Reader) *Scanner {
	return &Scanner{
		Reader:    bufio.NewReader(reader),
		LineCount: 1,
		CharCount: 0,
	}
}

func (s *Scanner) NextRune() error {
	char, _, err := s.Reader.ReadRune()
	if err != nil {
		return err
	}
	if char == '\n' {
		s.LineCount++
		s.CharCount = 0
	} else {
		s.CharCount++
	}
	s.Current = char
	return nil
}

func (s *Scanner) Position() Position {
	return NewPosition(s.LineCount, s.CharCount)
}

func (s *Scanner) Character() rune {
	return s.Current
}

func (s *Scanner) ASCIIchar() int {
	return int(s.Current)
}
