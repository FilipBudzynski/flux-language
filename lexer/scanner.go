package lexer

import (
	"bufio"
	"io"
	"log"
	"tkom/shared"
)

const EOF rune = -1

type Scanner struct {
	Reader    *bufio.Reader
	Current   rune
	LineCount int
	CharCount int
}

func NewScanner(reader io.Reader) (*Scanner, error) {
	scanner := &Scanner{
		Reader:    bufio.NewReader(reader),
		LineCount: 1,
		CharCount: 0,
	}

	scanner.NextRune()

	return scanner, nil
}

func (s *Scanner) readRune() rune {
	char, _, err := s.Reader.ReadRune()
	if err != nil {
		if err == io.EOF {
			char = EOF
		} else {
			log.Println("Unexpected error while reading source")
			return EOF
		}
	}
	return char
}

func (s *Scanner) NextRune() {
	if s.Current == '\n' {
		s.LineCount++
		s.CharCount = 0
	}

	if s.Current == EOF {
		return
	}

	char := s.readRune()

	if char == '\r' {
		nextChar := s.readRune()
		if nextChar == '\n' {
			char = nextChar
		} else {
			err := s.Reader.UnreadRune()
			if err != nil {
				log.Println("Unexpected error while reading source", err)
				char = EOF
			}
		}
	}

	s.CharCount++
	s.Current = char
}

func (s *Scanner) Position() shared.Position {
	return shared.NewPosition(s.LineCount, s.CharCount)
}

func (s *Scanner) Character() rune {
	return s.Current
}

func (s *Scanner) ASCIIchar() int {
	return int(s.Current)
}
