package main

import (
	"bufio"
	"fmt"
	. "github.com/flocks/gof/parse"
	"log"
	"os"
	"strings"

	"github.com/urfave/cli/v2"
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
		// Handle error.
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
	files := make(map[string]Filematch)

	lines := strings.Split(string(input), "\n")

	for _, line := range lines {
		file, err := ParseLine(line)
		if err == nil {
			fileWithPWD := updateFilePath(file)
			_, exist := files[fileWithPWD.FilePath]
			if FileExist(fileWithPWD.FilePath) && !exist {
				files[fileWithPWD.FilePath] = fileWithPWD
			}
		}
	}

	for _, val := range files {
		result = append(result, val)
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
