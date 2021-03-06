package gui

import (
	"fmt"
	"regexp"
	"strings"
)

type codeLines []codeLine

type codeLine struct {
	isSeparator bool
	lineIdx     int
	text        string
}

func (cls codeLines) first() codeLine {
	if len(cls) == 0 {
		return codeLine{}
	}

	for _, l := range cls {
		if !l.isSeparator {
			return l
		}
	}

	return cls[0]
}

func (cls codeLines) last() codeLine {
	l := len(cls)

	if l == 0 {
		return codeLine{}
	}

	for i := range cls {
		l := cls[l-i-1]
		if !l.isSeparator {
			return l
		}
	}

	return cls[l-1]
}

func (cls codeLines) next(curLine int) codeLine {
	if len(cls) == 0 {
		return codeLine{}
	}

	for _, l := range cls {
		if !l.isSeparator && l.lineIdx > curLine {
			return l
		}
	}

	return cls[len(cls)-1]
}

func (cls codeLines) previous(curLine int) codeLine {
	l := len(cls)
	if l == 0 {
		return codeLine{}
	}

	for i := range cls {
		l := cls[l-i-1]

		if !l.isSeparator && l.lineIdx < curLine {
			return l
		}
	}

	return cls[0]
}

func formatDisasm(line string, indentation *int, indentationStep int) string {
	opData := regexp.MustCompile(`OP_DATA_\d+ `)
	line = opData.ReplaceAllString(line, "")

	parts := strings.SplitN(line, ":", 3)
	if len(parts) != 3 {
		return line
	}

	code := strings.TrimSpace(parts[2])

	// Decrease indentation for OP_ELSE and OP_ENDIF statements
	if strings.HasPrefix(code, "OP_ELSE") ||
		strings.HasPrefix(code, "OP_ENDIF") {
		*indentation -= indentationStep
	}

	line = fmt.Sprintf("  %s%s", strings.Repeat(" ", *indentation), code)

	// Increase indentation for all line inside OP_IF and OP_ELSE
	if strings.HasPrefix(code, "OP_IF") ||
		strings.HasPrefix(code, "OP_NOTIF") ||
		strings.HasPrefix(code, "OP_ELSE") {
		*indentation += indentationStep
	}

	return line
}

func isPubkeyScript(line string) bool {
	return strings.HasPrefix(line, "01:")
}

func isSignatureScript(line string) bool {
	return strings.HasPrefix(line, "00:")
}

func isWitnessScript(line string) bool {
	return strings.HasPrefix(line, "02:")
}

func isFirstScriptLine(line string) bool {
	return strings.Contains(line, ":0000: ")
}
