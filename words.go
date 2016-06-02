package tools

import "unicode"

func SplitWords(c rune) bool {
	return unicode.IsSpace(c) || c == ',' || c == '.' || c == '!' || c == '?' || c == ':' || c == ';' ||
		c == '%' || c == '#' || c == '*' || c == '(' || c == ')' || c == '^' || c == '&' || c == '+' ||
		c == '^' || c == '<' || c == '>' || c == '[' || c == ']' || c == '{' || c == '}'
}

