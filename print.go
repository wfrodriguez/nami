package main

import (
	"fmt"
	"regexp"
	"strings"
)

const Green = "\033[32m"
const Blue = "\033[34m"
const Yellow = "\033[33m"
const Red = "\033[31m"
const Reset = "\033[0m"

var reKeywords = regexp.MustCompile(
	`\b(break|default|func|interface|select|case|defer|go|map|struct|chan|else|goto|package|switch|const|fallthrough|if|range|type|continue|for|import|return|var)\b`,
)

var reOperators = regexp.MustCompile(
	`(~|[\[\]\(\)\{\}\.\+\\\-*/%&\|^!,;:<>]+)`,
)

var reTypes = regexp.MustCompile(
	`\b(any|bool|int|int8|int16|int32|int64|uint|uint8|uint16|uint32|uint64|uintptr|float32|float64|complex64|complex128|string|nil)\b`,
)

func PrintDef(def string) {
	lines := strings.Split(def, "\n")

	pf := reOperators.ReplaceAllString(lines[0], Red+"$1"+Reset)
	pf = reKeywords.ReplaceAllString(pf, Yellow+"$1"+Reset)
	pf = reTypes.ReplaceAllString(pf, Blue+"$1"+Reset)
	df := reOperators.ReplaceAllString(lines[2], Red+"$1"+Reset)
	df = reKeywords.ReplaceAllString(df, Yellow+"$1"+Reset)
	df = reTypes.ReplaceAllString(df, Blue+"$1"+Reset)

	lines = lines[3:]
	fmt.Printf("%s╭─┤ %s %s├────────%s\n", Green, df, Green, Reset)
	fmt.Printf("%s│%s%s\n", Green, pf, Reset)
	fmt.Printf("%s│%s\n", Green, Reset)
	for _, l := range lines {
		fmt.Printf("%s│%s%s\n", Green, Reset, strings.Replace(l, "    ", "", 1))
	}
	fmt.Printf("%s╰────%s\n", Green, Reset)
}

/*

╭ ╰ ─ ┤ ├ │

╭─┤ ├─
│
│
╰────
~[]{},,;:()
*/
