package parse

import (
	"regexp"
	"strconv"
	"strings"
)

type Parsed struct {
	parsed string
	rest   string
}

type Filematch struct {
	FilePath string
	Line     int64
	Col      int64
	Desc     string // could potentially hold error/warning from program like linter
}

func (file1 Filematch) CompareWith(file2 Filematch) bool {
	return (file1.Col == file2.Col &&
		file1.Line == file2.Line &&
		file1.FilePath == file2.FilePath &&
		file1.Desc == file2.Desc)
}

func getParser(regexpString string) func(string) Parsed {
	return func(input string) Parsed {
		r, _ := regexp.Compile(regexpString)
		if len(input) > 0 {
			first := string(input[0])
			if r.MatchString(first) {
				return Parsed{
					parsed: first,
					rest:   input[1:],
				}
			}

		}
		return Parsed{parsed: "", rest: input}
	}
}

func ParseSeparator() func(string) Parsed {
	return getParser("[: | ,]")
}

func ParseFilepathChar() func(string) Parsed {
	return getParser("([a-zA-Z-0-9]|/|\\.)")
}
func ParseAnything() func(string) Parsed {
	return getParser("[^$]")
}

func ParseWhitespace() func(string) Parsed {
	return getParser("[[:space:]]")
}

func ParseNumber() func(string) Parsed {
	return getParser("\\d{1}")
}

func ParseFilePath(input string) Parsed {
	p := TakeWhile(input, ParseSeparator())
	file := TakeWhile(p.rest, ParseFilepathChar())

	return Parsed{
		parsed: file.parsed,
		rest:   file.rest,
	}
}
func ParseDesc(input string) Parsed {
	p := TakeWhile(input, ParseSeparator())
	desc := TakeWhile(p.rest, ParseAnything())

	return Parsed{
		parsed: desc.parsed,
		rest:   desc.rest,
	}
}

func ParseFilePosition(input string) Parsed {
	p := TakeWhile(input, ParseSeparator())
	line := TakeWhile(p.rest, ParseNumber())
	p = TakeWhile(line.rest, ParseSeparator())
	col := TakeWhile(p.rest, ParseNumber())

	return Parsed{
		parsed: line.parsed + ":" + col.parsed,
		rest:   col.rest,
	}
}

func ParseLine(input string) (Filematch, error) {
	p := TakeWhile(input, ParseSeparator())
	file := ParseFilePath(p.rest)
	position := ParseFilePosition(file.rest)
	desc := ParseDesc(position.rest)

	_position := strings.Split(position.parsed, ":")
	line := int64(0)
	col := int64(0)
	if len(_position) == 2 {
		line, _ = strconv.ParseInt(_position[0], 10, 0)
		col, _ = strconv.ParseInt(_position[1], 10, 0)
	}
	return Filematch{
		FilePath: file.parsed,
		Line:     line,
		Col:      col,
		Desc:     desc.parsed,
	}, nil

}

func TakeWhile(input string, parser func(i string) Parsed) Parsed {
	p := parser(input)
	current_parsed := p.parsed
	for p.parsed != "" {
		p = parser(p.rest)
		current_parsed = current_parsed + p.parsed
	}

	return Parsed{
		parsed: current_parsed,
		rest:   p.rest,
	}
}
