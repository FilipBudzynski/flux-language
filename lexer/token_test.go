package lexer

import "testing"

func TestTypeValueMatch(t *testing.T) {
	t.Run("Valid cases", func(t *testing.T) {
		intToken := NewToken(CONST_INT, Position{}, 42)
		floatToken := NewToken(CONST_FLOAT, Position{}, 3.14)
		stringToken := NewToken(CONST_STRING, Position{}, "hello")
		trueToken := NewToken(CONST_TRUE, Position{}, true)
		falseToken := NewToken(CONST_FALSE, Position{}, false)

		if intToken == nil || floatToken == nil || stringToken == nil || trueToken == nil || falseToken == nil {
			t.Error("Failed to create token")
		}
	})

	t.Run("Invalid cases", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("Expected panic, but got nil")
			}
		}()

		invalidIntToken := NewToken(CONST_INT, Position{}, "invalid")     // Should panic
		invalidFloatToken := NewToken(CONST_FLOAT, Position{}, "invalid") // Should panic
		invalidStringToken := NewToken(CONST_STRING, Position{}, 42)      // Should panic
		invalidIdentifierToken := NewToken(IDENTIFIER, Position{}, 42)    // Should panic
		invalidTrueToken := NewToken(CONST_TRUE, Position{}, false)       // Should panic
		invalidFalseToken := NewToken(CONST_FALSE, Position{}, true)      // Should panic

		_ = invalidIntToken
		_ = invalidFloatToken
		_ = invalidStringToken
		_ = invalidIdentifierToken
		_ = invalidTrueToken
		_ = invalidFalseToken
	})
}
