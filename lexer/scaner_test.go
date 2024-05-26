package lexer

import (
	"reflect"
	"strings"
	"testing"
	"tkom/shared"
)

func TestScanner(t *testing.T) {
	testCases := []struct {
		input         string
		expectedRunes []rune
		expectedPos   []shared.Position
	}{
		{
			input:         "abc",
			expectedRunes: []rune{'a', 'b', 'c', EOF},
			expectedPos:   []shared.Position{{1, 1}, {1, 2}, {1, 3}, {1, 4}},
		},
		{
			input:         "hello\nworld",
			expectedRunes: []rune{'h', 'e', 'l', 'l', 'o', '\n', 'w', 'o', 'r', 'l', 'd', EOF},
			expectedPos:   []shared.Position{{1, 1}, {1, 2}, {1, 3}, {1, 4}, {1, 5}, {1, 6}, {2, 1}, {2, 2}, {2, 3}, {2, 4}, {2, 5}, {2, 6}},
		},
		{
			input:         "line 1\nline 2\nline 3\n",
			expectedRunes: []rune{'l', 'i', 'n', 'e', ' ', '1', '\n', 'l', 'i', 'n', 'e', ' ', '2', '\n', 'l', 'i', 'n', 'e', ' ', '3', '\n', EOF},
			expectedPos:   []shared.Position{{1, 1}, {1, 2}, {1, 3}, {1, 4}, {1, 5}, {1, 6}, {1, 7}, {2, 1}, {2, 2}, {2, 3}, {2, 4}, {2, 5}, {2, 6}, {2, 7}, {3, 1}, {3, 2}, {3, 3}, {3, 4}, {3, 5}, {3, 6}, {3, 7}, {4, 1}},
		},
		{
			input:         "hello\n\n\n\n\n",
			expectedRunes: []rune{'h', 'e', 'l', 'l', 'o', '\n', '\n', '\n', '\n', '\n', EOF},
			expectedPos:   []shared.Position{{1, 1}, {1, 2}, {1, 3}, {1, 4}, {1, 5}, {1, 6}, {2, 1}, {3, 1}, {4, 1}, {5, 1}, {6, 1}},
		},
		{
			input:         "\r\nt\t\r\n",
			expectedRunes: []rune{'\n', 't', '\t', '\n', EOF},
			expectedPos:   []shared.Position{{1, 1}, {2, 1}, {2, 2}, {2, 3}, {3, 1}},
		},
		{
			input:         "\r\tt\t\r\n",
			expectedRunes: []rune{'\r', '\t', 't', '\t', '\n', EOF},
			expectedPos:   []shared.Position{{1, 1}, {1, 2}, {1, 3}, {1, 4}, {1, 5}, {2, 1}},
		},
	}

	for _, tc := range testCases {
		reader := strings.NewReader(tc.input)
		scanner, _ := NewScanner(reader)

		var runes []rune
		var positions []shared.Position
		for {
			char := scanner.Character()
			pos := scanner.Position()

			runes = append(runes, char)
			positions = append(positions, pos)

			if char == EOF {
				break
			}

			scanner.NextRune()
		}

		if !reflect.DeepEqual(runes, tc.expectedRunes) {
			t.Errorf("Input: %s\nExpected Runes: %+v\nGot Runes: %+v\n", tc.input, tc.expectedRunes, runes)
		}

		if !reflect.DeepEqual(positions, tc.expectedPos) {
			t.Errorf("Input: %s\nExpected Positions: %+v\nGot Positions: %+v\n", tc.input, tc.expectedPos, positions)
		}
	}
}

func TestScannerMultipleEOF(t *testing.T) {
	input := ""
	expectedRunes := []rune{EOF, EOF, EOF, EOF, EOF}
	expectedPos := []shared.Position{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}

	scanner, _ := NewScanner(strings.NewReader(input))

	for i, expectedRune := range expectedRunes {
		scanner.NextRune()
		actualPos := scanner.Position()

		actualRune := scanner.Character()
		if actualRune != expectedRune {
			t.Errorf("Expected rune %q at position %d, got %q", expectedRune, i, actualRune)
		}

		pos := scanner.Position()
		if pos != expectedPos[i] || pos != actualPos {
			t.Errorf("Expected position %v at position %d, got %v", expectedPos[i], i, pos)
		}
	}
}
