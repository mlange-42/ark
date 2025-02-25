package main

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"text/template"
)

type genFile struct {
	Source string
	Target string
}

var files = []genFile{
	{"./filter.go.template", "../filter_gen.go"},
	{"./query.go.template", "../query_gen.go"},
	{"./query_test.go.template", "../query_gen_test.go"},
	{"./maps.go.template", "../maps_gen.go"},
	{"./exchange.go.template", "../exchange_gen.go"},
}

func main() {
	funcMap := template.FuncMap{
		"makeRange":    makeRange,
		"lowerLetters": lowerLetters,
		"upperLetters": upperLetters,
		"concat":       concat,
		"join":         join,
		"arguments":    arguments,
		"blanks":       blanks,
	}

	for _, file := range files {
		t, err := template.New("template").Funcs(funcMap).ParseFiles(file.Source)
		if err != nil {
			panic(err)
		}

		var result bytes.Buffer
		err = t.ExecuteTemplate(&result, "template", nil)
		if err != nil {
			panic(err)
		}
		os.WriteFile(file.Target, result.Bytes(), 0644)
	}
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

func arguments(names []string, types []string, prefix string) string {
	str := make([]string, len(names))
	for i, name := range names {
		str[i] = fmt.Sprintf("%s *%s%s", name, prefix, types[i])
	}
	return strings.Join(str, ", ")
}

func blanks(count int) string {
	b := strings.Builder{}
	for i := range count {
		if i > 0 {
			b.WriteString(", ")
		}
		b.WriteString("_")
	}
	return b.String()
}
