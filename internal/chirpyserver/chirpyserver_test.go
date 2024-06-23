package chirpyserver

import (
	"testing"
)

func TestCleanChirp(t *testing.T) {
	cases := []struct {
		input       string
		expectedStr string
	}{
		{
			input:       "A clean chirp",
			expectedStr: "A clean chirp",
		},
		{
			input:       "A fornax dirty chirp",
			expectedStr: "A ****** dirty chirp",
		},
		{
			input:       "A fornax really kerfuffle  dirty chirp",
			expectedStr: "A ****** really *********  dirty chirp",
		},
		{
			input:       "Capitalized ShArBeRt chirp",
			expectedStr: "Capitalized ******** chirp",
		},
		{
			input:       "Punctuated sharbert! Fornax, and kerfuffle?",
			expectedStr: "Punctuated sharbert! Fornax, and kerfuffle?",
		},
		{
			input:       "",
			expectedStr: "",
		},
	}

	for _, cs := range cases {
		actual := cleanChirp(cs.input)
		if actual != cs.expectedStr {
			t.Errorf("Unequal chirps: '%v', '%v'", actual, cs.expectedStr)
		}
	}
}
