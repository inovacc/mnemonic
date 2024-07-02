package wordlist

import (
	"fmt"
	"github.com/inovacc/mnemonic/entropy"
	"strconv"
	"strings"
)

// Mnemonic represents a collection of human-readable words
// used for HD wallet seed generation
type Mnemonic struct {
	Words    []string
	Language LanguageStr
}

// New returns a new Mnemonic for the given entropy and language
func New(ent []byte, lang LanguageStr) (*Mnemonic, error) {
	const chunkSize = 11
	bits := entropy.CheckSummed(ent)
	length := len(bits)
	words := make([]string, length/11)
	for i := 0; i < length; i += chunkSize {
		stringVal := string(bits[i : chunkSize+i])
		intVal, err := strconv.ParseInt(stringVal, 2, 64)
		if err != nil {
			return nil, fmt.Errorf("could not convert %s to word index", stringVal)
		}
		words[(chunkSize+i)/11-1] = GetWord(lang, intVal)
	}
	return &Mnemonic{words, lang}, nil
}

// NewRandom returns a new Mnemonic with random entropy of the given length
// in bits
func NewRandom(length int, lang LanguageStr) (*Mnemonic, error) {
	ent, err := entropy.Random(length)
	if err != nil {
		return nil, fmt.Errorf("error generating random entropy: %s", err)
	}
	return New(ent, lang)
}

// Sentence returns a Mnemonic's word collection as a space separated
// sentence
func (m *Mnemonic) Sentence() string {
	if m.Language == Japanese {
		return strings.Join(m.Words, `　`)
	}
	return strings.Join(m.Words, " ")
}

// GenerateSeed returns a seed used for wallet generation per
// BIP-0032 or similar method. The internal Words set
// of the Mnemonic will be used
func (m *Mnemonic) GenerateSeed(passphrase string) *Seed {
	return NewSeed(m.Sentence(), passphrase)
}
