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
	for scanner.Scan() {
		stdin = stdin + "\n" + scanner.Text()
	}

	if scanner.Err() != nil {
		log.Fatal("Error while reading STDIN")
	}
	app := &cli.App{
		Name:  "gof",
		Usage: "extract files from stdin",
		Action: func(cCtx *cli.Context) error {
			files := FindFiles(stdin)
			printFiles(files)
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

	lines := strings.Split(string(input), "\n")

	for _, line := range lines {
		fullLine, _ := ParseLine(line, true)
		nestedLine, _ := ParseLine(line, false)

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

	if hasAtLeastOneError(result) {
		return removeFromList(result, func(file Filematch) bool {
			return file.Col == 0 && file.Line == 0
		})
	}

	return result
}

func hasOnlyFile(file Filematch) bool {
	return file.Col == 0 && file.Line == 0 && file.Desc == "" && file.FilePath != ""
}

func hasError(file Filematch) bool {
	return file.Col != 0 && file.Line != 0 && file.Desc != ""
}

func hasOnlyError(file Filematch) bool {
	return hasError(file) && file.FilePath == ""
}

func hasAtLeastOneError(files []Filematch) bool {
	for _, val := range files {
		if val.Line != 0 || val.Col != 0 {
			return true
		}
	}
	return false
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

func updateFilePath(file Filematch) Filematch {
	if !strings.HasPrefix(file.FilePath, "/") {
		pwd, _ := os.LookupEnv("PWD")
		filePath := pwd + "/" + file.FilePath
		file.FilePath = filePath
	}
	return file
}

func printFiles(files []Filematch) {
	for _, file := range files {
		fmt.Printf(`%v:%v:%v:%v`, file.FilePath, file.Line, file.Col, file.Desc)
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
