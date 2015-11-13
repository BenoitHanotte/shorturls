package handlers

import (
	"testing"
	"strconv"
)

type token struct {
	Original	string		// the suggested token
	Expected	string		// the expected token (a missing letter is a random char)
	ErrorExp    bool		// an error is expected? (eg: too many retires)
	Counts		int			// after how many calls should the lock be acquired
}

// the tokens
// an underscore is a char that should theoretically be replaced by a random char
// underscores are not possible random chars, so they should never be equal to the random chars
var tokens = []token{
	token{"token0", "token0", false, 0},
	token{"token1", "token", false, 1},
	token{"token2", "toke", false, 5},
	token{"token3", "tok", false, 7},
	token{"token4", "to", false, 10},
	token{"token5", "", true, 25},
	token{"", "", false, 0},
	token{"a", "a", false, 0},
	token{"bb", "bb", false, 0},
	token{"cccc", "cccc", false, 0},
	token{"ddddd", "ddddd", false, 0},
	token{"eeee", "eee", false, 2},
	token{"ffff", "", false, 10},
}

// returns true if the token was expected, other false and a human readable error string
func validToken(orig token, got string, gotErr error, expectedLength int) (bool, string) {

	if orig.ErrorExp && gotErr==nil {
		return false, "Should have raised error: for '"+orig.Original+"' got '"+ got+"'"
	} else if orig.ErrorExp && gotErr!=nil {
		return true, ""
	}

	if len(got)!=expectedLength {
		return false, "Wrong length: for '"+orig.Original+"' got '"+
		got+"', expected: '"+orig.Expected+"'"
	}

	for i:=0; i<len(orig.Expected); i++ {
		if orig.Expected[i]!=got[i] {
			return false, "Wrong char at pos "+strconv.Itoa(i)+": for '"+orig.Original+"' got '"+
			got+"', expected: '"+orig.Expected+"'"
		}
	}

	return true, ""
}

func TestGenerateToken(t *testing.T) {

	for _, orig := range tokens {
		randomGenerator := randomTokenGenerator(orig.Original, 6)
		var got string
		var err error
		for i:=0; i<= orig.Counts; i++ { got, err = randomGenerator()}
		ok, msg := validToken(orig, got, err, 6)
		if !ok { t.Error(msg) }
	}
}