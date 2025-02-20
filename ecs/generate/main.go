package main

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"text/template"
)

const queryTemplate = "./query.gtpl"
const queryOutput = "./query.go"

func main() {
	funcMap := template.FuncMap{
		"makeRange":    makeRange,
		"lowerLetters": lowerLetters,
		"upperLetters": upperLetters,
		"concat":       concat,
		"join":         join,
	}

	t, err := template.New("query").Funcs(funcMap).ParseFiles(queryTemplate)
	if err != nil {
		panic(err)
	}

	var result bytes.Buffer
	err = t.ExecuteTemplate(&result, "query", nil)
	if err != nil {
		panic(err)
	}
	os.WriteFile(queryOutput, result.Bytes(), 0644)
}

func makeRange(min, max int) []int {
	r := make([]int, max-min+1)
	for i := range r {
		r[i] = min + i
	}
	return r
}

func lowerLetters(n int) []string {
	letters := make([]string, n)
	for i := 0; i < n; i++ {
		letters[i] = string(rune('a' + i))
	}
	return letters
}

func upperLetters(n int) []string {
	letters := make([]string, n)
	for i := 0; i < n; i++ {
		letters[i] = string(rune('A' + i))
	}
	return letters
}

func concat(args ...interface{}) string {
	var result strings.Builder
	for _, arg := range args {
		result.WriteString(fmt.Sprintf("%v", arg))
	}
	return result.String()
}

func join(before, sep, after string, args []string) string {
	return fmt.Sprintf("%s%s%s", before, strings.Join(args, sep), after)
}
