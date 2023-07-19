package parse

import (
	"log"
	"testing"
)

func TestParseSeparator(t *testing.T) {
	parsed := ParseSeparator()(":test")
	expect(parsed, Parsed{parsed: ":", rest: "test"}, "separator")
}
func TestParseNotSeparator(t *testing.T) {
	parsed := ParseSeparator()("|test")
	expect(parsed, Parsed{parsed: "", rest: "|test"}, "separator")
}

func TestParseWhitespace(t *testing.T) {
	parsed := ParseWhitespace()(" bob l")
	expect(parsed, Parsed{parsed: " ", rest: "bob l"}, "whitespace")
}

func TestParseNumber(t *testing.T) {
	parsed := ParseNumber()("99rest9")
	expect(parsed, Parsed{parsed: "9", rest: "9rest9"}, "number")
}
func TestParseNumberNil(t *testing.T) {
	parsed := ParseNumber()("nodigits")
	expect(parsed, Parsed{parsed: "", rest: "nodigits"}, "number")
}

func TestTakeWhileNumber(t *testing.T) {
	parsed := TakeWhile("888rest 1", ParseNumber())
	expect(parsed, Parsed{parsed: "888", rest: "rest 1"}, "take-while-number")
}
func TestTakeWhileSeparator(t *testing.T) {
	parsed := TakeWhile("::::lol bob", ParseSeparator())
	expect(parsed, Parsed{parsed: "::::", rest: "lol bob"}, "take-while-separator")
}

func TestParsePosition(t *testing.T) {
	parsed := ParseFilePosition("     40:39 ")
	expect(parsed, Parsed{parsed: "40:39", rest: " "}, "parse-position")
}

func TestParseFile(t *testing.T) {
	parsed := ParseFilePath("     /home/bob/file.txt rest-of-the-line")
	expect(parsed, Parsed{parsed: "/home/bob/file.txt", rest: " rest-of-the-line"}, "parse-file")
}

func TestParseDesc(t *testing.T) {
	parsed := ParseDesc("this a description of the error")
	expect(parsed, Parsed{parsed: "this a description of the error", rest: ""}, "parse-desc")
}
func TestParseFullLine(t *testing.T) {
	file, _ := ParseLine(" /home/bob/file.txt 19:25 this is the error description", true)
	expectFile(file, Filematch{
		FilePath: "/home/bob/file.txt",
		Line:     19,
		Col:      25,
		Desc:     " this is the error description",
	})
}
func TestParseFullGrepLine(t *testing.T) {
	file, _ := ParseLine("src/components/DeviceInteraction/index.jsx:62:11:  appVersion: string,", true)
	expectFile(file, Filematch{
		FilePath: "src/components/DeviceInteraction/index.jsx",
		Line:     62,
		Col:      11,
		Desc:     "  appVersion: string,",
	})
}
func TestParseLintLineWithoutFile(t *testing.T) {
	file, _ := ParseLine(" 23:7  warning  'a' is assigned a value but never used  @typescript-eslint/no-unused-vars", false)
	expectFile(file, Filematch{
		Line: 23,
		Col:  7,
		Desc: "  warning  'a' is assigned a value but never used  @typescript-eslint/no-unused-vars",
	})
}
func TestParseLineWithOnlyFile(t *testing.T) {
	file, _ := ParseLine(" /home/bob/file.txt ", true)
	expectFile(file, Filematch{
		FilePath: "/home/bob/file.txt",
	})
}

func TestParseLineWithOnlyFileWithGarbage(t *testing.T) {
	file, _ := ParseLine(" /home/bob/file.txt    | 4 +--- ", true)
	expectFile(file, Filematch{
		FilePath: "/home/bob/file.txt",
	})
}

func expect(result Parsed, expected Parsed, label string) {
	if result.parsed != expected.parsed {
		log.Fatalf(`Invalid parsed value for %v, value: %v, expected: %v`, label, result.parsed, expected.parsed)
	}
	if result.rest != expected.rest {
		log.Fatalf(`Invalid rest value for %v, value: %v, expected: %v`, label, result.rest, expected.rest)
	}
}

func expectFile(result Filematch, expected Filematch) {
	if !expected.CompareWith(result) {
		log.Printf(`got %v %v,%v %v`, result.FilePath, result.Line, result.Col, result.Desc)
		log.Fatalf(`expected %v %v,%v %v`, expected.FilePath, expected.Line, expected.Col, expected.Desc)
	}
}
