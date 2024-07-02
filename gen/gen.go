//go:build generate

package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"strings"
	"text/template"
)

type LanguageStr string

const (
	English    LanguageStr = "English"
	Spanish    LanguageStr = "Spanish"
	Korean     LanguageStr = "Korean"
	ChineseS   LanguageStr = "ChineseS"
	ChineseT   LanguageStr = "ChineseT"
	Japanese   LanguageStr = "Japanese"
	French     LanguageStr = "French"
	Czech      LanguageStr = "Czech"
	Italian    LanguageStr = "Italian"
	Portuguese LanguageStr = "Portuguese"
)

const templateFile = `package mnemonic

import "math/rand/v2"

type LanguageStr string

const (
	English    LanguageStr = "English"
	Spanish    LanguageStr = "Spanish"
	Korean     LanguageStr = "Korean"
	ChineseS   LanguageStr = "ChineseS"
	ChineseT   LanguageStr = "ChineseT"
	Japanese   LanguageStr = "Japanese"
	French     LanguageStr = "French"
	Czech      LanguageStr = "Czech"
	Italian    LanguageStr = "Italian"
	Portuguese LanguageStr = "Portuguese"
)

var wordLists = WordLists{
	words: make(map[LanguageStr][]string),
}

type WordLists struct {
	words map[LanguageStr][]string
}

func init() {
	wordLists.words[English] = []string{
		{{ .EnglishWords }},
	}
	
	wordLists.words[Spanish] = []string{
		{{ .SpanishWords }},
	}
	
	wordLists.words[Korean] = []string{
		{{ .KoreanWords }},
	}
	
	wordLists.words[ChineseS] = []string{
		{{ .ChineseSWords }},
	}
	
	wordLists.words[ChineseT] = []string{
		{{ .ChineseTWords }},
	}
	
	wordLists.words[Japanese] = []string{
		{{ .JapaneseWords }},
	}
	
	wordLists.words[French] = []string{
		{{ .FrenchWords }},
	}
	
	wordLists.words[Czech] = []string{
		{{ .CzechWords }},
	}
	
	wordLists.words[Italian] = []string{
		{{ .ItalianWords }},
	}
	
	wordLists.words[Portuguese] = []string{
		{{ .PortugueseWords }},
	}
}

func GetWord(lang LanguageStr, idx int64) string {
	return wordLists.words[lang][idx]
}

func RandomWord(lang LanguageStr) string {
	return wordLists.words[lang][rand.IntN(len(wordLists.words[lang]))]
}

func GenerateMnemonic(size int64, lang LanguageStr) []string {
	var result []string
	for _ = range size {
		result = append(result, RandomWord(lang))
	}
	return result
}
`

//go:generate go run gen.go

func main() {
	urls := make(map[LanguageStr]string)
	urls[English] = "https://raw.githubusercontent.com/bitcoin/bips/master/bip-0039/english.txt"
	urls[Spanish] = "https://raw.githubusercontent.com/bitcoin/bips/master/bip-0039/spanish.txt"
	urls[Japanese] = "https://raw.githubusercontent.com/bitcoin/bips/master/bip-0039/japanese.txt"
	urls[Korean] = "https://raw.githubusercontent.com/bitcoin/bips/master/bip-0039/korean.txt"
	urls[ChineseS] = "https://raw.githubusercontent.com/bitcoin/bips/master/bip-0039/chinese_simplified.txt"
	urls[ChineseT] = "https://raw.githubusercontent.com/bitcoin/bips/master/bip-0039/chinese_traditional.txt"
	urls[French] = "https://raw.githubusercontent.com/bitcoin/bips/master/bip-0039/french.txt"
	urls[Italian] = "https://raw.githubusercontent.com/bitcoin/bips/master/bip-0039/italian.txt"
	urls[Czech] = "https://raw.githubusercontent.com/bitcoin/bips/master/bip-0039/czech.txt"
	urls[Portuguese] = "https://raw.githubusercontent.com/bitcoin/bips/master/bip-0039/portuguese.txt"

	words := make(map[LanguageStr][]string)
	for k, url := range urls {
		words[k], _ = downloadFile(url)
	}

	tmpl, err := template.New("wordlist").Parse(templateFile)
	if err != nil {
		panic(err)
	}

	file, err := os.Create("../wordlist.go")
	if err != nil {
		panic(err)
	}

	err = tmpl.Execute(file, map[string]string{
		"EnglishWords":    strings.Join(words[English], ","),
		"SpanishWords":    strings.Join(words[Spanish], ","),
		"KoreanWords":     strings.Join(words[Korean], ","),
		"ChineseSWords":   strings.Join(words[ChineseS], ","),
		"ChineseTWords":   strings.Join(words[ChineseT], ","),
		"JapaneseWords":   strings.Join(words[Japanese], ","),
		"FrenchWords":     strings.Join(words[French], ","),
		"CzechWords":      strings.Join(words[Czech], ","),
		"ItalianWords":    strings.Join(words[Italian], ","),
		"PortugueseWords": strings.Join(words[Portuguese], ","),
	})
	if err != nil {
		panic(err)
	}
}

func downloadFile(url string) ([]string, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	scanner := bufio.NewScanner(resp.Body)
	var words []string
	for scanner.Scan() {
		words = append(words, fmt.Sprintf("\"%s\"", scanner.Text()))
	}

	return words, nil
}
