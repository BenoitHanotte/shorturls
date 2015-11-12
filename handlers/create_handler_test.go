package handlers

import (
	"testing"
)

type token struct {
	Counts		int			// after how many calls should the function return that the token if free?
	ExpRandChar	int			// how many random char are expected at the end?
}

// the tokens
// an underscore is a char that should theoretically be replaced by a random char
// underscores are not possible random chars, so they should never be equal to the random chars
var tokens = map[string]token{
	"token0": token{0, 0},
	"token_": token{1, 1},
	"toke__": token{5, 2},
	"tok___": token{7, 3},
	"to____": token{10,4},
	"______": token{20,6},
}

// tmporary var
var lastTokenKey = "";
var curCount = 0;

// check if the token already exists
func mockCheckTokenExists(token string) (bool, error) {

	// if it is a token from the map: reinitiliaze the current count and last token
	_, ok := tokens[token]
	if ok {
		lastTokenKey = token
		curCount = 0;
	}

	tokenData := tokens[lastTokenKey]
	if curCount == tokenData.Counts {
		return false, nil
	}

	curCount++;
	return true, nil
}

// returns true if the token was expected, other false and a human readable error string
func compareTokens(orig string, got string) (bool, string) {
	tokenData := tokens[orig]

	if got[:len(orig)-tokenData.ExpRandChar] == orig[:len(orig)-tokenData.ExpRandChar] &&
	(tokenData.ExpRandChar==0 || got[len(orig)-tokenData.ExpRandChar:] != orig[len(orig)-tokenData.ExpRandChar:]) {
		return true, ""
	}
	return false, "For '"+orig+"', got '"+got+"'"
}

func TestGenerateToken(t *testing.T) {

	token, _ := generateToken("token0", 6, mockCheckTokenExists)
	ok, msg := compareTokens("token0", token)
	if !ok {
		t.Error(msg)
	}

	token, _ = generateToken("token_", 6, mockCheckTokenExists)
	ok, msg = compareTokens("token_", token)
	if !ok {
		t.Error(msg)
	}

	token, _ = generateToken("toke__", 6, mockCheckTokenExists)
	ok, msg = compareTokens("toke__", token)
	if !ok {
		t.Error(msg)
	}

	token, _ = generateToken("tok___", 6, mockCheckTokenExists)
	ok, msg = compareTokens("tok___", token)
	if !ok {
		t.Error(msg)
	}

	token, _ = generateToken("to____", 6, mockCheckTokenExists)
	ok, msg = compareTokens("to____", token)
	if !ok {
		t.Error(msg)
	}

	token, _ = generateToken("______", 6, mockCheckTokenExists)
	ok, msg = compareTokens("______", token)
	if !ok {
		t.Error(msg)
	}
}