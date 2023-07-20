package main

import (
	"bufio"
	"fmt"
	. "github.com/flocks/gof/parse"
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"strings"
)

var (
	FileExist = _fileExist
)

func main() {
	stdin := ""
	scanner := bufio.NewScanner(os.Stdin)
	var onlyPath int
	for scanner.Scan() {
		stdin = stdin + "\n" + scanner.Text()
	}

	if scanner.Err() != nil {
		log.Fatal("Error while reading STDIN")
	}
	app := &cli.App{
		Name:  "gof",
		Usage: "extract files from stdin",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "onlyPath",
				Aliases: []string{"o"},
				Usage:   "print only files path without line/col/desc",
				Count:   &onlyPath,
			},
		},
		Action: func(cCtx *cli.Context) error {
			files := FindFiles(stdin)
			printFiles(files, onlyPath > 0)
			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func FindFiles(input string) []Filematch {
	var result []Filematch
	var currentFileMatch *Filematch
	hasError := false

	lines := strings.Split(string(input), "\n")

	for _, line := range lines {
		fullLine, _ := ParseLine(line, true)
		nestedLine, _ := ParseLine(line, false)

		if hasFileError(fullLine) || hasFileError(nestedLine) {
			hasError = true
		}

		if hasOnlyError(nestedLine) && currentFileMatch != nil {
			nestedLine.SetFile(currentFileMatch.FilePath)
			result = append(result, nestedLine)
		} else if fullLine.FilePath != "" {
			if !strings.HasPrefix(fullLine.FilePath, "/") {
				pwd, _ := os.LookupEnv("PWD")
				filePath := pwd + "/" + fullLine.FilePath
				fullLine.SetFile(filePath)
			}
			if FileExist(fullLine.FilePath) {
				currentFileMatch = &fullLine
				result = append(result, fullLine)
			}
		}
	}

	// if we parse at least a line with error, we can consider all matches should have line/col/desc
	if hasError {
		return removeFromList(result, func(file Filematch) bool {
			return file.Col == 0 && file.Line == 0
		})
	}

	return result
}

func hasFileError(file Filematch) bool {
	return file.Col != 0 && file.Line != 0 && file.Desc != ""
}

func hasOnlyError(file Filematch) bool {
	return hasFileError(file) && file.FilePath == ""
}

func removeFromList(files []Filematch, predicate func(Filematch) bool) []Filematch {
	var result []Filematch
	for _, file := range files {
		if !predicate(file) {
			result = append(result, file)
		}
	}
	return result
}

func printFiles(files []Filematch, onlyPath bool) {
	for _, file := range files {
		if onlyPath {
			fmt.Printf(`%v`, file.FilePath)
		} else {
			fmt.Printf(`%v:%v:%v:%v`, file.FilePath, file.Line, file.Col, file.Desc)
		}
		fmt.Println()
	}
}

func _fileExist(filePath string) bool {
	if stat, err := os.Stat(filePath); err == nil {
		return !stat.IsDir()
	} else {
		return false
	}
}
